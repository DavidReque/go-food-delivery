package consumer

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/consumer"
	consumertracing "github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/otel/tracing/consumer"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/pipeline"
	messagingTypes "github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/types"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/utils"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/metadata"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/serializer"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/config"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/consumer/configurations"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/rabbitmqErrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/types"
	errorUtils "github.com/DavidReque/go-food-delivery/internal/pkg/utils/errorutils"

	"emperror.dev/errors"
	"github.com/ahmetb/go-linq/v3"
	"github.com/avast/retry-go"
	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

const (
	retryAttempts = 3                      // the number of times to retry a message
	retryDelay    = 300 * time.Millisecond // the delay between retry attempts
)

// retry options for the consumer
var retryOptions = []retry.Option{
	retry.Attempts(retryAttempts),       // the number of times to retry a message
	retry.Delay(retryDelay),             // the delay between retry attempts
	retry.DelayType(retry.BackOffDelay), // the delay type
	retry.LastErrorOnly(true),           // only return the last error
	retry.RetryIf(func(err error) bool {
		return err != nil
	}), // retry if the error is not nil
}

type rabbitMQConsumer struct {
	rabbitmqConsumerOptions *configurations.RabbitMQConsumerConfiguration
	connection              types.IConnection
	handlerDefault          consumer.ConsumerHandler
	channel                 *amqp091.Channel // the channel to use for the consumer
	deliveryRoutines        chan struct{}    // chan should init before using channel
	messageSerializer       serializer.MessageSerializer
	logger                  logger.Logger
	rabbitmqOptions         *config.RabbitmqOptions
	ErrChan                 chan error
	wg                      sync.WaitGroup             // wait group to wait for the consumer to finish
	handlers                []consumer.ConsumerHandler // the handlers to use for the consumer
	pipelines               []pipeline.ConsumerPipeline
	handlersLock            sync.Mutex // lock to protect the handlers
	isConsumedNotifications []func(message messagingTypes.IMessage)
}

// NewRabbitMQConsumer create a new generic RabbitMQ consumer
func NewRabbitMQConsumer(
	rabbitmqOptions *config.RabbitmqOptions,
	connection types.IConnection,
	consumerConfiguration *configurations.RabbitMQConsumerConfiguration,
	messageSerializer serializer.MessageSerializer,
	logger logger.Logger,
	isConsumedNotifications ...func(message messagingTypes.IMessage),
) (consumer.Consumer, error) {
	if consumerConfiguration == nil {
		return nil, errors.New("consumer configuration is required")
	}

	if consumerConfiguration.ConsumerMessageType == nil {
		return nil, errors.New(
			"consumer ConsumerMessageType property is required",
		)
	}

	// create a channel to limit the number of concurrent deliveries
	// this is used to prevent the consumer from processing more than the concurrency limit
	// ej: if the concurrency limit is 10, and the consumer is processing 10 messages,
	// the consumer will not process any more messages until the 10 messages are processed
	deliveryRoutines := make(
		chan struct{},
		consumerConfiguration.ConcurrencyLimit,
	)
	cons := &rabbitMQConsumer{
		messageSerializer:       messageSerializer,
		rabbitmqOptions:         rabbitmqOptions,
		logger:                  logger,
		rabbitmqConsumerOptions: consumerConfiguration,
		deliveryRoutines:        deliveryRoutines,
		ErrChan:                 make(chan error),
		connection:              connection,
		handlers:                consumerConfiguration.Handlers,
		pipelines:               consumerConfiguration.Pipelines,
	}

	cons.isConsumedNotifications = isConsumedNotifications

	return cons, nil
}

// IsConsumed is a function that adds a notification function to the consumer
// this is used to notify the consumer when a message is consumed
// this is useful for logging, metrics, etc.
func (r *rabbitMQConsumer) IsConsumed(h func(message messagingTypes.IMessage)) {
	r.isConsumedNotifications = append(r.isConsumedNotifications, h)
}

func (r *rabbitMQConsumer) Start(ctx context.Context) error {
	// https://github.com/rabbitmq/rabbitmq-tutorials/blob/master/go/receive.go
	// 1. Determine the topology of the exchange, queue and binding
	if r.connection == nil {
		return errors.New("connection is required")
	}

	var exchange string
	var queue string
	var routingKey string

	// if the exchange name is not empty, use the exchange name
	if r.rabbitmqConsumerOptions.ExchangeOptions.Name != "" {
		exchange = r.rabbitmqConsumerOptions.ExchangeOptions.Name
	} else { // if the exchange name is empty, use the topic or exchange name from the consumer message type
		exchange = utils.GetTopicOrExchangeNameFromType(r.rabbitmqConsumerOptions.ConsumerMessageType)
	}

	// if the routing key is not empty, use the routing key
	if r.rabbitmqConsumerOptions.BindingOptions.RoutingKey != "" {
		routingKey = r.rabbitmqConsumerOptions.BindingOptions.RoutingKey
	} else {
		// if the routing key is empty, use the routing key from the consumer message type
		routingKey = utils.GetRoutingKeyFromType(r.rabbitmqConsumerOptions.ConsumerMessageType)
	}

	// if the queue name is not empty, use the queue name
	if r.rabbitmqConsumerOptions.QueueOptions.Name != "" {
		queue = r.rabbitmqConsumerOptions.QueueOptions.Name
	} else {
		// if the queue name is empty, use the queue name from the consumer message type
		queue = utils.GetQueueNameFromType(r.rabbitmqConsumerOptions.ConsumerMessageType)
	}

	// 2. Prepare the consumer and topology declaration
	r.reConsumeOnDropConnection(ctx)

	// get a new channel on the connection - channel is unique for each consumer
	ch, err := r.connection.Channel()
	if err != nil {
		return rabbitmqErrors.ErrDisconnected
	}
	r.channel = ch

	// The prefetch count tells the Rabbit connection how many messages to retrieve from the server per request.
	// ej: if the concurrency limit is 10, and the prefetch count is 10, the consumer will retrieve 100 messages from the server per request.
	prefetchCount := r.rabbitmqConsumerOptions.ConcurrencyLimit * r.rabbitmqConsumerOptions.PrefetchCount
	if err := r.channel.Qos(prefetchCount, 0, false); err != nil {
		return err
	}

	// declare the exchange
	// if the exchange already exists, it will not be declared again
	err = r.channel.ExchangeDeclare(
		exchange,
		string(r.rabbitmqConsumerOptions.ExchangeOptions.Type),
		r.rabbitmqConsumerOptions.ExchangeOptions.Durable,
		r.rabbitmqConsumerOptions.ExchangeOptions.AutoDelete,
		false,
		r.rabbitmqConsumerOptions.NoWait,
		r.rabbitmqConsumerOptions.ExchangeOptions.Args,
	)
	if err != nil {
		return err
	}

	// declare the queue
	// if the queue already exists, it will not be declared again
	_, err = r.channel.QueueDeclare(
		queue, // queue name
		r.rabbitmqConsumerOptions.QueueOptions.Durable,    // durable
		r.rabbitmqConsumerOptions.QueueOptions.AutoDelete, // auto delete
		r.rabbitmqConsumerOptions.QueueOptions.Exclusive,  // exclusive
		r.rabbitmqConsumerOptions.NoWait,                  // no wait
		r.rabbitmqConsumerOptions.QueueOptions.Args,
	) // arguments
	if err != nil {
		return err
	}

	// bind the queue to the exchange
	err = r.channel.QueueBind(
		queue, // queue name
		routingKey,
		exchange,
		r.rabbitmqConsumerOptions.NoWait,
		r.rabbitmqConsumerOptions.BindingOptions.Args, // arguments
	)
	if err != nil {
		return err
	}

	// 3. Consume the messages
	// consume the messages
	msgs, err := r.channel.Consume(
		queue,                                // queue name
		r.rabbitmqConsumerOptions.ConsumerId, // consumer id
		r.rabbitmqConsumerOptions.AutoAck,    // When autoAck (also known as noAck) is true, the server will acknowledge deliveries to this consumer prior to writing the delivery to the network. When autoAck is true, the consumer should not call Delivery.Ack.
		r.rabbitmqConsumerOptions.QueueOptions.Exclusive,
		r.rabbitmqConsumerOptions.NoLocal,
		r.rabbitmqConsumerOptions.NoWait,
		nil,
	)
	if err != nil {
		return err
	}

	// This channel will receive a notification when a channel closed event happens.
	// https://github.com/streadway/amqp/blob/v1.0.0/channel.go#L447
	// https://github.com/rabbitmq/amqp091-go/blob/main/example_client_test.go#L75
	chClosedCh := make(chan *amqp091.Error, 1)
	ch.NotifyClose(chClosedCh)

	// https://blog.boot.dev/golang/connecting-to-rabbitmq-in-golang/
	// https://levelup.gitconnected.com/connecting-a-service-in-golang-to-a-rabbitmq-server-835294d8c914
	// https://medium.com/@dhanushgopinath/automatically-recovering-rabbitmq-connections-in-go-applications-7795a605ca59
	// https://github.com/rabbitmq/amqp091-go/blob/main/_examples/pubsub/pubsub.go
	for i := 0; i < r.rabbitmqConsumerOptions.ConcurrencyLimit; i++ {
		r.logger.Infof("Processing messages on thread %d", i)
		go func() {
			for { // infinite loop to consume messages
				select {
				// if the context is done, shutdown the consumer
				case <-ctx.Done():
					r.logger.Info("shutting down consumer")
					return
				case amqErr := <-chClosedCh:
					// This case handles the event of closed channel e.g. abnormal shutdown
					r.logger.Errorf("AMQP Channel closed due to: %s", amqErr)

					// Re-set channel to receive notifications
					chClosedCh = make(chan *amqp091.Error, 1)
					ch.NotifyClose(chClosedCh)

					// if the channel is closed, shutdown the consumer
					if amqErr != nil {
						r.logger.Errorf("AMQP Channel closed due to: %s", amqErr)
					}

				case msg, ok := <-msgs:
					if !ok {
						r.logger.Info("consumer connection dropped")
						return
					}

					// handle received message and remove message form queue with a manual ack
					r.handleReceived(ctx, msg)

				}
			}
		}()
	}

	return nil
}

func (r *rabbitMQConsumer) Stop() error {
	// cancel the arrival of new messages
	if r.channel != nil && !r.channel.IsClosed() {
		r.channel.Cancel(r.rabbitmqConsumerOptions.ConsumerId, false)
	}

	// 2. Esperar a que todas las goroutines terminen
	// wait all goroutines to finish
	r.wg.Wait() // block until the wait group is zero

	// close the channel (no need defer)
	if r.channel != nil && !r.channel.IsClosed() {
		r.channel.Close()
	}

	return nil
}

func (r *rabbitMQConsumer) ConnectionHandler(handler consumer.ConsumerHandler) {
	r.handlersLock.Lock()
	defer r.handlersLock.Unlock()
	r.handlerDefault = handler
}

// ConnectHandler connects a handler to the consumer
func (r *rabbitMQConsumer) ConnectHandler(handler consumer.ConsumerHandler) {
	r.handlersLock.Lock()
	defer r.handlersLock.Unlock()
	r.handlers = append(r.handlers, handler)
}

// GetName returns the name of the consumer
func (r *rabbitMQConsumer) GetName() string {
	return r.rabbitmqConsumerOptions.Name
}

func (r *rabbitMQConsumer) reConsumeOnDropConnection(ctx context.Context) {
	go func() {
		defer errorUtils.HandlePanic()
		for {
			select {
			case reconnect := <-r.connection.ReconnectedChannel():
				if reflect.ValueOf(reconnect).IsValid() {
					r.logger.Info("reconsume_on_drop_connection started")
					err := r.Start(ctx)
					if err != nil {
						r.logger.Error(
							"reconsume_on_drop_connection finished with error: %v",
							err,
						)
						continue
					}
					r.logger.Info(
						"reconsume_on_drop_connection finished successfully",
					)
					return
				}
			case <-ctx.Done():
				r.logger.Info("reconsume_on_drop_connection finished due to context done")
				return
			}
		}
	}()
}

// handleReceived is a function that handles the received message
// it creates a consume context, creates a consumer span, and handles the message
func (r *rabbitMQConsumer) handleReceived(
	ctx context.Context,
	delivery amqp091.Delivery,
) {
	// for ensuring our handlers execute completely after shutdown
	r.deliveryRoutines <- struct{}{}

	defer func() { <-r.deliveryRoutines }()

	var meta metadata.Metadata
	if delivery.Headers != nil {
		meta = metadata.MapToMetadata(delivery.Headers)
	}

	// start the consumer span
	consumerTraceOption := &consumertracing.ConsumerTracingOptions{
		MessagingSystem: "rabbitmq",
		DestinationKind: "queue",
		Destination:     r.rabbitmqConsumerOptions.QueueOptions.Name,
		OtherAttributes: []attribute.KeyValue{
			semconv.MessagingRabbitmqDestinationRoutingKey(delivery.RoutingKey),
		},
	}
	ctx, beforeConsumeSpan := consumertracing.StartConsumerSpan(
		ctx,
		&meta,
		string(delivery.Body),
		consumerTraceOption,
	)

	consumeContext, err := r.createConsumeContext(delivery)
	if err != nil {
		r.logger.Error(
			consumertracing.FinishConsumerSpan(beforeConsumeSpan, err),
		)
		return
	}

	var ack func()
	var nack func()

	// if auto-ack is enabled we should not call Ack method manually it could create some unexpected errors
	if r.rabbitmqConsumerOptions.AutoAck == false {
		// ack is a function that acknowledges the message
		// ack is used to tell the broker that the message was processed correctly
		ack = func() {
			if err := delivery.Ack(false); err != nil {
				r.logger.Error(
					"error sending ACK to RabbitMQ consumer: %v",
					consumertracing.FinishConsumerSpan(beforeConsumeSpan, err),
				)
				return
			}
			_ = consumertracing.FinishConsumerSpan(beforeConsumeSpan, nil)
			if len(r.isConsumedNotifications) > 0 {
				for _, notification := range r.isConsumedNotifications {
					if notification != nil {
						notification(consumeContext.Message())
					}
				}
			}
		}

		// nack is a function that negatively acknowledges the message
		// nack is used to tell the broker that the message was not processed correctly
		nack = func() {
			if err := delivery.Nack(false, true); err != nil {
				r.logger.Error(
					"error in sending Nack to RabbitMQ consumer: %v",
					consumertracing.FinishConsumerSpan(beforeConsumeSpan, err),
				)
				return
			}
			_ = consumertracing.FinishConsumerSpan(beforeConsumeSpan, nil)
		}
	}

	r.handle(ctx, ack, nack, consumeContext)
}

// handle is a function that handles the message
// it runs the handlers, and if there is an error, it nacks the message
// if there is no error, it acks the message
func (r *rabbitMQConsumer) handle(
	ctx context.Context,
	ack func(), // the function to acknowledge the message
	nack func(), // the function to negatively acknowledge the message
	messageConsumeContext messagingTypes.MessageConsumeContext, // the context of the message
) {
	var err error
	for _, handler := range r.handlers {
		err = r.runHandlersWithRetry(ctx, handler, messageConsumeContext)
		if err != nil {
			break // break the loop if there is an error
		}
	}

	if err != nil {
		r.logger.Error(
			"[rabbitMQConsumer.Handle] error in handling consume message of RabbitmqMQ, prepare for nacking message",
		)
		if nack != nil && r.rabbitmqConsumerOptions.AutoAck == false {
			nack()
		}
	} else if err == nil && ack != nil && r.rabbitmqConsumerOptions.AutoAck == false {
		ack()
	}
}

func (r *rabbitMQConsumer) runHandlersWithRetry(
	ctx context.Context,
	handler consumer.ConsumerHandler,
	messageConsumeContext messagingTypes.MessageConsumeContext,
) error {
	// if there is an error, it will retry the handler
	err := retry.Do(func() error {
		var lastHandler pipeline.ConsumerHandlerFunc

		if r.pipelines != nil && len(r.pipelines) > 0 {
			// reverse the pipelines
			// the last pipeline is the first one
			// the first pipeline is the last one
			reversPipes := r.reversOrder(r.pipelines)

			// create the last handler
			// the last handler is the handler that will be called last
			// the last handler is the handler that will be called last
			lastHandler = func(ctx context.Context) error {
				handler := handler.(consumer.ConsumerHandler)
				return handler.Handle(ctx, messageConsumeContext)
			}

			aggregateResult := linq.From(reversPipes).
				AggregateWithSeedT(lastHandler, func(next pipeline.ConsumerHandlerFunc, pipe pipeline.ConsumerPipeline) pipeline.ConsumerHandlerFunc {
					pipeValue := pipe
					nexValue := next

					var handlerFunc pipeline.ConsumerHandlerFunc = func(ctx context.Context) error {
						genericContext, ok := messageConsumeContext.(messagingTypes.MessageConsumeContext)
						if ok {
							return pipeValue.Handle(
								ctx,
								genericContext,
								nexValue,
							)
						}
						return pipeValue.Handle(
							ctx,
							messageConsumeContext.(messagingTypes.MessageConsumeContext),
							nexValue,
						)
					}
					return handlerFunc
				})

			v := aggregateResult.(pipeline.ConsumerHandlerFunc)
			err := v(ctx)
			if err != nil {
				return errors.Wrap(
					err,
					"error handling consumer handlers pipeline",
				)
			}
			return nil
		} else {
			err := handler.Handle(ctx, messageConsumeContext.(messagingTypes.MessageConsumeContext))
			if err != nil {
				return err
			}
		}
		return nil
	}, append(retryOptions, retry.Context(ctx))...)

	return err
}

func (r *rabbitMQConsumer) createConsumeContext(
	delivery amqp091.Delivery,
) (messagingTypes.MessageConsumeContext, error) {
	message := r.deserializeData(
		delivery.ContentType,
		delivery.Type,
		delivery.Body,
	)

	var meta metadata.Metadata
	if delivery.Headers != nil {
		meta = metadata.MapToMetadata(delivery.Headers)
	}

	consumeContext := messagingTypes.NewMessageConsumeContext(
		message,
		meta,
		delivery.ContentType,
		delivery.Type,
		delivery.Timestamp,
		delivery.DeliveryTag,
		delivery.MessageId,
		delivery.CorrelationId,
	)
	return consumeContext, nil
}

// deserializeData is a function that deserializes the data
// it deserializes the data from the body
// it returns the deserialized data
func (r *rabbitMQConsumer) deserializeData(
	contentType string,
	eventType string,
	body []byte,
) messagingTypes.IMessage {
	if contentType == "" {
		contentType = "application/json"
	}

	if body == nil || len(body) == 0 {
		r.logger.Error("message body is nil or empty in the consumer")
		return nil
	}

	if contentType == "application/json" {
		// r.rabbitmqConsumerOptions.ConsumerMessageType --> actual type
		// deserialize, err := r.messageSerializer.DeserializeType(body, r.rabbitmqConsumerOptions.ConsumerMessageType, contentType)
		deserialize, err := r.messageSerializer.Deserialize(
			body,
			eventType,
			contentType,
		) // or this to explicit type deserialization
		if err != nil {
			r.logger.Errorf(
				fmt.Sprintf(
					"error in deserilizng of type '%s' in the consumer",
					eventType,
				),
			)
			return nil
		}

		return deserialize
	}

	return nil
}

// reversOrder is a function that reverses the order of the pipelines
// it reverses the order of the pipelines
// it returns the reversed pipelines
func (r *rabbitMQConsumer) reversOrder(
	values []pipeline.ConsumerPipeline,
) []pipeline.ConsumerPipeline {
	var reverseValues []pipeline.ConsumerPipeline

	for i := len(values) - 1; i >= 0; i-- {
		reverseValues = append(reverseValues, values[i])
	}

	return reverseValues
}

func (r *rabbitMQConsumer) existsPipeType(p reflect.Type) bool {
	for _, pipe := range r.pipelines {
		if reflect.TypeOf(pipe) == p {
			return true
		}
	}

	return false
}
