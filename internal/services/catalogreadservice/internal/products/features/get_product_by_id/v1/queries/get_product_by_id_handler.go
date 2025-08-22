package queries

import (
	"context"
	"fmt"

	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/mapper"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/contracts/data"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/dto"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/features/get_product_by_id/v1/dtos"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/models"
)

type GetProductByIdHandler struct {
	log             logger.Logger
	mongoRepository data.ProductRepository
	redisRepository data.ProductCacheRepository
	tracer          tracing.AppTracer
}

func NewGetProductByIdHandler(
	log logger.Logger,
	mongoRepository data.ProductRepository,
	redisRepository data.ProductCacheRepository,
	tracer tracing.AppTracer,
) *GetProductByIdHandler {
	return &GetProductByIdHandler{
		log:             log,
		mongoRepository: mongoRepository,
		redisRepository: redisRepository,
		tracer:          tracer,
	}
}

func (q *GetProductByIdHandler) Handle(
	ctx context.Context,
	query *GetProductById,
) (*dtos.GetProductByIdResponseDto, error) {
	// get product from redis
	redisProduct, err := q.redisRepository.GetProductById(
		ctx,
		query.Id.String(),
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf(
				"error in getting product with id %d in the redis repository",
				query.Id,
			),
		)
	}

	// if product is not in redis, get it from mongo
	var product *models.Product

	// if product is in redis, use it
	if redisProduct != nil {
		product = redisProduct
	} else {
		// if product is not in redis, get it from mongo
		var mongoProduct *models.Product
		mongoProduct, err = q.mongoRepository.GetProductById(ctx, query.Id.String())
		if err != nil {
			return nil, customErrors.NewApplicationErrorWrap(err, fmt.Sprintf("error in getting product with id %d in the mongo repository", query.Id))
		}
		// if product is not in mongo, get it by product id
		if mongoProduct == nil {
			mongoProduct, err = q.mongoRepository.GetProductByProductId(ctx, query.Id.String())
		}
		if err != nil {
			return nil, err
		}

		// if product is in mongo, use it
		product = mongoProduct

		// put product in redis
		err = q.redisRepository.PutProduct(ctx, product.Id, product)
		if err != nil {
			return new(dtos.GetProductByIdResponseDto), err
		}
	}

	// map product to product dto
	productDto, err := mapper.Map[*dto.ProductDto](product)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping product",
		)
	}

	q.log.Infow(
		fmt.Sprintf(
			"product with id: {%s} fetched",
			query.Id,
		),
		logger.Fields{"ProductId": product.ProductId, "Id": product.Id},
	)

	return &dtos.GetProductByIdResponseDto{Product: productDto}, nil
}
