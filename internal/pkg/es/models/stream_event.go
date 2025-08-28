package models

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/domain"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/metadata"

	"github.com/google/uuid"
)

type StreamEvent struct {
	EventID  uuid.UUID
	Version  int64
	Position int64
	Event    domain.IDomainEvent
	Metadata metadata.Metadata
}
