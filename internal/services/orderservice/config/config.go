package config

import (
	"strings"

	"github.com/DavidReque/go-food-delivery/internal/pkg/config"
	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
)

type Config struct {
	AppOptions AppOptions `mapstructure:"appOptions"`
}

func NewConfig(environment environment.Environment) (*Config, error) {
	cfg, err := config.BindConfigKey[Config]("orderservice", config.WithEnvironment(environment))
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

type AppOptions struct {
	DeliveryType string `mapstructure:"deliveryType"`
	ServiceName  string `mapstructure:"serviceName"`
}

func (cfg *AppOptions) GetMicroserviceNameUpper() string {
	return strings.ToUpper(cfg.ServiceName)
}

func (cfg *AppOptions) GetMicroserviceName() string {
	return cfg.ServiceName
}