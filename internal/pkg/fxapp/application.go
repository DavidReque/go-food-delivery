package fxapp

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	"github.com/DavidReque/go-food-delivery/internal/pkg/fxapp/contracts"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"go.uber.org/fx"
)

type application struct {
	provides    []interface{}
	decorates   []interface{}
	invokes     []interface{}
	options     []fx.Option
	logger      logger.Logger
	fxapp       *fx.App
	environment environment.Environment
}

func NewApplication(
	provides []interface{},
	decorates []interface{},
	invokes []interface{},
	options []fx.Option,
	logger logger.Logger,
	env environment.Environment,
) contracts.Application {
	return &application{
		provides:    provides,
		decorates:   decorates,
		options:     options,
		logger:      logger,
		environment: env,
	}
}

// ResolveFunc registra una función para ser ejecutada cuando la aplicación se inicia
func (a *application) ResolveFunc(function interface{}) {
	a.invokes = append(a.invokes, function)
}

// ResolveFuncWithParamTag registra una función para ser ejecutada cuando la aplicación se inicia con un tag de parámetro
func (a *application) ResolveFuncWithParamTag(function interface{}, paramTagName string) {
	a.invokes = append(a.invokes, fx.Annotate(function, fx.ParamTags(paramTagName)))
}

// RegisterHook registra una función para ser ejecutada cuando la aplicación se inicia
func (a *application) RegisterHook(function interface{}) {
	a.invokes = append(a.invokes, function)
}

func (a *application) Run() {
	fxApp := CreateFxApp(a)

	a.fxapp = fxApp

	fxApp.Run()
}

func (a *application) Start(ctx context.Context) error {
	// build phase of container will do in this stage, containing provides and invokes but app not started yet and will be started in the future with `fxApp.Register`
	fxApp := CreateFxApp(a)
	a.fxapp = fxApp

	return fxApp.Start(ctx)
}

func (a *application) Stop(ctx context.Context) error {
	if a.fxapp == nil {
		a.logger.Fatal("Failed to stop because application not started.")
	}
	return a.fxapp.Stop(ctx)
}

// Wait espera a que la aplicación se detenga
func (a *application) Wait() <-chan fx.ShutdownSignal {
	if a.fxapp == nil {
		a.logger.Fatal("Failed to wait because application not started.")
	}
	return a.fxapp.Wait()
}

func (a *application) Logger() logger.Logger {
	return a.logger
}

func (a *application) Environment() environment.Environment {
	return a.environment
}
