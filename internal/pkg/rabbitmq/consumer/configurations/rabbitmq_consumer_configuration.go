package configurations

import (
	"fmt"
	"reflect"

	consumer2 "github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/consumer"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/pipeline"
	types2 "github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/types"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/utils"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/consumer/options"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/types"
)

type RabbitMQConsumerConfiguration struct {
	Name                string
	ConsumerMessageType reflect.Type
	Pipelines           []pipeline.ConsumerPipeline
	Handlers            []consumer2.ConsumerHandler
	*consumer2.ConsumerOptions
	ConcurrencyLimit int
	// The prefetch count tells the Rabbit connection how many messages to retrieve from the server per request.
	PrefetchCount   int // prefetch count is the number of messages to prefetch from the server
	AutoAck         bool // if true, the consumer will automatically acknowledge the message after it is processed
	NoLocal         bool // if true, the consumer will not receive messages published on the same connection
	NoWait          bool // if true, the consumer will not wait for the server to confirm the message
	BindingOptions  *options.RabbitMQBindingOptions
	QueueOptions    *options.RabbitMQQueueOptions
	ExchangeOptions *options.RabbitMQExchangeOptions
}

func NewDefaultRabbitMQConsumerConfiguration(
	messageType types2.IMessage,
) *RabbitMQConsumerConfiguration {
	name := fmt.Sprintf("%s_consumer", utils.GetMessageName(messageType))

	return &RabbitMQConsumerConfiguration{
		ConsumerOptions:  &consumer2.ConsumerOptions{ExitOnError: false, ConsumerId: ""},
		ConcurrencyLimit: 1,     // how many messages we can handle at once
		PrefetchCount:    4,     // how many messages we can handle at once
		NoLocal:          false, // if true, the consumer will not receive messages published on the same connection
		NoWait:           true,  // if true, the consumer will not wait for the server to confirm the message
		BindingOptions: &options.RabbitMQBindingOptions{
			RoutingKey: utils.GetRoutingKey(messageType),
		},
		ExchangeOptions: &options.RabbitMQExchangeOptions{
			Durable: true,
			Type:    types.ExchangeTopic,
			Name:    utils.GetTopicOrExchangeName(messageType),
		},
		QueueOptions: &options.RabbitMQQueueOptions{
			Durable: true,
			Name:    utils.GetQueueName(messageType),
		},
		ConsumerMessageType: utils.GetMessageBaseReflectType(messageType),
		Name:                name,
	}
}
