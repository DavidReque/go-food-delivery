package infrastructure

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/fxapp/contracts"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	loggingpipelines "github.com/DavidReque/go-food-delivery/internal/pkg/logger/pipelines"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/metrics"
	metricspipelines "github.com/DavidReque/go-food-delivery/internal/pkg/otel/metrics/mediatr/pipelines"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"
	tracingpipelines "github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing/mediatr/pipelines"

	"github.com/mehdihadeli/go-mediatr"
)

type InfrastructureConfigurator struct {
	contracts.Application
}

func NewInfrastructureConfigurator(
	app contracts.Application,
) *InfrastructureConfigurator {
	return &InfrastructureConfigurator{
		Application: app,
	}
}

func (ic *InfrastructureConfigurator) ConfigInfrastructures() {
	ic.ResolveFunc(
		func(l logger.Logger, tracer tracing.AppTracer, metrics metrics.AppMetrics) error {
			err := mediatr.RegisterRequestPipelineBehaviors(
				loggingpipelines.NewMediatorLoggingPipeline(l),
				tracingpipelines.NewMediatorTracingPipeline(
					tracer,
					tracingpipelines.WithLogger(l),
				),
				metricspipelines.NewMediatorMetricsPipeline(
					metrics,
					metricspipelines.WithLogger(l),
				),
			)

			return err
		},
	)
}
