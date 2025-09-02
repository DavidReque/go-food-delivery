package main

import (
	"os"

	"github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/shared/app"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "orders-microservice",
	Short:            "orders-microservice based on vertical slice architecture",
	Long:             `This is a command runner or cli for api architecture in golang.`,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		app.NewApp().Run()
	},
}

// https://github.com/swaggo/swag#how-to-use-it-with-gin

// @title Orders Service Api
// @version 1.0
// @description Orders Service Api for order management
// @termsOfService http://swagger.io/terms/

// @contact.name David Requeno
// @contact.url https://github.com/DavidReque
// @contact.email davidrequeno52@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter "Bearer" followed by a space and the JWT token
func main() {
	pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("Orders", pterm.FgLightGreen.ToStyle()),
		putils.LettersFromStringWithStyle(" Service", pterm.FgLightMagenta.ToStyle())).
		Render()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
