package mongodb

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/config"
	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	typeMapper "github.com/DavidReque/go-food-delivery/internal/pkg/reflection/typemapper"

	"github.com/iancoleman/strcase"
)

type MongoDbOptions struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	User          string `mapstructure:"user"`
	Password      string `mapstructure:"password"`
	Database      string `mapstructure:"database"`
	UseAuth       bool   `mapstructure:"useAuth"`
	EnableTracing bool   `mapstructure:"enableTracing" default:"true"`
}

func provideConfig(
	environment environment.Environment,
) (*MongoDbOptions, error) {
	optionName := strcase.ToLowerCamel(
		typeMapper.GetGenericTypeNameByT[MongoDbOptions](),
	)
	return config.BindConfigKey[MongoDbOptions](optionName)
}
