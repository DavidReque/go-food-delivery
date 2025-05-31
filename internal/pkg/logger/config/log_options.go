package config

import "github.com/DavidReque/go-food-delivery/internal/pkg/logger/models"

// LogOptions representa las opciones de configuración para el logger
type LogOptions struct {
	LogType       models.LogType // Tipo de logger a utilizar
	CallerEnabled bool            // Indica si se debe incluir el nombre del archivo y línea en los logs
}
