package tracing_headers

import "github.com/DavidReque/go-food-delivery/internal/pkg/core/metadata"

func GetTracingTraceId(m metadata.Metadata) string {
	return m.GetString(TraceId)
}

func GetTracingParentSpanId(m metadata.Metadata) string {
	return m.GetString(ParentSpanId)
}

func GetTracingTraceparent(m metadata.Metadata) string {
	return m.GetString(Traceparent)
}