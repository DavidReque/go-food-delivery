package config

import "github.com/spf13/viper"

// BindConfigKey vincula una clave de configuración a una estructura de tipo genérico T.
// Retorna un puntero a una nueva instancia de T con los valores de configuración.
// key: nombre de la clave en la configuración.
func BindConfigKey[T any](key string) (*T, error) {
	var cfg T
	if err := viper.UnmarshalKey(key, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
