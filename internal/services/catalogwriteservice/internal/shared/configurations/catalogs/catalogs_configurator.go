package catalogs

import (
	"fmt"
	"net/http"

	//"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	"github.com/DavidReque/go-food-delivery/internal/pkg/fxapp/contracts"
	echocontracts "github.com/DavidReque/go-food-delivery/internal/pkg/http/customecho/contracts"
	migrationcontracts "github.com/DavidReque/go-food-delivery/internal/pkg/migration/contracts"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/config"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/configurations"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/shared/configurations/catalogs/infrastructure"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CatalogsServiceConfigurator struct {
	contracts.Application      // dependency injection
	infrastructureConfigurator *infrastructure.InfrastructureConfigurator
	productsModuleConfigurator *configurations.ProductsModuleConfigurator
}

func NewCatalogsServiceConfigurator(
	app contracts.Application,
) *CatalogsServiceConfigurator {
	infraConfigurator := infrastructure.NewInfrastructureConfigurator(app)
	productModuleConfigurator := configurations.NewProductsModuleConfigurator(
		app,
	)

	return &CatalogsServiceConfigurator{
		Application:                app,
		infrastructureConfigurator: infraConfigurator,
		productsModuleConfigurator: productModuleConfigurator,
	}
}

func (ic *CatalogsServiceConfigurator) ConfigureCatalogs() error {
	// shared
	// infrastructure
	ic.infrastructureConfigurator.ConfigInfrastructures()

	// shared
	// catalogs configurations
	// config database
	ic.ResolveFunc(
		func(db *gorm.DB, postgresMigrationRunner migrationcontracts.PostgresMigrationRunner) error {
			// migrate database
			err := ic.migrateCatalogs(postgresMigrationRunner)
			if err != nil {
				return err
			}

			// if we are not in test environment, seed the database
			// TEMPORARILY DISABLED: Seeding requires tables to exist (created by migrations)
			// if ic.Environment() != environment.Test {
			// 	err = ic.seedCatalogs(db)
			// 	if err != nil {
			// 		return err
			// 	}
			// }

			return nil
		},
	)

	// Modules
	// Product Module
	err := ic.productsModuleConfigurator.ConfigureProductsModule()

	return err
}

func (ic *CatalogsServiceConfigurator) MapCatalogsEndpoints() error {
	// Shared
	// config catalogs endpoints
	ic.ResolveFunc(
		func(catalogsServer echocontracts.EchoHttpServer, options *config.AppOptions) error {
			catalogsServer.SetupDefaultMiddlewares()

			// config catalogs root endpoint
			catalogsServer.RouteBuilder().
				RegisterRoutes(func(e *echo.Echo) {
					e.GET("", func(ec echo.Context) error {
						return ec.String(
							http.StatusOK,
							fmt.Sprintf(
								"%s is running...",
								options.GetMicroserviceNameUpper(),
							),
						)
					})
				})

			// config catalogs swagger
			//ic.configSwagger(catalogsServer.RouteBuilder())

			return nil
		},
	)

	// Modules
	// Products CatalogsServiceModule endpoints
	// map products endpoints
	err := ic.productsModuleConfigurator.MapProductsEndpoints()

	return err
}
