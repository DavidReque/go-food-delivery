package contracts

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/models"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
}