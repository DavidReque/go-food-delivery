package externalEvents

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/consumer"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/types"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"
	v1 "github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/features/creating_product/v1"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/features/creating_product/v1/dtos"

	"github.com/go-playground/validator/v10"
	"github.com/mehdihadeli/go-mediatr"
)

type productCreatedConsumer struct {
	logger    logger.Logger
	validator *validator.Validate
	tracer    tracing.AppTracer
}

func NewProductCreatedConsumer(
	logger logger.Logger,
	validator *validator.Validate,
	tracer tracing.AppTracer,
) consumer.ConsumerHandler {
	return &productCreatedConsumer{
		logger:    logger,
		validator: validator,
		tracer:    tracer,
	}
}

func (c *productCreatedConsumer) Handle(
	ctx context.Context,
	consumeContext types.MessageConsumeContext,
) error {
	// get product from message
	product, ok := consumeContext.Message().(*ProductCreatedV1)
	if !ok {
		return errors.New("error in casting message to ProductCreatedV1")
	}

	// create product command
	command, err := v1.NewCreateProduct(
		product.ProductId,
		product.Name,
		product.Description,
		product.Price,
		product.CreatedAt,
	)
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"command validation failed",
		)

		return validationErr
	}

	// send create product command to mediator
	_, err = mediatr.Send[*v1.CreateProduct, *dtos.CreateProductResponseDto](
		ctx,
		command,
	)
	if err != nil {
		return errors.WithMessage(
			err,
			fmt.Sprintf(
				"error in sending CreateProduct with id: {%s}",
				command.ProductId,
			),
		)
	}
	c.logger.Info("Product consumer handled.")

	return err
}
