package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing/attribute"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing/utils"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/contracts/data"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/models"

	"emperror.dev/errors"
	"github.com/redis/go-redis/v9"
	attribute2 "go.opentelemetry.io/otel/attribute"
)

// redisProductPrefixKey is the prefix key for the product read service
// prefix key is used to store the product in the redis cache
const (
	redisProductPrefixKey = "product_read_service"
)

type redisProductRepository struct {
	log         logger.Logger
	redisClient redis.UniversalClient
	tracer      tracing.AppTracer
}

func NewRedisProductRepository(
	log logger.Logger,
	redisClient redis.UniversalClient,
	tracer tracing.AppTracer,
) data.ProductCacheRepository {
	return &redisProductRepository{
		log:         log,
		redisClient: redisClient,
		tracer:      tracer,
	}
}

// PutProduct put a product in the redis cache
func (r *redisProductRepository) PutProduct(
	ctx context.Context,
	key string,
	product *models.Product,
) error {
	// Start a new span
	ctx, span := r.tracer.Start(ctx, "redisRepository.PutProduct")
	// Set the prefix key attribute
	span.SetAttributes(
		attribute2.String("PrefixKey", r.getRedisProductPrefixKey()),
	)
	span.SetAttributes(attribute2.String("Key", key))
	defer span.End()

	// Marshal the product to json
	productBytes, err := json.Marshal(product)
	if err != nil {
		return utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"error marshalling product",
			),
		)
	}

	if err := r.redisClient.HSetNX(ctx, r.getRedisProductPrefixKey(), key, productBytes).Err(); err != nil {
		return utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"error in updating product with key %s",
					key,
				),
			),
		)
	}

	span.SetAttributes(attribute.Object("Product", product))

	r.log.Infow(
		fmt.Sprintf(
			"product with key '%s', prefix '%s'  updated successfully",
			key,
			r.getRedisProductPrefixKey(),
		),
		logger.Fields{
			"Product":   product,
			"Id":        product.ProductId,
			"Key":       key,
			"PrefixKey": r.getRedisProductPrefixKey(),
		},
	)

	return nil
}

// GetProductById get a product by id from the redis cache
func (r *redisProductRepository) GetProductById(
	ctx context.Context,
	key string,
) (*models.Product, error) {
	// Start a new span
	ctx, span := r.tracer.Start(ctx, "redisRepository.GetProductById")
	// Set the prefix key attribute
	span.SetAttributes(
		attribute2.String("PrefixKey", r.getRedisProductPrefixKey()),
	)
	span.SetAttributes(attribute2.String("Key", key))
	defer span.End()

	// Get the product from the redis cache
	productBytes, err := r.redisClient.HGet(ctx, r.getRedisProductPrefixKey(), key).
		Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"error in getting product with Key %s from database",
					key,
				),
			),
		)
	}

	// Unmarshal the product from json
	var product models.Product
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, utils.TraceErrStatusFromSpan(span, err)
	}

	span.SetAttributes(attribute.Object("Product", product))

	r.log.Infow(
		fmt.Sprintf(
			"product with with key '%s', prefix '%s' laoded",
			key,
			r.getRedisProductPrefixKey(),
		),
		logger.Fields{
			"Product":   product,
			"Id":        product.ProductId,
			"Key":       key,
			"PrefixKey": r.getRedisProductPrefixKey(),
		},
	)

	return &product, nil
}

// DeleteProduct delete a product from the redis cache
func (r *redisProductRepository) DeleteProduct(
	ctx context.Context,
	key string, // key is the id of the product
) error {
	// Start a new span
	ctx, span := r.tracer.Start(ctx, "redisRepository.DeleteProduct")
	// Set the prefix key attribute
	span.SetAttributes(
		attribute2.String("PrefixKey", r.getRedisProductPrefixKey()),
	)
	span.SetAttributes(attribute2.String("Key", key))
	defer span.End()

	if err := r.redisClient.HDel(ctx, r.getRedisProductPrefixKey(), key).Err(); err != nil {
		return utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"error in deleting product with key %s",
					key,
				),
			),
		)
	}

	r.log.Infow(
		fmt.Sprintf(
			"product with key %s, prefix: %s deleted successfully",
			key,
			r.getRedisProductPrefixKey(),
		),
		logger.Fields{"Key": key, "PrefixKey": r.getRedisProductPrefixKey()},
	)

	return nil
}

// DeleteAllProducts delete all products from the redis cache
func (r *redisProductRepository) DeleteAllProducts(ctx context.Context) error {
	ctx, span := r.tracer.Start(ctx, "redisRepository.DeleteAllProducts")
	span.SetAttributes(
		attribute2.String("PrefixKey", r.getRedisProductPrefixKey()),
	)
	defer span.End()

	// Delete all products from the redis cache
	if err := r.redisClient.Del(ctx, r.getRedisProductPrefixKey()).Err(); err != nil {
		return utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"error in deleting all products",
			),
		)
	}

	r.log.Infow(
		"all products deleted",
		logger.Fields{"PrefixKey": r.getRedisProductPrefixKey()},
	)

	return nil
}

// getRedisProductPrefixKey get the prefix key for the product read service
func (r *redisProductRepository) getRedisProductPrefixKey() string {
	return redisProductPrefixKey
}
