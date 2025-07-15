package infrastructure

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core"
	"github.com/DavidReque/go-food-delivery/internal/pkg/grpc"

	"github.com/DavidReque/go-food-delivery/internal/pkg/http/customecho"
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm"

	"go.uber.org/fx"
)

// https://pmihaylov.com/shared-components-go-microservices/


var Module = fx.Module(
	"infrastructurefx",
	// Modules
	core.Module,
	customecho.Module,
	grpc.Module,
	postgresgorm.Module,
)