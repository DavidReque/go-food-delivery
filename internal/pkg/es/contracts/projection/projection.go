package projection

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/es/models"
)

type IProjection interface {
	ProcessEvent(ctx context.Context, streamEvent *models.StreamEvent) error
}