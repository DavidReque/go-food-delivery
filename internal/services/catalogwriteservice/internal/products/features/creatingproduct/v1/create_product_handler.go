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
	dtosv1 "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/dtos"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/events/integrationevents"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/models"
	"github.com/mehdihadeli/go-mediatr"
)

type createProductHandler struct {
	fxparams.ProductHandlerParams
}

func NewCreateProductHandler(
	params fxparams.ProductHandlerParams,
) cqrs.RequestHandlerWithRegisterer[*CreateProduct, *dtos.CreateProductResponseDto] {
	return &createProductHandler{
		ProductHandlerParams: params,
	}
}

func (c *createProductHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler(c)
}

func (c *createProductHandler) Handle(
	ctx context.Context,
	command *CreateProduct,
) (*dtos.CreateProductResponseDto, error) {
	product := &models.Product{
		Id:          command.ProductID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
		CreatedAt:   command.CreatedAt,
	}

	// Add the product to the database
	result, err := gormdbcontext.AddModel[*datamodels.ProductDataModel, *models.Product](
		ctx,
		c.CatalogsDBContext,
		product,
	)
	if err != nil {
		return nil, err
	}

	// Map the product to a ProductDto
	productDto, err := mapper.Map[*dtosv1.ProductDto](result)
	if err != nil {
		return nil, customErrors.NewApplicationError(
			err,
			"error in the mapping ProductDto",
		)
	}

	// Create the product creation event
	productCreated := integrationevents.NewProductCreatedV1(productDto)

	// Publish the product creation event
	err = c.RabbitmqProducer.PublishMessage(ctx, productCreated, nil)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in publishing ProductCreated integration_events event",
		)
	}

	// Log the product creation
	c.Log.Infow(
		fmt.Sprintf(
			"ProductCreated message with messageId `%s` published to the rabbitmq broker",
			productCreated.MessageId,
		),
		logger.Fields{"MessageId": productCreated.MessageId},
	)

	// Create the product creation result
	createProductResult := &dtos.CreateProductResponseDto{
		ProductID: product.Id,
	}

	// Log the product creation result
	c.Log.Infow(
		fmt.Sprintf(
			"product with id '%s' created",
			command.ProductID,
		),
		logger.Fields{
			"Id":        command.ProductID,
			"MessageId": productCreated.MessageId,
		},
	)

	// Return the product creation result
	return createProductResult, nil
}
