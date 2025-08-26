package redis

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/config"
	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	typeMapper "github.com/DavidReque/go-food-delivery/internal/pkg/reflection/typemapper"
	"github.com/iancoleman/strcase"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[RedisOptions]())

type RedisOptions struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	Username      string `mapstructure:"username"`
	Password      string `mapstructure:"password"`
	Database      int    `mapstructure:"database"`
	PoolSize      int    `mapstructure:"poolSize"`
	EnableTracing bool   `mapstructure:"enableTracing" default:"true"`
}

func provideConfig(environment environment.Environment) (*RedisOptions, error) {
	return config.BindConfigKey[RedisOptions](optionName)
}
