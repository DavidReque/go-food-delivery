package migration

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/config"
	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	typeMapper "github.com/DavidReque/go-food-delivery/internal/pkg/reflection/typemapper"
	"github.com/iancoleman/strcase"
)

type CommandType string

const (
	Up   CommandType = "up"
	Down CommandType = "down"
)

type MigrationOptions struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	User          string `mapstructure:"user"`
	DBName        string `mapstructure:"dbName"`
	SSLMode       bool   `mapstructure:"sslMode"`
	Password      string `mapstructure:"password"`
	// El nombre de la tabla que la herramienta de migración usará para llevar un
	// registro de qué migraciones ya se han aplicado
	VersionTable  string `mapstructure:"versionTable"`
	// La ruta a la carpeta que contiene las migraciones
	MigrationsDir string `mapstructure:"migrationsDir"`
	// Si se debe saltar la migración
	SkipMigration bool   `mapstructure:"skipMigration"`
}

func ProvideConfig(environment environment.Environment) (*MigrationOptions, error) {
	optionName := strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[MigrationOptions]())

	return config.BindConfigKey[MigrationOptions](optionName)
}
