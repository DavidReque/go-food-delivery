package configurations

import (
	consumerConfigurations "github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/consumer/configurations"
	producerConfigurations "github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/producer/configurations"
)

type RabbitMQConfiguration struct {
	ProducersConfigurations []*producerConfigurations.RabbitMQProducerConfiguration
	ConsumersConfigurations []*consumerConfigurations.RabbitMQConsumerConfiguration
}
