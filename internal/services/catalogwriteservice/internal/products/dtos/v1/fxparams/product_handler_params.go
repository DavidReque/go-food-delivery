package fxparams

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"go.uber.org/fx"
)

// ProductHandlerParams define las dependencias necesarias para los handlers de productos
// Este struct utiliza Uber FX para inyección de dependencias automática
// Todas las dependencias listadas aquí serán proporcionadas automáticamente por el contenedor FX
type ProductHandlerParams struct {
	fx.In

	Log logger.Logger
}
