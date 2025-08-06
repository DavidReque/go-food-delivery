package consumer

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/consumer"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/types"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/serializer"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/config"
	consumerConfigurations "github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/consumer/configurations"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/consumer/consumercontracts"
	types2 "github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/types"
)

type consumerFactory struct {
	connection      types2.IConnection           // the connection to the RabbitMQ server
	eventSerializer serializer.MessageSerializer // the serializer for the event
	logger          logger.Logger
	rabbitmqOptions *config.RabbitmqOptions
}

func NewConsumerFactory(
	rabbitmqOptions *config.RabbitmqOptions,
	connection types2.IConnection,
	eventSerializer serializer.MessageSerializer,
	logger logger.Logger,
) consumercontracts.ConsumerFactory {
	return &consumerFactory{
		rabbitmqOptions: rabbitmqOptions,
		logger:          logger,
		eventSerializer: eventSerializer,
		connection:      connection,
	}
}

func (c *consumerFactory) CreateConsumer(
	consumerConfiguration *consumerConfigurations.RabbitMQConsumerConfiguration,
	isConsumedNotifications ...func(message types.IMessage),
) (consumer.Consumer, error) {
	return NewRabbitMQConsumer(
		c.rabbitmqOptions,
		c.connection,
		consumerConfiguration,
		c.eventSerializer,
		c.logger,
		isConsumedNotifications...,
	)
}

func (c *consumerFactory) Connection() types2.IConnection {
	return c.connection
}
