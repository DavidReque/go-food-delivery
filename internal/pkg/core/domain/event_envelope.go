package domain

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/metadata"
)

type EventEnvelope struct {
	EventData interface{}
	Metadata  metadata.Metadata
}



