package queries

import (
	"context"

	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/contracts/data"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/dto"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/features/getting_products/v1/dtos"
)

type GetProductsHandler struct {
	log             logger.Logger
	mongoRepository data.ProductRepository
	tracer          tracing.AppTracer
}

func NewGetProductsHandler(
	log logger.Logger,
	mongoRepository data.ProductRepository,
	tracer tracing.AppTracer,
) *GetProductsHandler {
	return &GetProductsHandler{
		log:             log,
		mongoRepository: mongoRepository,
		tracer:          tracer,
	}
}

func (c *GetProductsHandler) Handle(
	ctx context.Context,
	query *GetProducts,
) (*dtos.GetProductsResponseDto, error) {
	// get products from mongo repository
	products, err := c.mongoRepository.GetAllProducts(ctx, query.ListQuery)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in getting products in the repository",
		)
	}

	// map products to list result dto
	listResultDto, err := utils.ListResultToListResultDto[*dto.ProductDto](
		products,
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping ListResultToListResultDto",
		)
	}

	c.log.Info("products fetched")

	return &dtos.GetProductsResponseDto{Products: listResultDto}, nil
}
