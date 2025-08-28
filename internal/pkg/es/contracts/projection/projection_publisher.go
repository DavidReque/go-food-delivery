package projection

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/es/models"
)

type IProjectionPublisher interface {
	Publish(ctx context.Context, streamEvent *models.StreamEvent) error
}