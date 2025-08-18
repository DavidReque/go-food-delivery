package v1

import (
	"context"
	"fmt"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/cqrs"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/gormdbcontext"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/data/datamodels"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	integrationEvents "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/deletingproduct/v1/events/integrationevents"

	"github.com/mehdihadeli/go-mediatr"
)

type deleteProductHandler struct {
	fxparams.ProductHandlerParams
}

func NewDeleteProductHandler(
	params fxparams.ProductHandlerParams,
) cqrs.RequestHandlerWithRegisterer[*DeleteProduct, *mediatr.Unit] {
	return &deleteProductHandler{
		ProductHandlerParams: params,
	}
}

// RegisterHandler registers the handler for the delete product command
func (c *deleteProductHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*DeleteProduct, *mediatr.Unit](
		c,
	)
}

// IsTxRequest for enabling transactions on the mediatr pipeline
func (c *deleteProductHandler) isTxRequest() {
}

// Handle handles the delete product command
func (c *deleteProductHandler) Handle(
	ctx context.Context,
	command *DeleteProduct,
) (*mediatr.Unit, error) {
	// Delete the product from the database
	err := gormdbcontext.DeleteDataModelByID[*datamodels.ProductDataModel](ctx, c.CatalogsDBContext, command.ProductID)
	if err != nil {
		return nil, err
	}

	// Create the product deletion event
	productDeleted := integrationEvents.NewProductDeletedV1(
		command.ProductID.String(),
	)

	// Publish the product deletion event to RabbitMQ
	if err = c.RabbitmqProducer.PublishMessage(ctx, productDeleted); err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in publishing 'ProductDeleted' message",
		)
	}

	// Log the product deletion event
	c.Log.Infow(
		fmt.Sprintf(
			"ProductDeleted message with messageId '%s' published to the rabbitmq broker",
			productDeleted.MessageId,
		),
		logger.Fields{"MessageId": productDeleted.MessageId},
	)

	// Log the product deletion
	c.Log.Infow(
		fmt.Sprintf(
			"product with id '%s' deleted",
			command.ProductID,
		),
		logger.Fields{"Id": command.ProductID},
	)

	// Return the result
	return &mediatr.Unit{}, err
}
