package consumer

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/types"
)

// Consumer es una interfaz que define los métodos que debe implementar un consumidor.
type Consumer interface {
	Start(ctx context.Context) error
	Stop() error
	ConnectionHandler(handler ConsumerHandler)
	IsConsumed(func(message types.IMessage))
	GetName() string
}
