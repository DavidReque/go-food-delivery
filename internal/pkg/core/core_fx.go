package core

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/serializer/json"
	"go.uber.org/fx"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"corefx",
	fx.Provide(
		json.NewDefaultJsonSerializer,
		json.NewDefaultEventJsonSerializer,
		json.NewDefaultMessageJsonSerializer,
		json.NewDefaultMetadataJsonSerializer,
	),
)
