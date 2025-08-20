package data

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/models"
)

// ProductRepository is the interface for the product repository
type ProductRepository interface {
	GetAllProducts(
		ctx context.Context,
		listQuery *utils.ListQuery,
	) (*utils.ListResult[*models.Product], error)
	SearchProducts(
		ctx context.Context,
		searchText string,
		listQuery *utils.ListQuery,
	) (*utils.ListResult[*models.Product], error)
	GetProductById(ctx context.Context, uuid string) (*models.Product, error)
	GetProductByProductId(ctx context.Context, uuid string) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	DeleteProductByID(ctx context.Context, uuid string) error
}
