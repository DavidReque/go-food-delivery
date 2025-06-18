package fxlog

// Package fxlog proporciona una implementación personalizada del logger para el framework Uber FX.
// Este paquete se encarga de manejar todos los eventos de logging del ciclo de vida de la aplicación.

import (
	"strings"

	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var FxLogger = fx.WithLogger(func(logger logger.Logger) fxevent.Logger {
	return NewCustomFxLogger(logger)
},
)

// FxCustomLogger es una estructura que implementa la interfaz fxevent.Logger
// Envuelve nuestro logger personalizado para adaptarlo al sistema de eventos de FX
type FxCustomLogger struct {
	logger.Logger
}

// NewCustomFxLogger crea una nueva instancia de nuestro logger personalizado
func NewCustomFxLogger(logger logger.Logger) fxevent.Logger {
	return &FxCustomLogger{Logger: logger}
}

// Printf implementa el método Printf de fxevent.Logger
// str: string a imprimir
// args: argumentos para el string
func (l FxCustomLogger) Printf(str string, args ...interface{}) {
	if len(args) > 0 {
		l.Debugf(str, args)
	}
	l.Debug(str)
}

// LogEvent es el método principal que maneja todos los eventos del ciclo de vida de FX
// Recibe diferentes tipos de eventos y los registra con el nivel de log apropiado
// Los eventos incluyen:
// - Inicio y parada de la aplicación
// - Ejecución de hooks OnStart y OnStop
// - Inyección de dependencias
// - Errores y rollbacks
func (l *FxCustomLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		// Registra cuando se está ejecutando un hook OnStart
		l.Debugw("OnStart hook executing", logger.Fields{"caller": e.CallerName, "function": e.FunctionName})
	case *fxevent.OnStartExecuted:
		// Registra el resultado de la ejecución de un hook OnStart
		if e.Err != nil {
			l.Errorw("OnStart hook failed",
				logger.Fields{"caller": e.CallerName, "callee": e.CallerName, "error": e.Err},
			)
		} else {
			l.Debugw("OnStart hook executed", logger.Fields{"caller": e.CallerName, "callee": e.FunctionName, "runtime": e.Runtime.String()})
		}
	case *fxevent.OnStopExecuting:
		// Registra cuando se está ejecutando un hook OnStop
		l.Debugw("OnStop hook executing", logger.Fields{"callee": e.FunctionName, "caller": e.CallerName})
	case *fxevent.OnStopExecuted:
		// Registra el resultado de la ejecución de un hook OnStop
		if e.Err != nil {
			l.Errorw("OnStop hook failed",
				logger.Fields{"caller": e.CallerName, "callee": e.CallerName, "error": e.Err},
			)
		} else {
			l.Debugw("OnStop hook executed", logger.Fields{"caller": e.CallerName, "callee": e.FunctionName, "runtime": e.Runtime.String()})
		}
	case *fxevent.Supplied:
		// Registra cuando se suministran dependencias
		if e.Err != nil {
			l.Errorw("error encountered while applying options",
				logger.Fields{"type": e.TypeName, "stacktrace": e.StackTrace, "module": e.ModuleName, "error": e.Err},
			)
		} else {
			l.Debugw("supplied", logger.Fields{"type": e.TypeName, "stacktrace": e.StackTrace, "module": e.ModuleName})
		}
	case *fxevent.Provided:
		// Registra cuando se proporcionan nuevos tipos/servicios
		for _, rtype := range e.OutputTypeNames {
			l.Debugw("provided", logger.Fields{"constructor": e.ConstructorName, "stacktrace": e.StackTrace, "module": e.ModuleName, "type": rtype, "private": e.Private})
		}
		if e.Err != nil {
			l.Errorw("error encountered while applying options",
				logger.Fields{"module": e.ModuleName, "stacktrace": e.StackTrace, "error": e.Err},
			)
		}
	case *fxevent.Replaced:
		// Registra cuando se reemplazan tipos existentes
		for _, rtype := range e.OutputTypeNames {
			l.Debugw("replaced", logger.Fields{"stacktrace": e.StackTrace, "module": e.ModuleName, "type": rtype})
		}
		if e.Err != nil {
			l.Errorw("error encountered while replacing",
				logger.Fields{"module": e.ModuleName, "stacktrace": e.StackTrace, "error": e.Err},
			)
		}
	case *fxevent.Decorated:
		// Registra cuando se decoran tipos existentes
		for _, rtype := range e.OutputTypeNames {
			l.Debugw("decorated", logger.Fields{"decorator": e.DecoratorName, "stacktrace": e.StackTrace, "module": e.ModuleName, "type": rtype})
		}
		if e.Err != nil {
			l.Errorw("error encountered while applying options",
				logger.Fields{"module": e.ModuleName, "stacktrace": e.StackTrace, "error": e.Err},
			)
		}
	case *fxevent.Run:
		// Registra eventos de ejecución
		if e.Err != nil {
			l.Errorw("error returned",
				logger.Fields{"module": e.ModuleName, "name": e.Name, "kind": e.Kind, "error": e.Err},
			)
		} else {
			l.Debugw("run", logger.Fields{"module": e.ModuleName, "name": e.Name, "kind": e.Kind})
		}
	case *fxevent.Invoking:
		// Registra cuando se está invocando una función
		l.Debugw("invoking", logger.Fields{"module": e.ModuleName, "function": e.FunctionName})
	case *fxevent.Invoked:
		// Registra el resultado de una invocación
		if e.Err != nil {
			l.Errorw("invoke failed",
				logger.Fields{"error": e.Err, "stack": e.Trace, "function": e.FunctionName, "module": e.ModuleName},
			)
		}
	case *fxevent.Stopping:
		// Registra cuando se recibe una señal de parada
		l.Debugw("received signal", logger.Fields{"signal": strings.ToUpper(e.Signal.String())})
	case *fxevent.Stopped:
		// Registra cuando la aplicación se ha detenido
		if e.Err != nil {
			l.Errorw("stop failed",
				logger.Fields{"error": e.Err},
			)
		}
	case *fxevent.RollingBack:
		// Registra cuando se inicia un rollback debido a un error
		l.Errorw("start failed, rolling back",
			logger.Fields{"error": e.StartErr},
		)
	case *fxevent.RolledBack:
		// Registra el resultado de un rollback
		if e.Err != nil {
			l.Errorw("rollback failed",
				logger.Fields{"error": e.Err},
			)
		}
	case *fxevent.Started:
		// Registra cuando la aplicación ha iniciado
		if e.Err != nil {
			l.Errorw("start failed",
				logger.Fields{"error": e.Err},
			)
		} else {
			l.Debug("started")
		}
	case *fxevent.LoggerInitialized:
		// Registra cuando el logger personalizado ha sido inicializado
		if e.Err != nil {
			l.Errorw("custom logger initialization failed",
				logger.Fields{"error": e.Err},
			)
		} else {
			l.Debugw("initialized custom fxevent.logger", logger.Fields{"function": e.ConstructorName})
		}
	}
}
