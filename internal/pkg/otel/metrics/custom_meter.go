package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type AppMetrics interface {
	metric.Meter
}

type appMetrics struct {
	metric.Meter
}

func NewAppMeter(name string, opts ...metric.MeterOption) AppMetrics {
	meter := otel.Meter(name, opts...)
	return &appMetrics{Meter: meter}
}
