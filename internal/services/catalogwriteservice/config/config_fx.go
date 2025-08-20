package config

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("appconfigfx",
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	fx.Provide(
		NewAppOptions,
	),
	fx.Invoke(loadServiceConfig),
)

func loadServiceConfig(env environment.Environment) error {
	// Seleccionar archivo según entorno
	configPath := "config/config.development.json"
	if env.IsProduction() {
		configPath = "config/config.production.json"
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")
	return viper.ReadInConfig()
}
