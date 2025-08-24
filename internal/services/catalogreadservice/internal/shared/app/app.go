package app

import "github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/shared/configurations/catalogs"

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	// configure dependencies
	appBuilder := NewCatalogsReadApplicationBuilder()
	appBuilder.ProvideModule(catalogs.CatalogsServiceModule)

	app := appBuilder.Build()

	// configure application
	app.ConfigureCatalogs()

	// map endpoints
	app.MapCatalogsEndpoints()

	app.Logger().Info("Starting catalog_service application")
	app.Run()
}
