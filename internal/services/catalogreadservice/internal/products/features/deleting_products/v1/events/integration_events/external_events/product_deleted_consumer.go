package externaleEvents

import (
	"context"
	"errors"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/consumer"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/types"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/features/deleting_products/v1/commands"

	"github.com/go-playground/validator/v10"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
)

type productDeletedConsumer struct {
	logger    logger.Logger
	validator *validator.Validate
	tracer    tracing.AppTracer
}

func NewProductDeletedConsumer(
	logger logger.Logger,
	validator *validator.Validate,
	tracer tracing.AppTracer,
) consumer.ConsumerHandler {
	return &productDeletedConsumer{
		logger:    logger,
		validator: validator,
		tracer:    tracer,
	}
}

func (c *productDeletedConsumer) Handle(
	ctx context.Context,
	consumeContext types.MessageConsumeContext,
) error {
	// get message from consume context
	message, ok := consumeContext.Message().(*ProductDeletedV1)
	if !ok {
		return errors.New("error in casting message to ProductDeletedV1")
	}

	// convert product id to uuid
	productUUID, err := uuid.FromString(message.ProductId)
	if err != nil {
		badRequestErr := customErrors.NewBadRequestErrorWrap(
			err,
			"error in the converting uuid",
		)

		return badRequestErr
	}

	// create command
	command, err := commands.NewDeleteProduct(productUUID)
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"command validation failed",
		)

		return validationErr
	}

	// send command to delete product
	_, err = mediatr.Send[*commands.DeleteProduct, *mediatr.Unit](ctx, command)
	if err != nil {
		return err
	}

	// log success
	c.logger.Info("productDeletedConsumer executed successfully.")

	return err
}
