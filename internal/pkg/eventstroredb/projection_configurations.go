package eventstroredb

import "github.com/DavidReque/go-food-delivery/internal/pkg/es/contracts/projection"

type ProjectionsConfigurations struct {
	Projections []projection.IProjection
}
