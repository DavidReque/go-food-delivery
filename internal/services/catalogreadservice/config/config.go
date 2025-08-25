package config

import (
	"strings"

	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	"github.com/spf13/viper"
)

type Config struct {
	AppOptions AppOptions `mapstructure:"appOptions" env:"AppOptions"`
}

func NewConfig(env environment.Environment) (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

type AppOptions struct {
	DeliveryType string `mapstructure:"deliveryType" env:"DeliveryType"`
	ServiceName  string `mapstructure:"serviceName"  env:"serviceName"`
}

func (cfg *AppOptions) GetMicroserviceNameUpper() string {
	return strings.ToUpper(cfg.ServiceName)
}

func (cfg *AppOptions) GetMicroserviceName() string {
	return cfg.ServiceName
}
