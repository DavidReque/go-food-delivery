package config

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	"github.com/spf13/viper"
)

// ConfigOptions contiene opciones para la configuración
type ConfigOptions struct {
	Environment environment.Environment
}

// ConfigOption es una función que modifica ConfigOptions
type ConfigOption func(*ConfigOptions)

// WithEnvironment establece el environment en las opciones
func WithEnvironment(env environment.Environment) ConfigOption {
	return func(o *ConfigOptions) {
		o.Environment = env
	}
}

// BindConfigKey vincula una clave de configuración a una estructura de tipo genérico T.
// Retorna un puntero a una nueva instancia de T con los valores de configuración.
// key: nombre de la clave en la configuración.
// opts: opciones adicionales de configuración
func BindConfigKey[T any](key string, opts ...ConfigOption) (*T, error) {
	options := &ConfigOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var cfg T
	if err := viper.UnmarshalKey(key, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
