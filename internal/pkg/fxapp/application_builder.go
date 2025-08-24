package fxapp

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	"github.com/DavidReque/go-food-delivery/internal/pkg/fxapp/contracts"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger/config"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger/zap"

	"go.uber.org/fx"
)

// applicationBuilder es una estructura que contiene las opciones de la aplicación
type applicationBuilder struct {
	// constructores que fx debera registrar
	provides []interface{}
	// decoradores que modifican o envuelven esos providers.
	decorates []interface{}
	// módulos FX completos que también se añadirán.
	options     []fx.Option
	// invokes que se ejecutarán
	invokes     []interface{}
	logger      logger.Logger
	environment environment.Environment
}

func NewApplicationBuilder(environments ...environment.Environment) contracts.ApplicationBuilder {
	env := environment.ConfigAppEnv(environments...)

	var logger logger.Logger
	logoption, err := config.ProvideLogConfig(env)
	if err != nil {
		panic(err)
	}
	logger = zap.NewZapLogger(logoption, env)

	return &applicationBuilder{logger: logger, environment: env}
}

// AddModule agrega un módulo a la aplicación
func (a *applicationBuilder) AddModule(module ...fx.Option) {
	a.options = append(a.options, module...)
}

// ProvideModule agrega un módulo a la aplicación
func (a *applicationBuilder) ProvideModule(module fx.Option) {
	a.options = append(a.options, module)
}

// Provide agrega un constructor a la aplicación
func (a *applicationBuilder) Provide(constructors ...interface{}) {
	a.provides = append(a.provides, constructors...)
}

// Decorate agrega un decorador a la aplicación
func (a *applicationBuilder) Decorate(constructors ...interface{}) {
	a.decorates = append(a.decorates, constructors...)
}

// Build construye la aplicación
func (a *applicationBuilder) Build() contracts.Application {
	app := NewApplication(a.provides, a.decorates, a.invokes, a.options, a.logger, a.environment)
	return app
}

// GetProvides devuelve los constructores proporcionados
func (a *applicationBuilder) GetProvides() []interface{} {
	return a.provides
}

func (a *applicationBuilder) GetInvokes() []interface{} {
	return a.invokes
}

// GetDecorates devuelve los decoradores
func (a *applicationBuilder) GetDecorates() []interface{} {
	return a.decorates
}

// Options devuelve las opciones de la aplicación
func (a *applicationBuilder) Options() []fx.Option {
	return a.options
}

// Logger devuelve el logger de la aplicación
func (a *applicationBuilder) Logger() logger.Logger {
	return a.logger
}

// Environment devuelve el entorno de la aplicación
func (a *applicationBuilder) Environment() environment.Environment {
	return a.environment
}
