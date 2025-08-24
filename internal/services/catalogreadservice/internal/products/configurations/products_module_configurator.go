package configurations

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/web/route"
	"github.com/DavidReque/go-food-delivery/internal/pkg/fxapp/contracts"
	logger2 "github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/configurations/mappings"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/configurations/mediator"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/contracts/data"
)

type ProductsModuleConfigurator struct {
	contracts.Application
}

func NewProductsModuleConfigurator(
	app contracts.Application,
) *ProductsModuleConfigurator {
	return &ProductsModuleConfigurator{
		Application: app,
	}
}

func (c *ProductsModuleConfigurator) ConfigureProductsModule() {
	c.ResolveFunc(
		func(logger logger2.Logger, mongoRepository data.ProductRepository, cacheRepository data.ProductCacheRepository, tracer tracing.AppTracer) error {
			// config Products Mediators
			err := mediator.ConfigProductsMediator(
				logger,
				mongoRepository,
				cacheRepository,
				tracer,
			)
			if err != nil {
				return err
			}

			// config Products Mappings
			err = mappings.ConfigureProductsMappings()
			if err != nil {
				return err
			}
			return nil
		},
	)
}

func (c *ProductsModuleConfigurator) MapProductsEndpoints() {
	// config Products Http Endpoints
	c.ResolveFuncWithParamTag(func(endpoints []route.Endpoint) {
		for _, endpoint := range endpoints {
			endpoint.MapEndpoint()
		}
	}, `group:"product-routes"`,
	)
}
