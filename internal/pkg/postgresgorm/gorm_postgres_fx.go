package postgresgorm

import (
	"fmt"

	"github.com/DavidReque/go-food-delivery/internal/pkg/health/contracts"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"gormpostgresfx",
	fx.Provide(
		provideConfig,
		NewGorm,
		NewSQLDB,

		fx.Annotate(
			NewGormHealthChecker,
			fx.As(new(contracts.Health)),
			fx.ResultTags(fmt.Sprintf(`group:"%s"`, "healths")),
		),
	),
)