package zap

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger/config"
	"go.uber.org/fx"
)

var Module = fx.Module("zapfx",

	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	fx.Provide(
		config.ProvideLogConfig,
		NewZapLogger,
		fx.Annotate(
			NewZapLogger,
			fx.As(new(logger.Logger))),
	),
)

// ModuleFunc es una función que devuelve un módulo de FX
var ModuleFunc = func(l logger.Logger) fx.Option {
	return fx.Module(
		"zapfx",
		fx.Provide(config.ProvideLogConfig),
		fx.Supply(fx.Annotate(l, fx.As(new(logger.Logger)))),
		fx.Supply(fx.Annotate(l, fx.As(new(zapLogger)))),
	)
}
