package infrastructure

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core"
	"go.uber.org/fx"
)

// https://pmihaylov.com/shared-components-go-microservices/


var Module = fx.Module(
	"infrastructurefx",
	// Modules
	core.Module,

	
)