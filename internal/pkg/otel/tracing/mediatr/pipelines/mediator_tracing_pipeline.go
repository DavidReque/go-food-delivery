package pipelines

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"
	"github.com/DavidReque/go-food-delivery/internal/pkg/reflection/typemapper"
	"github.com/mehdihadeli/go-mediatr"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// mediatorTracingPipeline es un comportamiento de pipeline de MediatR que añade
// capacidades de tracing distribuido a cada solicitud manejada.
type mediatorTracingPipeline struct {
	config *config           // Almacena la configuración del pipeline (ej. serviceName, logger).
	tracer tracing.AppTracer // El tracer de OpenTelemetry para crear spans.
}

// NewMediatorTracingPipeline crea una nueva instancia del pipeline de tracing.
// Sigue el patrón de Opciones Funcionales, permitiendo una configuración flexible.
// Parámetros:
//   - appTracer: La instancia del tracer de la aplicación.
//   - opts: Una lista variable de opciones de configuración (ej. WithServiceName).
//
// Retorna:
//   - Un mediatr.PipelineBehavior que puede ser registrado en MediatR.
func NewMediatorTracingPipeline(
	appTracer tracing.AppTracer,
	opts ...Option,
) mediatr.PipelineBehavior {
	// Comienza con la configuración por defecto.
	cfg := defaultConfig
	// Aplica todas las opciones proporcionadas por el usuario para sobreescribir los valores por defecto.
	for _, opt := range opts {
		opt.apply(cfg)
	}

	// Devuelve una nueva instancia del pipeline con la configuración final.
	return &mediatorTracingPipeline{
		config: cfg,
		tracer: appTracer,
	}
}

// Handle es el método que intercepta cada solicitud que pasa por MediatR.
// Implementa la lógica principal del tracing para la solicitud.
func (r *mediatorTracingPipeline) Handle(
	ctx context.Context, // El contexto de la solicitud, que puede contener un span padre.
	request interface{}, // La solicitud (comando, query o evento) a ser procesada.
	next mediatr.RequestHandlerFunc, // La función que invoca al siguiente manejador en la cadena.
) (interface{}, error) {
	// Obtiene el nombre de la solicitud en formato snake_case (ej. 'create_order_command').
	requestName := typemapper.GetSnakeTypeName(request)
	// Obtiene el nombre del paquete para determinar si es un comando, query o evento.
	packageName := typemapper.GetPackageName(request)

	// Define un conjunto de nombres y etiquetas por defecto para el tracing.
	componentName := "RequestHandler"
	requestNameTag := "app.request_name"
	requestTag := "app.request"
	requestResultNameTag := "app.request_result_name"
	requestResultTag := "app.request_result"

	// Personaliza las etiquetas y el nombre del componente basándose en el tipo de solicitud.
	// Esto enriquece los traces, haciendo más fácil filtrar por tipo de operación.
	switch {
	case strings.Contains(packageName, "command"):
		componentName = "CommandHandler"
		requestNameTag = "app.command_name"
		requestTag = "app.command"
		requestResultNameTag = "app.command_result_name"
		requestResultTag = "app.command_result"
	case strings.Contains(packageName, "query"):
		componentName = "QueryHandler"
		requestNameTag = "app.query_name"
		requestTag = "app.query"
		requestResultNameTag = "app.query_result_name"
		requestResultTag = "app.query_result"
	case strings.Contains(packageName, "event"):
		componentName = "EventHandler"
		requestNameTag = "app.event_name"
		requestTag = "app.event"
		requestResultNameTag = "app.event_result_name"
		requestResultTag = "app.event_result"
	}

	// Construye un nombre de operación descriptivo para el span.
	operationName := fmt.Sprintf("%s_handler", requestName)
	// Construye el nombre del span siguiendo una convención estandarizada.
	// Ej: "CommandHandler.create_order_command_handler/create_order_command"
	spanName := fmt.Sprintf("%s.%s/%s", componentName, operationName, requestName)

	// Inicia un nuevo span de tracing. El `newCtx` contiene la información del nuevo span
	// y será propagado al siguiente manejador.
	newCtx, span := r.tracer.Start(ctx, spanName)
	// `defer span.End()` asegura que el span se cierre siempre al final de la función,
	// calculando su duración total.
	defer span.End()

	// Establece los atributos (tags) iniciales en el span.
	// Incluye el nombre de la solicitud y el contenido de la solicitud como un JSON.
	span.SetAttributes(
		attribute.String(requestNameTag, requestName),
		jsonObject(requestTag, request),
	)

	// Invoca al siguiente manejador en el pipeline, pasando el nuevo contexto con el span activo.
	response, err := next(newCtx)

	// Una vez que el manejador ha terminado, se añade la información de la respuesta al span.
	responseName := typemapper.GetSnakeTypeName(response)
	span.SetAttributes(
		attribute.String(requestResultNameTag, responseName),
		jsonObject(requestResultTag, response),
	)

	// Si ocurrió un error durante el procesamiento, se registra en el span.
	if err != nil {
		span.RecordError(err)                    // Registra el error detallado.
		span.SetStatus(codes.Error, err.Error()) // Marca el status del span como erróneo.
	}

	// Devuelve la respuesta y el error al llamador anterior en la cadena.
	return response, err
}

// jsonObject es una función de utilidad que serializa un objeto a una cadena JSON.
// Esto permite registrar objetos complejos como atributos de texto en un span.
// Parámetros:
//   - key: El nombre del atributo (tag) para el span.
//   - value: El objeto a ser serializado.
//
// Retorna:
//   - Un attribute.KeyValue que puede ser usado con `span.SetAttributes`.
func jsonObject(key string, value interface{}) attribute.KeyValue {
	if value == nil {
		return attribute.String(key, "")
	}
	// Intenta serializar el valor a JSON.
	bytes, err := json.Marshal(value)
	if err != nil {
		// Si falla la serialización, registra el error como valor del atributo.
		return attribute.String(key, fmt.Sprintf("Error marshaling to JSON: %v", err))
	}
	// Devuelve el atributo con el JSON como valor.
	return attribute.String(key, string(bytes))
}
