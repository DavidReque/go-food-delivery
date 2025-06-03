package producer

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/types"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/metadata"
)

// Producer es una interfaz que define los métodos para publicar mensajes en un sistema de mensajería
type Producer interface {
	// PublishMessage publica un mensaje en el sistema de mensajería
	PublishMessage(ctx context.Context, message types.IMessage, metadata metadata.Metadata) error
	// PublishMessageWithTopicName publica un mensaje en el sistema de mensajería con un nombre de tema específico
	PublishMessageWithTopicName(
		ctx context.Context,
		message types.IMessage,
		metadata metadata.Metadata,
		totopicOrExchangeName string,
	) error
	// IsProduced verifica si un mensaje ha sido publicado
	IsProduced(func(message types.IMessage) bool)
}
