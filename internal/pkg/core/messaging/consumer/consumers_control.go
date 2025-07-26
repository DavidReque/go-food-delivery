package consumer

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/types"
)

type BusControl interface {
	// Start starts all consumers
	Start(ctx context.Context) error
	// Stop stops all consumers
	Stop() error
	// IsConsumed checks if a message has been consumed
	IsConsumed(func(message types.IMessage))
}
