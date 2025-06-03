package utils

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// traceContextKeyType es una clave para almacenar el span padre en el contexto
type traceContextKeyType int

// parentSpanKey es la clave para almacenar el span padre en el contexto
const parentSpanKey traceContextKeyType = iota + 1

// ContextWithParentSpan agrega el span padre al contexto
func ContextWithParentSpan(
	parent context.Context,
	span trace.Span,
) context.Context {
	return context.WithValue(parent, parentSpanKey, span)
}
