package postgresgorm

import (
	"fmt"
	"path/filepath"

	"github.com/DavidReque/go-food-delivery/internal/pkg/config"
	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	//typeMapper "github.com/DavidReque/go-food-delivery/internal/pkg/reflection/typemapper"
)

var optionName = "gorm" // strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[GormOptions]())

type GormOptions struct {
	UseInMemory   bool   `mapstructure:"useInMemory" json:"useInMemory"`
	UseSQLLite    bool   `mapstructure:"useSqlLite" json:"useSqlLite"`
	Host          string `mapstructure:"host" json:"host"`
	Port          int    `mapstructure:"port" json:"port"`
	User          string `mapstructure:"user" json:"user"`
	DBName        string `mapstructure:"dbName" json:"dbName"`
	SSLMode       bool   `mapstructure:"sslMode" json:"sslMode"`
	Password      string `mapstructure:"password" json:"password"`
	EnableTracing bool   `mapstructure:"enableTracing" default:"true" json:"enableTracing"`
}

func (h *GormOptions) Dns() string {
	if h.UseInMemory {
		return ""
	}

	if h.UseSQLLite {
		projectRootDir := environment.GetProjectRootWorkingDirectory()
		dbFilePath := filepath.Join(projectRootDir, fmt.Sprintf("%s.db", h.DBName))

		return dbFilePath
	}

	datasource := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		h.User,
		h.Password,
		h.Host,
		h.Port,
		h.DBName,
	)

	return datasource
}

func provideConfig(environment environment.Environment) (*GormOptions, error) {
	return config.BindConfigKey[GormOptions](optionName)
}
