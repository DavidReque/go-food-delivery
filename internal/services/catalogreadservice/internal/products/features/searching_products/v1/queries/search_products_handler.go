package queries

import (
	"context"

	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/contracts/data"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/dto"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/features/searching_products/v1/dtos"
)

type SearchProductsHandler struct {
	log             logger.Logger
	mongoRepository data.ProductRepository
	tracer          tracing.AppTracer
}

func NewSearchProductsHandler(
	log logger.Logger,
	repository data.ProductRepository,
	tracer tracing.AppTracer,
) *SearchProductsHandler {
	return &SearchProductsHandler{
		log:             log,
		mongoRepository: repository,
		tracer:          tracer,
	}
}

func (c *SearchProductsHandler) Handle(
	ctx context.Context,
	query *SearchProducts,
) (*dtos.SearchProductsResponseDto, error) {
	// search products in the repository
	products, err := c.mongoRepository.SearchProducts(
		ctx,
		query.SearchText,
		query.ListQuery,
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in searching products in the repository",
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

	return &dtos.SearchProductsResponseDto{Products: listResultDto}, nil
}
