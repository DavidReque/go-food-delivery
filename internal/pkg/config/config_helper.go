package config

import "github.com/spf13/viper"

// BindConfigKey vincula una clave de configuración a una estructura
// key: nombre de la clave en la configuración
// cfg: puntero a la estructura donde se almacenará la configuración
func BindConfigKey(key string, cfg interface{}) error {
	return viper.UnmarshalKey(key, cfg)
}
