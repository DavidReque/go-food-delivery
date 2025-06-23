package metrics

import (
	"context"
	"fmt"

	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	"github.com/DavidReque/go-food-delivery/internal/pkg/http/customecho/contracts"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/labstack/echo/v4"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	prometheusexporter "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// SimpleMetrics es una implementación simplificada del sistema de métricas
// que se centra en las funcionalidades más comunes y Prometheus como exportador
type SimpleMetrics struct {
	config      *MetricsOptions
	logger      logger.Logger
	environment environment.Environment
	provider    *sdkmetric.MeterProvider
	registry    *prom.Registry
}

// NewSimpleMetrics crea una nueva instancia del sistema de métricas simplificado
func NewSimpleMetrics(
	config *MetricsOptions,
	logger logger.Logger,
	env environment.Environment,
) (*SimpleMetrics, error) {
	sm := &SimpleMetrics{
		config:      config,
		logger:      logger,
		environment: env,
		registry:    prom.NewRegistry(),
	}

	if err := sm.initialize(); err != nil {
		return nil, fmt.Errorf("error inicializando métricas: %w", err)
	}

	return sm, nil
}

// initialize configura el sistema de métricas
func (sm *SimpleMetrics) initialize() error {
	// Crear el recurso con información básica del servicio
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(sm.config.ServiceName),
			semconv.ServiceVersion(sm.config.Version),
			attribute.String("environment", sm.environment.GetEnvironmentName()),
		),
	)
	if err != nil {
		return fmt.Errorf("error creando recurso: %w", err)
	}

	exporter, err := prometheusexporter.New(prometheusexporter.WithRegisterer(sm.registry))
	if err != nil {
		return fmt.Errorf("error creando exportador prometheus: %w", err)
	}

	sm.provider = sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(exporter),
	)

	return nil
}

// RegisterEndpoint configura el endpoint para exponer las métricas
func (sm *SimpleMetrics) RegisterEndpoint(server contracts.EchoHttpServer) {
	path := sm.config.MetricsRoutePath
	if path == "" {
		path = "/metrics"
	}

	handler := promhttp.HandlerFor(sm.registry, promhttp.HandlerOpts{})
	server.GetEchoInstance().GET(path, echo.WrapHandler(handler))
}

// Métricas comunes pre-configuradas

// RequestCounter crea y registra un contador para peticiones HTTP
func (sm *SimpleMetrics) RequestCounter() (metric.Int64Counter, error) {
	meter := sm.provider.Meter(sm.config.ServiceName)
	return meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total de peticiones HTTP procesadas"),
	)
}

// RequestDuration crea y registra un histograma para la duración de peticiones
func (sm *SimpleMetrics) RequestDuration() (metric.Float64Histogram, error) {
	meter := sm.provider.Meter(sm.config.ServiceName)
	return meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("Duración de las peticiones HTTP en segundos"),
	)
}

// ActiveConnections crea y registra un contador up/down para conexiones activas
func (sm *SimpleMetrics) ActiveConnections() (metric.Int64UpDownCounter, error) {
	meter := sm.provider.Meter(sm.config.ServiceName)
	return meter.Int64UpDownCounter(
		"active_connections",
		metric.WithDescription("Número actual de conexiones activas"),
	)
}

// Helpers para registrar métricas comunes

// RecordRequest registra una petición HTTP con su duración y estado
func (sm *SimpleMetrics) RecordRequest(path string, duration float64, status int) {
	ctx := context.Background()
	attrs := []attribute.KeyValue{
		attribute.String("path", path),
		attribute.Int("status", status),
	}

	// Registrar el contador de peticiones
	if counter, err := sm.RequestCounter(); err == nil {
		counter.Add(ctx, 1, metric.WithAttributes(attrs...))
	}

	// Registrar la duración
	if histogram, err := sm.RequestDuration(); err == nil {
		histogram.Record(ctx, duration, metric.WithAttributes(attrs...))
	}
}

// UpdateActiveConnections actualiza el contador de conexiones activas
func (sm *SimpleMetrics) UpdateActiveConnections(delta int64) {
	if counter, err := sm.ActiveConnections(); err == nil {
		counter.Add(context.Background(), delta)
	}
}

// Shutdown limpia los recursos del sistema de métricas
func (sm *SimpleMetrics) Shutdown(ctx context.Context) error {
	if sm.provider != nil {
		return sm.provider.Shutdown(ctx)
	}
	return nil
}
