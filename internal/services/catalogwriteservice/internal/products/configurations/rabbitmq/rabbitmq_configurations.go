package rabbitmq

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/configurations"
	producerConfigurations "github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/producer/configurations"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/events/integrationevents"
)

// ConfigProductsRabbitMQ configures the rabbitmq for the products
func ConfigProductsRabbitMQ(
	builder configurations.RabbitMQConfigurationBuilder,
) {
	// Add producer for the product created event
	builder.AddProducer(
		integrationevents.ProductCreatedV1{},
		func(builder producerConfigurations.RabbitMQProducerConfigurationBuilder) {
		},
	)
}
