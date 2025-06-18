package config

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	"go.uber.org/fx"
)

// Module es un módulo estático de FX para la configuración.
// Ventajas:
// - Simple de usar, no requiere parámetros adicionales
// - Se inicializa una sola vez al inicio de la aplicación
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"configfx",
	fx.Provide(func() environment.Environment {
		return environment.ConfigAppEnv()
	}),
)

// ModuleFunc es una función que crea un módulo FX personalizable para la configuración.
// A diferencia de Module, esta función:
// - Acepta un parámetro de environment.Environment
// - Permite inyectar configuraciones personalizadas
// - Es ideal para pruebas o diferentes entornos
var ModuleFunc = func(e environment.Environment) fx.Option {
	return fx.Module(
		"configfx",
		fx.Provide(func() environment.Environment {
			return environment.ConfigAppEnv(e)
		}),
	)
}
