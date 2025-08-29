package models

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/domain"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/metadata"

	"github.com/google/uuid"
)


type StreamEvent struct {
	EventID  uuid.UUID // Unique identifier for the event
	Version  int64     // Version of the event
	Position int64     // Position of the event in the stream
	Event    domain.IDomainEvent // The actual event
	Metadata metadata.Metadata // Metadata associated with the event
}
