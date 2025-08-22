package externalEvents

import (
	"context"
	"fmt"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/consumer"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/types"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing/attribute"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing/utils"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/features/updating_products/v1/commands"

	"emperror.dev/errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	satoriuuid "github.com/satori/go.uuid"
)

type productUpdatedConsumer struct {
	logger    logger.Logger
	validator *validator.Validate
	tracer    tracing.AppTracer
}

func NewProductUpdatedConsumer(
	logger logger.Logger,
	validator *validator.Validate,
	tracer tracing.AppTracer,
) consumer.ConsumerHandler {
	return &productUpdatedConsumer{
		logger:    logger,
		validator: validator,
		tracer:    tracer,
	}
}

func (c *productUpdatedConsumer) Handle(
	ctx context.Context,
	consumeContext types.MessageConsumeContext,
) error {
	// get product from message
	message, ok := consumeContext.Message().(*ProductUpdatedV1)
	if !ok {
		return errors.New("error in casting message to ProductUpdatedV1")
	}

	// start span
	ctx, span := c.tracer.Start(ctx, "productUpdatedConsumer.Handle")
	span.SetAttributes(attribute.Object("Message", consumeContext.Message()))
	defer span.End()

	// convert product id to satori uuid
	productSatoriUUID, err := satoriuuid.FromString(message.ProductId)
	if err != nil {
		c.logger.WarnMsg("uuid.FromString", err)
		badRequestErr := customErrors.NewBadRequestErrorWrap(
			err,
			"[updateProductConsumer_Consume.uuid.FromString] error in the converting uuid",
		)
		c.logger.Errorf(
			fmt.Sprintf(
				"[updateProductConsumer_Consume.uuid.FromString] err: %v",
				utils.TraceErrStatusFromSpan(span, badRequestErr),
			),
		)
		return err
	}

	// convert satori uuid to google uuid
	productUUID, err := uuid.Parse(productSatoriUUID.String())
	if err != nil {
		c.logger.WarnMsg("uuid.Parse", err)
		badRequestErr := customErrors.NewBadRequestErrorWrap(
			err,
			"[updateProductConsumer_Consume.uuid.Parse] error in the converting uuid",
		)
		c.logger.Errorf(
			fmt.Sprintf(
				"[updateProductConsumer_Consume.uuid.Parse] err: %v",
				utils.TraceErrStatusFromSpan(span, badRequestErr),
			),
		)
		return err
	}

	// create update product command
	command, err := commands.NewUpdateProduct(
		productUUID,
		message.Name,
		message.Description,
		message.Price,
	)
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"[updateProductConsumer_Consume.NewValidationErrorWrap] command validation failed",
		)
		c.logger.Errorf(
			fmt.Sprintf(
				"[updateProductConsumer_Consume.StructCtx] err: {%v}",
				utils.TraceErrStatusFromSpan(span, validationErr),
			),
		)
		return err
	}

	// send update product command to mediator
	_, err = mediatr.Send[*commands.UpdateProduct, *mediatr.Unit](ctx, command)
	if err != nil {
		err = errors.WithMessage(
			err,
			"[updateProductConsumer_Consume.Send] error in sending UpdateProduct",
		)
		c.logger.Errorw(
			fmt.Sprintf(
				"[updateProductConsumer_Consume.Send] id: {%s}, err: {%v}",
				command.ProductId,
				utils.TraceErrStatusFromSpan(span, err),
			),
			logger.Fields{"Id": command.ProductId},
		)
		return err
	}

	return nil
}
