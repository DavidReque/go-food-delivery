package contracts

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"

	"go.uber.org/fx"
)

type Application interface {
	// Container es el contenedor de dependencias
	Container
	// RegisterHook registra una función para ser ejecutada cuando la aplicación se inicia
	RegisterHook(function interface{})
	// Run ejecuta la aplicación
	Run()
	// Start inicia la aplicación
	Start(ctx context.Context) error
	// Stop detiene la aplicación
	Stop(ctx context.Context) error
	// Wait espera a que la aplicación se detenga
	Wait() <-chan fx.ShutdownSignal
	// Logger devuelve el logger de la aplicación
	Logger() logger.Logger
	// Environment devuelve el entorno de la aplicación
	Environment() environment.Environment
}
