package metrics

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/http/customecho/contracts"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// OtelMetrics is the main metrics service
type OtelMetrics struct {
	config     *MetricsOptions
	appMetrics AppMetrics
}

// NewOtelMetrics creates a new instance of OtelMetrics
func NewOtelMetrics(config *MetricsOptions) (*OtelMetrics, error) {
	appMetrics := NewAppMeter(config.ServiceName)

	return &OtelMetrics{
		config:     config,
		appMetrics: appMetrics,
	}, nil
}

// RegisterMetricsEndpoint registers the endpoint to expose metrics
func (m *OtelMetrics) RegisterMetricsEndpoint(server contracts.EchoHttpServer) {
	metricsPath := m.config.MetricsRoutePath
	if metricsPath == "" {
		metricsPath = "/metrics"
	}

	// Register Prometheus endpoint
	server.GetEchoInstance().GET(metricsPath, echo.WrapHandler(promhttp.Handler()))
}

// Shutdown closes the metrics service
func (m *OtelMetrics) Shutdown(ctx context.Context) error {
	// We don't need to implement anything here
	// as the metrics provider is handled globally
	return nil
} 