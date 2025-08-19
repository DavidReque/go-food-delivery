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
	dto "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/updatingproduct/v1/events/integrationevents"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/models"

	"github.com/mehdihadeli/go-mediatr"
	satori_uuid "github.com/satori/go.uuid"
)

type updateProductHandler struct {
	fxparams.ProductHandlerParams
	cqrs.HandlerRegisterer
}

func NewUpdateProductHandler(
	params fxparams.ProductHandlerParams,
) cqrs.RequestHandlerWithRegisterer[*UpdateProduct, *mediatr.Unit] {
	return &updateProductHandler{
		ProductHandlerParams: params,
	}
}

// RegisterHandler registra el manejador para actualizar un producto
func (c *updateProductHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*UpdateProduct, *mediatr.Unit](
		c,
	)
}

// IsTxRequest for enabling transactions on the mediatr pipeline
func (c *updateProductHandler) isTxRequest() {
}

// Handle maneja la solicitud para actualizar un producto
func (c *updateProductHandler) Handle(
	ctx context.Context,
	command *UpdateProduct, // Comando para actualizar un producto
) (*mediatr.Unit, error) {
	// Convert google UUID to satori UUID
	satoriUUID := satori_uuid.FromStringOrNil(command.ProductID.String())

	// Buscar el producto en la base de datos
	product, err := gormdbcontext.FindModelByID[*datamodels.ProductDataModel, *models.Product](
		ctx,
		c.CatalogsDBContext, // Contexto de la base de datos
		satoriUUID,          // ID del producto
	)
	if err != nil {
		return nil, customErrors.NewNotFoundErrorWrap(
			err,
			fmt.Sprintf(
				"product with id `%s` not found",
				command.ProductID,
			),
		)
	}

	// Actualizar el producto
	product.Name = command.Name
	product.Price = command.Price
	product.Description = command.Description // Descripción del producto
	product.UpdatedAt = command.UpdatedAt

	// Actualizar el producto en la base de datos
	updatedProduct, err := gormdbcontext.UpdateModel[*datamodels.ProductDataModel, *models.Product](
		ctx,
		c.CatalogsDBContext, // Contexto de la base de datos
		product,
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in updating product in the repository",
		)
	}

	// Mapear el producto actualizado a un DTO
	productDto, err := mapper.Map[*dto.ProductDto](updatedProduct)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping ProductDto",
		)
	}

	// Crear el evento de actualización de un producto
	productUpdated := integrationevents.NewProductUpdatedV1(productDto)

	// Publicar el evento de actualización de un producto
	err = c.RabbitmqProducer.PublishMessage(ctx, productUpdated)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in publishing 'ProductUpdated' message",
		)
	}

	c.Log.Infow(
		fmt.Sprintf(
			"product with id '%s' updated",
			command.ProductID,
		),
		logger.Fields{"Id": command.ProductID},
	)

	c.Log.Infow(
		fmt.Sprintf(
			"ProductUpdated message with messageId `%s` published to the rabbitmq broker",
			productUpdated.MessageId,
		),
		logger.Fields{"MessageId": productUpdated.MessageId},
	)

	return &mediatr.Unit{}, err
}
