package tracing

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// AppTracer es una interfaz que define los métodos para iniciar y finalizar spans en la aplicación
type AppTracer interface {
	trace.Tracer
}

// appTracer es una implementación de AppTracer que utiliza un tracer de OpenTelemetry
type appTracer struct {
	trace.Tracer
}

// Start inicia un nuevo span en la aplicación
func (c *appTracer) Start(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	parentSpan := trace.SpanFromContext(ctx)
	if parentSpan != nil {
		utils.ContextWithParentSpan(ctx, parentSpan)
	}

	return c.Tracer.Start(ctx, spanName, opts...)
}

// NewAppTracer crea una nueva instancia de AppTracer
func NewAppTracer(name string, options ...trace.TracerOption) AppTracer {
	tracer := otel.Tracer(name, options...)
	return &appTracer{Tracer: tracer}
}
