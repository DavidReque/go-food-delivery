package utils

import (
	"context"

	"go.opentelemetry.io/otel/codes"
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

// TraceStatusFromContext establece el estado del trace basado en el error
// Si hay un error, establece el estado como ERROR y registra el error
// Si no hay error, establece el estado como OK
func TraceStatusFromContext(ctx context.Context, err error) error {
	span := trace.SpanFromContext(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// GetParentSpanFromContext obtiene el span padre del contexto
func GetParentSpanFromContext(ctx context.Context) trace.Span {
	if span, ok := ctx.Value(parentSpanKey).(trace.Span); ok {
		return span
	}
	return nil
}

// TraceStatusFromSpan establece el estado del trace basado en el error
// Si hay un error, establece el estado como ERROR y registra el error
// Si no hay error, establece el estado como OK
// A diferencia de TraceStatusFromContext, esta función recibe el span directamente
func TraceStatusFromSpan(span trace.Span, err error) error {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}

	span.SetStatus(codes.Ok, "")
	return nil
}
