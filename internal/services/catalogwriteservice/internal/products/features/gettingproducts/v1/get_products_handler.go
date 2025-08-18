package v1

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/cqrs"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/helpers/gormextensions"
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/data/datamodels"
	dtosv1 "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/gettingproducts/v1/dtos"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/models"

	"github.com/mehdihadeli/go-mediatr"
)

type getProductsHandler struct {
	fxparams.ProductHandlerParams
}

// NewGetProductsHandler crea un nuevo manejador para obtener productos
func NewGetProductsHandler(
	params fxparams.ProductHandlerParams,
) cqrs.RequestHandlerWithRegisterer[*GetProducts, *dtos.GetProductsResponseDto] {
	return &getProductsHandler{
		ProductHandlerParams: params,
	}
}

// RegisterHandler registra el manejador para obtener productos
func (c *getProductsHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*GetProducts, *dtos.GetProductsResponseDto](
		c,
	)
}

// Handle maneja la solicitud para obtener productos
func (c *getProductsHandler) Handle(
	ctx context.Context,
	query *GetProducts, // Consulta para obtener productos
) (*dtos.GetProductsResponseDto, error) {
	// Obtener los productos de la base de datos
	products, err := gormextensions.Paginate[*datamodels.ProductDataModel, *models.Product](
		ctx,
		query.ListQuery,          // Parámetros de paginación
		c.CatalogsDBContext.DB(), // Contexto de la base de datos
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the fetching products",
		)
	}

	// Mapear los productos a un DTO
	listResultDto, err := utils.ListResultToListResultDto[*dtosv1.ProductDto](
		products,
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping",
		)
	}

	c.Log.Info("products fetched")

	// Devolver el resultado de la consulta
	return &dtos.GetProductsResponseDto{Products: listResultDto}, nil
}
