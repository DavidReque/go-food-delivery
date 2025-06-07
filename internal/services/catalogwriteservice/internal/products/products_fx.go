package products

import (
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/data/repositories"
	"go.uber.org/fx"
)

// Module is the module for the products service.
var Module = fx.Module(
	"productsfx",

	// Other providers
	fx.Provide(repositories.NewPostgresProductRepository),
)
