package goose

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/migration"
	"go.uber.org/fx"
)

var (
	// Module provided to fxlog
	// https://uber-go.github.io/fx/modules.html
	Module = fx.Module( //nolint:gochecknoglobals
		"goosefx",
		mongoProviders,
	)

	mongoProviders = fx.Provide( //nolint:gochecknoglobals
		migration.ProvideConfig,
		NewGoosePostgres,
	)
)
