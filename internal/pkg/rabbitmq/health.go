package rabbitmq

import (
	"context"

	"emperror.dev/errors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/health/contracts"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/types"
)

type gormHealthChecker struct {
	connection types.IConnection
}

func NewRabbitMQHealthChecker(connection types.IConnection) contracts.Health {
	return &gormHealthChecker{connection}
}

func (g *gormHealthChecker) CheckHealth(ctx context.Context) error {
	if g.connection.IsConnected() {
		return nil
	} else {
		return errors.New("rabbitmq is not available")
	}
}

func (g gormHealthChecker) GetHealthName() string {
	return "rabbitmq"
}
