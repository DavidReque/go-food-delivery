package infrastructure

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core"
	"github.com/DavidReque/go-food-delivery/internal/pkg/grpc"
	"github.com/DavidReque/go-food-delivery/internal/pkg/migration/goose"
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresmessaging"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/configurations"
	rabbitmq2 "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/configurations/rabbitmq"

	"github.com/DavidReque/go-food-delivery/internal/pkg/http/customecho"
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm"

	"go.uber.org/fx"
)

// https://pmihaylov.com/shared-components-go-microservices/


var Module = fx.Module(
	"infrastructurefx",
	// Modules
	core.Module,
	customecho.Module,
	grpc.Module,
	postgresgorm.Module,
	postgresmessaging.Module,
	goose.Module,
	rabbitmq.ModuleFunc(
		func() configurations.RabbitMQConfigurationBuilderFuc {
			return func(builder configurations.RabbitMQConfigurationBuilder) {
				rabbitmq2.ConfigProductsRabbitMQ(builder)
			}
		},
	),
)