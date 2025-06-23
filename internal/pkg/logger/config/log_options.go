package config

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/config"
	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger/models"
)

// Nombre de la opción de configuración para el logger
var optionName = "logOptions"

// LogOptions representa las opciones de configuración para el logger
type LogOptions struct {
	LogLevel       string         `mapstructure:"level"`
	LogType        models.LogType `mapstructure:"logType"`                      // Tipo de logger a utilizar
	CallerEnabled  bool           `mapstructure:"callerEnabled"`                // Indica si se debe incluir el nombre del archivo y línea en los logs
	EnabledTracing bool           `mapstructure:"enableTracing" default:"true"` // Indica si se debe incluir el tracing en los logs
}

// ProvideLogConfig proporciona las opciones de configuración para el logger
// env: entorno de la aplicación
// devuelve las opciones de configuración para el logger
func ProvideLogConfig(env environment.Environment) (*LogOptions, error) {
	return config.BindConfigKey[LogOptions](optionName)
}
