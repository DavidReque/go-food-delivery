package rabbitmq

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/consumer"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"
	rabbitmqConfigurations "github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/configurations"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/consumer/configurations"
	createProductExternalEventV1 "github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/features/creating_product/v1/events/integrationevents/externalevents"
	deleteProductExternalEventV1 "github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/features/deleting_products/v1/events/integration_events/external_events"
	updateProductExternalEventsV1 "github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/features/updating_products/v1/events/integration_events/external_events"

	"github.com/go-playground/validator/v10"
)

func ConfigProductsRabbitMQ(
	builder rabbitmqConfigurations.RabbitMQConfigurationBuilder,
	logger logger.Logger,
	validator *validator.Validate,
	tracer tracing.AppTracer,
) {
	// add custom message type mappings
	// utils.RegisterCustomMessageTypesToRegistrty(map[string]types.IMessage{"productCreatedV1": &creatingProductIntegration.ProductCreatedV1{}})

	// add create product consumer
	builder.
		AddConsumer(
			createProductExternalEventV1.ProductCreatedV1{},
			func(builder configurations.RabbitMQConsumerConfigurationBuilder) {
				builder.WithHandlers(
					func(handlersBuilder consumer.ConsumerHandlerConfigurationBuilder) {
						handlersBuilder.AddHandler(
							createProductExternalEventV1.NewProductCreatedConsumer(
								logger,
								validator,
								tracer,
							),
						)
					},
				)
			}).
		// add delete product consumer
		AddConsumer(
			deleteProductExternalEventV1.ProductDeletedV1{},
			func(builder configurations.RabbitMQConsumerConfigurationBuilder) {
				builder.WithHandlers(
					func(handlersBuilder consumer.ConsumerHandlerConfigurationBuilder) {
						handlersBuilder.AddHandler(
							deleteProductExternalEventV1.NewProductDeletedConsumer(
								logger,
								validator,
								tracer,
							),
						)
						deleteProductExternalEventV1.NewProductDeletedConsumer(
							logger,
							validator,
							tracer,
						)
					},
				)
			}).
		// add update product consumer
		AddConsumer(
			updateProductExternalEventsV1.ProductUpdatedV1{},
			func(builder configurations.RabbitMQConsumerConfigurationBuilder) {
				builder.WithHandlers(
					func(handlersBuilder consumer.ConsumerHandlerConfigurationBuilder) {
						handlersBuilder.AddHandler(
							updateProductExternalEventsV1.NewProductUpdatedConsumer(
								logger,
								validator,
								tracer,
							),
						)
						updateProductExternalEventsV1.NewProductUpdatedConsumer(
							logger,
							validator,
							tracer,
						)
					},
				)
			})
}
