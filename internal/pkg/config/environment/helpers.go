package environment

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"emperror.dev/errors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/constants"
	"github.com/spf13/viper"
)

// MustAtoi convierte un string a entero y hace panic si falla.
// s: valor string a convertir.
// name: nombre de la variable (para mensajes de error descriptivos).
func MustAtoi(s, name string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("valor %q para %s no es un entero válido: %v", s, name, err))
	}
	return i
}

// ValidateURL comprueba que el string es una URL válida (tiene esquema y host).
// raw: string a validar como URL.
// name: nombre de la variable (para mensajes de error descriptivos).
func ValidateURL(raw, name string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		panic(fmt.Sprintf("valor %q para %s no es una URL válida", raw, name))
	}
	return raw
}

func FixProjectRootWorkingDirectoryPath() {
	currentWD, _ := os.Getwd()
	log.Printf("Current working directory is: `%s`", currentWD)

	rootDir := GetProjectRootWorkingDirectory()
	// change working directory
	_ = os.Chdir(rootDir)
	newWD, _ := os.Getwd()

	log.Printf("New fixed working directory is: `%s`", newWD)
}

func GetProjectRootWorkingDirectory() string {
	var rootWorkingDirectory string
	// https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
	// when we `Set` a viper with string value, we should get it from viper with `viper.GetString`, elsewhere we get empty string
	// viper will get it from `os env` or a .env file
	pn := viper.GetString(constants.PROJECT_NAME_ENV)
	if pn != "" {
		rootWorkingDirectory = getProjectRootDirectoryFromProjectName(pn)
	} else {
		wd, _ := os.Getwd()
		dir, err := searchRootDirectory(wd)
		if err != nil {
			log.Fatal(err)
		}
		rootWorkingDirectory = dir
	}

	absoluteRootWorkingDirectory, _ := filepath.Abs(rootWorkingDirectory)

	return absoluteRootWorkingDirectory
}

func getProjectRootDirectoryFromProjectName(pn string) string {
	// set root working directory of our app in the viper
	// https://stackoverflow.com/a/47785436/581476
	wd, _ := os.Getwd()

	for !strings.HasSuffix(wd, pn) {
		wd = filepath.Dir(wd)
	}

	return wd
}

func searchRootDirectory(
	dir string,
) (string, error) {
	// List files and directories in the current directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", errors.WrapIf(err, "Error reading directory")
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			if strings.EqualFold(
				fileName,
				"go.mod",
			) {
				return dir, nil
			}
		}
	}

	// If no config file found in this directory, recursively search its parent
	parentDir := filepath.Dir(dir)
	if parentDir == dir {
		// We've reached the root directory, and no go.mod file was found
		return "", errors.WrapIf(err, "No go.mod file found")
	}

	return searchRootDirectory(parentDir)
}

// Explicación:
// - MustAtoi asegura que los valores numéricos requeridos en la configuración sean válidos y detiene la ejecución si no lo son.
// - ValidateURL garantiza que las URLs requeridas sean válidas antes de que la aplicación continúe.
// - Ambas funciones ayudan a detectar errores de configuración temprano y con mensajes claros.
