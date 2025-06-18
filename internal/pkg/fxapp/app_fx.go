package fxapp

import (
	"time"

	"github.com/DavidReque/go-food-delivery/internal/pkg/config"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger/external/fxlog"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger/zap"
	"go.uber.org/fx"
)

func CreateFxApp(app *application) *fx.App {
	var opts []fx.Option

	// Registrar todos los proveedores (constructores) reunidos en el application
	opts = append(opts, fx.Provide(app.provides...))

	// Registrar los decoradores (funciones que modifican o envuelven proveedores)
	opts = append(opts, fx.Decorate(app.decorates...))

	// Registrar los decoradores (funciones que modifican o envuelven proveedores)
	opts = append(opts, fx.Invoke(app.invokes...))

	// Agregar esas opciones al slice general de la app
	app.options = append(app.options, opts...)

	// Crear el módulo de la aplicación
	AppModule := fx.Module(
		"fxapp",
		app.options...,
	)
	// Módulo de logging (Zap) y de configuración
	logModule := zap.ModuleFunc(app.logger)
	// Configurar el tiempo de espera para el inicio de la aplicación
	duration := 30 * time.Second

	// Construir el fx.App con:
	//    - timeout de arranque
	//    - módulo de configuración
	//    - módulo de logging externo para FX
	//    - error hook (para capturar fallos de fx.Invoke)
	//    - el módulo que engloba todos los providers/decorators/invokes
	fxApp := fx.New(
		fx.StartTimeout(duration),
		config.ModuleFunc(app.environment),
		logModule,
		fxlog.FxLogger,
		fx.ErrorHook(NewFxErrorHandler(app.logger)),
		AppModule,
	)

	return fxApp
}
