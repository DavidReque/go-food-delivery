package es

import (
	"fmt"

	"github.com/DavidReque/go-food-delivery/internal/pkg/es/contracts/projection"
	"go.uber.org/fx"
)

func AsProjection(handler interface{}) interface{} {
	return fx.Annotate(
		handler,
		fx.As(new(projection.IProjection)),
		fx.ResultTags(fmt.Sprintf(`group:"projections"`)),
	)
}
