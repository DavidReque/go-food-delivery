package rabbitmq

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/configurations"
	producerConfigurations "github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/producer/configurations"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/types"
	creatingproductevents "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/events/integrationevents"
	deletingproductevents "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/deletingproduct/v1/events/integrationevents"
	updatingproductevents "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/updatingproduct/v1/events/integrationevents"
)

// ConfigProductsRabbitMQ configures the rabbitmq for the products
func ConfigProductsRabbitMQ(
	builder configurations.RabbitMQConfigurationBuilder,
) {
	// Add producer for the product created event
	builder.AddProducer(
		creatingproductevents.ProductCreatedV1{},
		func(builder producerConfigurations.RabbitMQProducerConfigurationBuilder) {
			builder.WithExchangeName("catalog.products.exchange"). // Exchange name
				WithExchangeType(types.ExchangeTopic). // Exchange type
				WithRoutingKey("products.created"). // Routing key
				WithDurable(true) // Durable
		},
	)

	// Add producer for the product updated event
	builder.AddProducer(
		updatingproductevents.ProductUpdatedV1{},
		func(builder producerConfigurations.RabbitMQProducerConfigurationBuilder) {
			builder.WithExchangeName("catalog.products.exchange"). // Exchange name
				WithExchangeType(types.ExchangeTopic). // Exchange type
				WithRoutingKey("products.updated"). // Routing key
				WithDurable(true) // Durable
		},
	)

	// Add producer for the product deleted event
	builder.AddProducer(
		deletingproductevents.ProductDeletedV1{},
		func(builder producerConfigurations.RabbitMQProducerConfigurationBuilder) {
			builder.WithExchangeName("catalog.products.exchange"). // Exchange name
				WithExchangeType(types.ExchangeTopic). // Exchange type
				WithRoutingKey("products.deleted"). // Routing key
				WithDurable(true) // Durable
		},
	)
}
