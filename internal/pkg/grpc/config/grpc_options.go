package config

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/config"
	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	typeMapper "github.com/DavidReque/go-food-delivery/internal/pkg/reflection/typemapper"
	"github.com/iancoleman/strcase"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[GrpcOptions]())

type GrpcOptions struct {
	Port        string `mapstructure:"port"        env:"TcpPort"`
	Host        string `mapstructure:"host"        env:"Host"`
	Development bool   `mapstructure:"development" env:"Development"`
	Name        string `mapstructure:"name"        env:"ShortTypeName"`
}

func ProvideConfig(environment environment.Environment) (*GrpcOptions, error) {
	return config.BindConfigKey[GrpcOptions](optionName)
}
