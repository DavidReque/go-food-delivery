package v1

import (
	"context"
	"fmt"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/cqrs"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/mapper"
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/gormdbcontext"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/data/datamodels"
	dtoV1 "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/gettingproductbyid/v1/dtos"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/models"

	"github.com/mehdihadeli/go-mediatr"
)

type GetProductByIDHandler struct {
	fxparams.ProductHandlerParams
}

func NewGetProductByIDHandler(
	params fxparams.ProductHandlerParams,
) cqrs.RequestHandlerWithRegisterer[*GetProductById, *dtos.GetProductByIdResponseDto] {
	return &GetProductByIDHandler{
		ProductHandlerParams: params,
	}
}

// RegisterHandler registra el manejador para obtener un producto por su ID
func (c *GetProductByIDHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*GetProductById, *dtos.GetProductByIdResponseDto](
		c,
	)
}

// Handle maneja la solicitud para obtener un producto por su ID
func (c *GetProductByIDHandler) Handle(
	ctx context.Context,
	query *GetProductById, // Consulta para obtener un producto por su ID
) (*dtos.GetProductByIdResponseDto, error) {
	// Obtener el producto de la base de datos
	product, err := gormdbcontext.FindModelByID[*datamodels.ProductDataModel, *models.Product](
		ctx,                 // Contexto de la solicitud
		c.CatalogsDBContext, // Contexto de la base de datos
		query.ProductID,     // ID del producto
	)
	if err != nil {
		return nil, err
	}

	// Mapear el producto a un DTO
	productDto, err := mapper.Map[*dtoV1.ProductDto](product)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping product",
		)
	}

	// Log del producto obtenido
	c.Log.Infow(
		fmt.Sprintf(
			"product with id: {%s} fetched",
			query.ProductID,
		),
		logger.Fields{"Id": query.ProductID.String()},
	)

	// Devolver el producto obtenido
	return &dtos.GetProductByIdResponseDto{Product: productDto}, nil
}
