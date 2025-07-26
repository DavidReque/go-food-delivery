package consumer

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/types"
)

// ConsumerHandler es una interfaz que define los métodos que debe implementar un manejador de consumo de mensajes.
type ConsumerHandler interface {
	Handle(ctx context.Context, consumeContext types.MessageConsumeContext) error
}
