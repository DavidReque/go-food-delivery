package tracing

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"

	"go.opentelemetry.io/otel/trace"
)

var MessagingTracer trace.Tracer

func init() {
	MessagingTracer = tracing.NewAppTracer(
		"github.com/DavidReque/go-food-delivery/internal/pkg/messaging",
	) // instrumentation name
}
