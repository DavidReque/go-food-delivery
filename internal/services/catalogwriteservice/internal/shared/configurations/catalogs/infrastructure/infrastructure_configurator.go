package infrastructure

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/fxapp/contracts"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	loggingpipelines "github.com/DavidReque/go-food-delivery/internal/pkg/logger/pipelines"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/metrics"
	metricspipelines "github.com/DavidReque/go-food-delivery/internal/pkg/otel/metrics/mediatr/pipelines"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"
	tracingpipelines "github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing/mediatr/pipelines"
	postgrespipelines "github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/pipelines"
	validationpipeline "github.com/DavidReque/go-food-delivery/internal/pkg/validation/pipeline"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

type InfrastructureConfigurator struct {
	contracts.Application
}

func (ic *InfrastructureConfigurator) ConfigInfrastructures() {
	ic.ResolveFunc(
		func(l logger.Logger, tracer tracing.AppTracer, metrics metrics.AppMetrics, db *gorm.DB) error {
			return mediatr.RegisterRequestPipelineBehaviors(
				loggingpipelines.NewMediatorLoggingPipeline(l),
				validationpipeline.NewMediatorValidationPipeline(l),
				tracingpipelines.NewMediatorTracingPipeline(
					tracer,
					tracingpipelines.WithLogger(l),
				),
				metricspipelines.NewMediatorMetricsPipeline(
					metrics,
					metricspipelines.WithLogger(l),
				),
				postgrespipelines.NewMediatorTransactionPipeline(l, db),
			)
		},
	)
}
