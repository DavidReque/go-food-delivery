package metrics

import (
	"context"
	"fmt"

	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	"github.com/DavidReque/go-food-delivery/internal/pkg/http/customecho/contracts"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// MetricsService maneja la recolección y exposición de métricas
type MetricsService struct {
	config      *MetricsOptions
	logger      logger.Logger
	environment environment.Environment
	provider    *sdkmetric.MeterProvider
	registry    *prometheus.Registry
}

func NewMetricsService(
	config *MetricsOptions,
	logger logger.Logger,
	env environment.Environment,
) (*MetricsService, error) {
	ms := &MetricsService{
		config:      config,
		logger:      logger,
		environment: env,
		registry:    prometheus.NewRegistry(),
	}

	if err := ms.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize metrics service: %w", err)
	}

	return ms, nil
}

// initialize configura los recursos y el proveedor de métricas
func (ms *MetricsService) initialize() error {
	// Crear el recurso con los atributos básicos
	res, err := ms.createResource()
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Configurar el proveedor de métricas
	ms.provider = sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithView(ms.defaultViews()...),
	)

	return nil
}

// createResource crea un recurso con los atributos básicos del servicio
func (ms *MetricsService) createResource() (*resource.Resource, error) {
	return resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(ms.config.ServiceName),
			semconv.ServiceVersion(ms.config.Version),
			attribute.String("environment", ms.environment.GetEnvironmentName()),
		),
		resource.WithHost(),
		resource.WithTelemetrySDK(),
	)
}

// defaultViews retorna las vistas por defecto para las métricas
func (ms *MetricsService) defaultViews() []sdkmetric.View {
	return []sdkmetric.View{
		// agregar vistas personalizadas para tus métricas
		// Por ejemplo, histogramas con buckets específicos
	}
}

// RegisterMetricsEndpoint registra el endpoint para exponer las métricas
func (ms *MetricsService) RegisterMetricsEndpoint(server contracts.EchoHttpServer) {
	metricsPath := ms.config.MetricsRoutePath
	if metricsPath == "" {
		metricsPath = "/metrics"
	}

	// Crear un handler de Prometheus que use nuestro registro personalizado
	handler := promhttp.HandlerFor(ms.registry, promhttp.HandlerOpts{})

	// Registrar el endpoint en Echo
	server.GetEchoInstance().GET(metricsPath, echo.WrapHandler(handler))
}

// CreateCounter crea un nuevo contador
func (ms *MetricsService) CreateCounter(name, description string) (metric.Int64Counter, error) {
	meter := ms.provider.Meter(ms.config.InstrumentationName)
	return meter.Int64Counter(
		name,
		metric.WithDescription(description),
	)
}

// CreateHistogram crea un nuevo histograma
func (ms *MetricsService) CreateHistogram(name, description string) (metric.Float64Histogram, error) {
	meter := ms.provider.Meter(ms.config.InstrumentationName)
	return meter.Float64Histogram(
		name,
		metric.WithDescription(description),
	)
}

// CreateUpDownCounter crea un nuevo contador que puede incrementar y decrementar
func (ms *MetricsService) CreateUpDownCounter(name, description string) (metric.Int64UpDownCounter, error) {
	meter := ms.provider.Meter(ms.config.InstrumentationName)
	return meter.Int64UpDownCounter(
		name,
		metric.WithDescription(description),
	)
}

// Shutdown limpia los recursos del servicio de métricas
func (ms *MetricsService) Shutdown(ctx context.Context) error {
	if ms.provider != nil {
		return ms.provider.Shutdown(ctx)
	}
	return nil
}

// Ejemplos de uso:

// RecordRequestDuration registra la duración de una petición HTTP
func (ms *MetricsService) RecordRequestDuration(path string, duration float64) {
	if histogram, err := ms.CreateHistogram(
		"http_request_duration_seconds",
		"Duración de las peticiones HTTP en segundos",
	); err == nil {
		histogram.Record(context.Background(), duration,
			metric.WithAttributes(attribute.String("path", path)),
		)
	}
}

// IncrementRequestCounter incrementa el contador de peticiones
func (ms *MetricsService) IncrementRequestCounter(path string, status int) {
	if counter, err := ms.CreateCounter(
		"http_requests_total",
		"Número total de peticiones HTTP",
	); err == nil {
		counter.Add(context.Background(), 1,
			metric.WithAttributes(
				attribute.String("path", path),
				attribute.Int("status", status),
			),
		)
	}
}
