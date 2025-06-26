package metrics

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/http/customecho/contracts"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"

	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
)

var (
	// Module proporcionado a fx
	Module = fx.Module(
		"otelmetrixfx",
		metricsProviders,
		metricsInvokes,
	)

	metricsProviders = fx.Options(fx.Provide(
		ProvideMetricsConfig,
		NewOtelMetrics,
		fx.Annotate(
			provideMeter,
			fx.ParamTags(`optional:"true"`),
			fx.As(new(AppMetrics)),
			fx.As(new(metric.Meter))),
	))

	metricsInvokes = fx.Options(
		fx.Invoke(registerHooks),
		fx.Invoke(func(m *OtelMetrics, server contracts.EchoHttpServer) {
			m.RegisterMetricsEndpoint(server)
		}),
	)
)

func provideMeter(otelMetrics *OtelMetrics) AppMetrics {
	return otelMetrics.appMetrics
}

// registerHooks registra los hooks del ciclo de vida para las métricas
func registerHooks(
	lc fx.Lifecycle,
	metrics *OtelMetrics,
	logger logger.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if metrics.appMetrics == nil {
				return nil
			}

			if metrics.config.EnableHostMetrics {
				logger.Info("Starting host instrumentation:")
				err := host.Start(
					host.WithMeterProvider(otel.GetMeterProvider()),
				)
				if err != nil {
					logger.Errorf(
						"error starting host instrumentation: %s",
						err,
					)
				}
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if metrics.appMetrics == nil {
				return nil
			}

			if err := metrics.Shutdown(ctx); err != nil {
				logger.Errorf(
					"error shutting down metrics provider: %v",
					err,
				)
			} else {
				logger.Info("metrics provider shutdown gracefully")
			}

			return nil
		},
	})
}
