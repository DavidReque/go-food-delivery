package pipelines

import (
	"context"
	"fmt"
	"time"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/cqrs"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/events"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/constants/telemetrytags"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/metrics"
	customAttribute "github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing/attribute"
	"github.com/DavidReque/go-food-delivery/internal/pkg/reflection/typemapper"

	"go.opentelemetry.io/otel/metric"

	"github.com/mehdihadeli/go-mediatr"
	"go.opentelemetry.io/otel/attribute"
)

type mediatorMetricsPipeline struct {
	config *config            // config para pipeline
	meter  metrics.AppMetrics // medidor de métricas
}

func NewMediatorMetricsPipeline(
	appMetrics metrics.AppMetrics,
	opts ...Option,
) mediatr.PipelineBehavior {
	cfg := defaultConfig
	for _, opt := range opts {
		opt.apply(cfg)
	}

	return &mediatorMetricsPipeline{
		config: cfg,
		meter:  appMetrics,
	}
}

func (r *mediatorMetricsPipeline) Handle(
	ctx context.Context, // contexto de la solicitud
	request interface{}, // solicitud a procesar
	next mediatr.RequestHandlerFunc, // siguiente handler
) (interface{}, error) {
	payloadSnakeTypeName := typemapper.GetSnakeTypeName(request)
	typeName := typemapper.GetTypeName(request)

	nameTag := telemetrytags.App.RequestName
	typeNameTag := telemetrytags.App.RequestType
	payloadTag := telemetrytags.App.Command
	resultSnakeTypeNameTag := telemetrytags.App.RequestResultName
	resultTag := telemetrytags.App.RequestResult

	if cqrs.IsCommand(request) {
		nameTag = telemetrytags.App.CommandName
		typeNameTag = telemetrytags.App.CommandType
		payloadTag = telemetrytags.App.Command
		resultSnakeTypeNameTag = telemetrytags.App.CommandResultName
		resultTag = telemetrytags.App.CommandResult
	} else if cqrs.IsQuery(request) {
		nameTag = telemetrytags.App.QueryName
		typeNameTag = telemetrytags.App.QueryType
		payloadTag = telemetrytags.App.Query
		resultSnakeTypeNameTag = telemetrytags.App.QueryResultName
		resultTag = telemetrytags.App.QueryResult
	} else if events.IsEvent(request) {
		nameTag = telemetrytags.App.EventName
		typeNameTag = telemetrytags.App.EventType
		payloadTag = telemetrytags.App.Event
		resultSnakeTypeNameTag = telemetrytags.App.EventResultName
		resultTag = telemetrytags.App.EventResult
	}

	successRequestCounter, err := r.meter.Int64Counter(
		fmt.Sprintf("%s.failed_total", payloadSnakeTypeName),
		metric.WithUnit("count"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the number of failed '%s' (%s)",
				payloadSnakeTypeName,
				typeName,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	totalRequestsCounter, err := r.meter.Int64Counter(
		fmt.Sprintf("%s.total", payloadSnakeTypeName),
		metric.WithUnit("count"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the total number of '%s' (%s)",
				payloadSnakeTypeName,
				typeName,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	failedRequestsCounter, err := r.meter.Int64Counter(
		fmt.Sprintf("%s.failed_total", payloadSnakeTypeName),
		metric.WithUnit("count"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the number of failed '%s' (%s)",
				payloadSnakeTypeName,
				typeName,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	durationValueRecorder, err := r.meter.Int64Histogram(
		fmt.Sprintf("%s.duration", payloadSnakeTypeName),
		metric.WithUnit("ms"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the duration of '%s' (%s)",
				payloadSnakeTypeName,
				typeName,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	// start recording the duration
	startTime := time.Now()

	response, err := next(ctx)

	// calculate the duration
	duration := time.Since(startTime).Microseconds()

	// response will be nil if there is an error
	responseSnakeName := typemapper.GetSnakeTypeName(response)

	opt := metric.WithAttributes(
		attribute.String(nameTag, payloadSnakeTypeName),
		attribute.String(typeNameTag, typeName),
		customAttribute.Object(payloadTag, request),
		attribute.String(resultSnakeTypeNameTag, responseSnakeName),
		customAttribute.Object(resultTag, response),
	)

	// record metrics
	totalRequestsCounter.Add(ctx, 1, opt)

	if err != nil {
		successRequestCounter.Add(ctx, 0, opt)
	} else {
		failedRequestsCounter.Add(ctx, 1, opt)
	}

	durationValueRecorder.Record(ctx, duration, opt)

	return response, err

}
