package configurations

import (
	fxcontracts "github.com/DavidReque/go-food-delivery/internal/pkg/fxapp/contracts"
	grpcServer "github.com/DavidReque/go-food-delivery/internal/pkg/grpc"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/configurations/endpoints"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/configurations/mappings"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/configurations/mediator"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/shared/grpc"
	productsservice "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/shared/grpc/genproto"

	googleGrpc "google.golang.org/grpc"
)

type ProductsModuleConfigurator struct {
	fxcontracts.Application
}

func NewProductsModuleConfigurator(
	fxapp fxcontracts.Application,
) *ProductsModuleConfigurator {
	return &ProductsModuleConfigurator{
		Application: fxapp,
	}
}

func (c *ProductsModuleConfigurator) ConfigureProductsModule() error {
	// config products mappings
	err := mappings.ConfigureProductsMappings()
	if err != nil {
		return err
	}

	// register products request handler on mediator
	c.ResolveFuncWithParamTag(
		mediator.RegisterMediatorHandlers,
		`group:"product-handlers"`,
	)

	return nil
}

// MapProductsEndpoints configures the endpoints for the products module.
// It registers the endpoints for the products module and configures the grpc endpoints.
func (c *ProductsModuleConfigurator) MapProductsEndpoints() error {
	// config endpoints
	c.ResolveFuncWithParamTag(
		endpoints.RegisterEndpoints,
		`group:"product-routes"`,
	)

	// config Products Grpc Endpoints
	c.ResolveFunc(
		func(catalogsGrpcServer grpcServer.GrpcServer, grpcService *grpc.ProductGrpcServiceServer) error {
			catalogsGrpcServer.GrpcServiceBuilder().
				RegisterRoutes(func(server *googleGrpc.Server) {
					productsservice.RegisterProductsServiceServer(
						server,
						grpcService,
					)
				})

			return nil
		},
	)

	return nil
}
