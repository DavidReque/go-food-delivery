package v1

import (
	"time"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/cqrs"
	uuid "github.com/satori/go.uuid"
)

type CreateProduct struct {
	cqrs.Command
	ProductID   uuid.UUID
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
}