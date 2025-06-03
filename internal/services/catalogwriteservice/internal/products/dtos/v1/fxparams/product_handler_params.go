package fxparams

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/producer"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/shared/data/dbcontext"
	"go.uber.org/fx"
)

// ProductHandlerParams define las dependencias necesarias para los handlers de productos
// Este struct utiliza Uber FX para inyección de dependencias automática
// Todas las dependencias listadas aquí serán proporcionadas automáticamente por el contenedor FX
type ProductHandlerParams struct {
	// fx.In es una macro de Uber FX que define las dependencias inyectadas automáticamente
	fx.In

	// Logger es la instancia de logging para registrar eventos de los handlers
	Log logger.Logger

	// CatalogsDBContext es el contexto de la base de datos para operaciones de productos
	// Se utiliza para interactuar con la base de datos de productos
	CatalogsDBContext *dbcontext.CatalogsGormDBContext

	// RabbitmqProducer es el productor de mensajes para la cola de RabbitMQ
	// Se utiliza para publicar mensajes en la cola de RabbitMQ
	RabbitmqProducer producer.Producer

	Tracer tracing.AppTracer
}
