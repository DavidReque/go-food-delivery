package utils

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/constants/telemetrytags"
	errorUtils "github.com/DavidReque/go-food-delivery/internal/pkg/utils/errorutils"

	"github.com/ahmetb/go-linq/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	trace2 "go.opentelemetry.io/otel/sdk/trace"
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

func GetSpanAttributeFromCurrentContext(
	ctx context.Context,
	attributeName string,
) attribute.KeyValue {
	span := trace.SpanFromContext(ctx)
	readWriteSpan, ok := span.(trace2.ReadWriteSpan)
	if !ok {
		return *new(attribute.KeyValue)
	}
	att := linq.From(readWriteSpan.Attributes()).
		FirstWithT(func(att attribute.KeyValue) bool { return string(att.Key) == attributeName })

	return att.(attribute.KeyValue)
}

func GetSpanAttribute(
	span trace.Span,
	attributeName string,
) attribute.KeyValue {
	readWriteSpan, ok := span.(trace2.ReadWriteSpan)
	if !ok {
		return *new(attribute.KeyValue)
	}

	att := linq.From(readWriteSpan.Attributes()).
		FirstWithT(func(att attribute.KeyValue) bool { return string(att.Key) == attributeName })

	return att.(attribute.KeyValue)
}

func TraceErrStatusFromSpan(span trace.Span, err error) error {
	isError := err != nil

	span.SetStatus(codes.Error, err.Error())

	if isError {
		stackTraceError := errorUtils.ErrorsWithStack(err)

		// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
		span.SetAttributes(
			attribute.String(telemetrytags.Exceptions.Message, err.Error()),
			attribute.String(telemetrytags.Exceptions.Stacktrace, stackTraceError),
		)
		span.RecordError(err)
	}

	return err
}

func TraceErrStatusFromContext(ctx context.Context, err error) error {
	// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
	span := trace.SpanFromContext(ctx)

	defer span.End()

	return TraceErrStatusFromSpan(span, err)
}
