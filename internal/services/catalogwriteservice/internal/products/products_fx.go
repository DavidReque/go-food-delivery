package products

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/cqrs"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/web/route"
	"github.com/DavidReque/go-food-delivery/internal/pkg/http/customecho/contracts"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/data/repositories"
	creatingproductv1 "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/shared/grpc"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// Module is the module for the products service.
var Module = fx.Module(
	"productsfx",

	// Other providers
	fx.Provide(repositories.NewPostgresProductRepository),
	fx.Provide(grpc.NewProductGrpcService),

	// Provee un grupo de rutas para el módulo de productos utilizando fx (framework de inyección de dependencias)
	fx.Provide(
		// Anotamos la función proveedora para poder identificar el grupo de rutas resultante
		fx.Annotate(
			// Esta función recibe el servidor HTTP de catálogos como dependencia
			// y retorna un grupo de Echo para las rutas de productos
			func(catalogsServer contracts.EchoHttpServer) *echo.Group {
				// Variable para almacenar y retornar el grupo de rutas
				var g *echo.Group

				// Configuramos las rutas usando el RouteBuilder
				catalogsServer.RouteBuilder().
					// Creamos un grupo base para la versión 1 de la API
					RegisterGroupFunc("/api/v1", func(v1 *echo.Group) {
						// Dentro del grupo v1, creamos un subgrupo específico para productos
						// Todas las rutas aquí tendrán el prefijo /api/v1/products
						group := v1.Group("/products")
						g = group
					})

				return g
			},
			// Etiquetamos el grupo resultante para que pueda ser inyectado en otros componentes
			// usando el nombre "product-echo-group"
			fx.ResultTags(`name:"product-echo-group"`)),
	),

	// add cqrs handlers to DI
	fx.Provide(
		cqrs.AsHandler(
			creatingproductv1.NewCreateProductHandler,
			"product-handlers",
		),
	),

	// add endpoints to DI
	fx.Provide(
		route.AsRoute(
			creatingproductv1.NewCreteProductEndpoint,
			"product-routes",
		),
	),
)
