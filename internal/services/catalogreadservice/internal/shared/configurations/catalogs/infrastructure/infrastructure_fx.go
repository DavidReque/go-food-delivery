package infrastructure

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core"
	"github.com/DavidReque/go-food-delivery/internal/pkg/grpc"
	"github.com/DavidReque/go-food-delivery/internal/pkg/health"
	customEcho "github.com/DavidReque/go-food-delivery/internal/pkg/http/customecho"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/mongodb"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/metrics"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/configurations"
	"github.com/DavidReque/go-food-delivery/internal/pkg/redis"
	rabbitmq2 "github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/configurations/rabbitmq"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

// https://pmihaylov.com/shared-components-go-microservices/
var Module = fx.Module(
	"infrastructurefx",
	// Modules
	core.Module,
	customEcho.Module,
	grpc.Module,
	mongodb.Module,
	redis.Module,
	rabbitmq.ModuleFunc(
		func(v *validator.Validate, l logger.Logger, tracer tracing.AppTracer) configurations.RabbitMQConfigurationBuilderFuc {
			return func(builder configurations.RabbitMQConfigurationBuilder) {
				rabbitmq2.ConfigProductsRabbitMQ(builder, l, v, tracer)
			}
		},
	),
	health.Module,
	tracing.Module,
	metrics.Module,

	// Other provides
	fx.Provide(validator.New),
)
