package gomigrate

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/migration"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"gomigratefx",
		mangoProviders,
	)

	mangoProviders = fx.Provide(
		migration.ProvideConfig,
		NewGoMigratorPostgres,
	)
)
