package config

import (
	"fmt"

	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	"github.com/spf13/viper"
)

// LoadConfig lee config.<ENV>.json, vuelca en Config y aplica overrides de ENV.
func LoadConfig() (*Config, error) {
    // 1. Perfil: development / test / production
    env := environment.GetEnv("ENVIRONMENT", "development")

    // 2. Configuración de Viper para leer JSON
    viper.SetConfigType("json")
    viper.AddConfigPath("./config")
    viper.SetConfigName("config." + env)

    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("leyendo config.%s.json: %w", env, err)
    }

    // 3. Deserializar en nuestro struct
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("unmarshal config: %w", err)
    }

    // 4. Overrides desde ENV (si deseas forzar alguno)
    //    Ejemplo: cambiar puerto de Postgres o credenciales sin tocar JSON
    cfg.GormOptions.Host = environment.GetEnv("DB_HOST", cfg.GormOptions.Host)
    cfg.GormOptions.Port = environment.MustAtoi(
        environment.GetEnv("DB_PORT", fmt.Sprint(cfg.GormOptions.Port)),
        "DB_PORT",
    )
    cfg.GormOptions.User = environment.GetEnv("DB_USER", cfg.GormOptions.User)
    cfg.GormOptions.Password = environment.GetEnv("DB_PASSWORD", cfg.GormOptions.Password)
    cfg.GormOptions.DBName = environment.GetEnv("DB_NAME", cfg.GormOptions.DBName)

    // RabbitMQ overrides
    hostOpts := &cfg.RabbitMQOptions.RabbitMQHostOptions
    hostOpts.HostName = environment.GetEnv("RABBIT_HOST", hostOpts.HostName)
    hostOpts.Port = environment.MustAtoi(
        environment.GetEnv("RABBIT_PORT", fmt.Sprint(hostOpts.Port)),
        "RABBIT_PORT",
    )
    hostOpts.UserName = environment.GetEnv("RABBIT_USER", hostOpts.UserName)
    hostOpts.Password = environment.GetEnv("RABBIT_PASS", hostOpts.Password)

    return &cfg, nil
}
