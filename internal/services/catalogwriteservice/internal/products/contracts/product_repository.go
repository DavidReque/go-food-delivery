package contracts

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/models"
	uuid "github.com/satori/go.uuid"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	DeleteProductByID(ctx context.Context, uuid uuid.UUID) error
}