package environment

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/DavidReque/go-food-delivery/internal/pkg/constants"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Environment string

var (
	Development = Environment(constants.Dev)
	Test        = Environment(constants.Test)
	Production  = Environment(constants.Production)
)

// ConfigAppEnv configura el entorno de la aplicación
// environments: lista de entornos a considerar, si no se proporciona, se usa el entorno de desarrollo
// devuelve el entorno configurado
func ConfigAppEnv(environments ...Environment) Environment {
	environment := Environment("")
	if len(environments) > 0 {
		environment = environments[0]
	} else {
		environment = Development
	}

	// Configuramos viper para que lea las variables de entorno
	viper.AutomaticEnv()

	err := loadEnvFilesRecursive()
	if err != nil {
		log.Printf(".env file cannot be found, err: %v", err)
	}

	//setRootWorkingDirectoryEnvironment()

	//FixProjectRootWorkingDirectoryPath()

	manualEnv := os.Getenv(constants.AppEnv)

	if manualEnv != "" {
		environment = Environment(manualEnv)
	}

	return environment
}

func IsDevelopment(env Environment) bool {
	return env == Development
}

func (env Environment) IsProduction() bool {
	return env == Production
}

func (env Environment) IsTest() bool {
	return env == Test
}

func (env Environment) GetEnvironmentName() string {
	return string(env)
}

func EnvString(key, fallback string) string {
	if v, ok := syscall.Getenv(key); ok {
		return v
	}

	return fallback
}

// loadEnvFilesRecursive carga el archivo .env.local si existe en el directorio actual o en sus padres
func loadEnvFilesRecursive() error {
	// Cargamos el archivo .env.local si existe
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		envFilePath := filepath.Join(dir, ".env")
		err := godotenv.Load(envFilePath)
		if err == nil {
			return nil
		}

		// mover al directorio padre
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			// reached the root directory, stop searching
			break
		}

		dir = parentDir
	}

	return errors.New(".env file not found in the project hierarchy")
}

// setRootWorkingDirectoryEnvironment establece la variable de entorno APP_ROOT_PATH en viper
func setRootWorkingDirectoryEnvironment() {
	absoluteRootWorkingDirectory := GetProjectRootWorkingDirectory()

	// when we `Set` a viper with string value, we should get it from viper with `viper.GetString`, elsewhere we get empty string
	viper.Set(constants.AppRootPath, absoluteRootWorkingDirectory)
}
