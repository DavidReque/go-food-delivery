package fxparams

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/shared/contracts"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// ProductRouteParams define todas las dependencias necesarias para configurar las rutas de productos
// Este struct utiliza Uber FX para inyección automática de dependencias
// Contiene todas las herramientas necesarias para manejar requests HTTP de productos
type ProductRouteParams struct {
	fx.In

	// CatalogsMetrics proporciona contadores y métricas para observabilidad
	// Se utiliza para rastrear requests, errores, latencias y eventos de productos
	// Ejemplo: params.CatalogsMetrics.CreateProductGrpcRequests.Add(ctx, 1)
	CatalogsMetrics *contracts.CatalogsMetrics

	// Logger es la instancia de logging para registrar eventos de las rutas
	// Permite logging estructurado con contexto adicional para debugging y monitoring

	Logger logger.Logger

	// ProductsGroup es el grupo de rutas específico para endpoints de productos
	// Se utiliza para registrar todas las rutas bajo el prefijo "/products"

	ProductsGroup *echo.Group `name:"product-echo-group"`

	// Validator proporciona validación de datos de entrada para requests HTTP
	// Se utiliza para validar structs de request antes de procesarlos
	Validator *validator.Validate
}
