package main

import (
	"os"

	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/shared/app"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "catalogs-read-microservices",
	Short:            "catalogs-read-microservices based on vertical slice architecture",
	Long:             `This is a command runner or cli for api architecture in golang.`,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		app.NewApp().Run()
	},
}

// https://github.com/swaggo/swag#how-to-use-it-with-gin

// @title Catalogs Read-Service Api
// @version 1.0
// @description Catalogs Read-Service Api for product reading and searching with comprehensive endpoints
// @termsOfService http://swagger.io/terms/

// @contact.name David Requeno
// @contact.url https://github.com/DavidReque
// @contact.email davidrequeno52@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:7001
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @description Enter "Bearer" followed by a space and the JWT token

// @tag.name Products
// @tag.description Operations about products

// @tag.name Catalogs
// @tag.description Catalog reading operations

// @tag.name Authentication
// @tag.description Authentication operations

// @x-extension-openapi {"example": "value on a json level"}

// @servers.url http://localhost:7001
// @servers.description Development server

// @servers.url https://api.production.com
// @servers.description Production server
func main() {
	pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("Catalogs", pterm.FgLightGreen.ToStyle()),
		putils.LettersFromStringWithStyle(" Read Service", pterm.FgLightMagenta.ToStyle())).
		Render()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
