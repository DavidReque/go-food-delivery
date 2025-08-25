package mongodb

import (
	"fmt"

	"github.com/DavidReque/go-food-delivery/internal/pkg/config"
	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
)

type MongoDbOptions struct {
	// Configuración para MongoDB local
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	UseAuth  bool   `mapstructure:"useAuth"`

	// MongoDB Atlas configuration
	AtlasURI string `mapstructure:"atlasUri"`
	UseAtlas bool   `mapstructure:"useAtlas"`

	// common configuration
	EnableTracing bool `mapstructure:"enableTracing" default:"true"`
}

func provideConfig(
	environment environment.Environment,
) (*MongoDbOptions, error) {
	// Use hardcoded name as in other services that work
	optionName := "mongoDbOptions"

	cfg, err := config.BindConfigKey[MongoDbOptions](optionName)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
