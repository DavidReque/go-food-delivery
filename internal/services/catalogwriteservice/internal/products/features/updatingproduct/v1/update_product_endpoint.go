package v1

import (
	"net/http"

	"emperror.dev/errors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/web/route"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/updatingproduct/v1/dtos"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type updateProductEndpoint struct {
	fxparams.ProductRouteParams
}

func NewUpdateProductEndpoint(
	params fxparams.ProductRouteParams,
) route.Endpoint {
	return &updateProductEndpoint{ProductRouteParams: params}
}

// MapEndpoint mapea el endpoint para actualizar un producto
func (ep *updateProductEndpoint) MapEndpoint() {
	ep.ProductsGroup.PUT("/:id", ep.handler())
}

// UpdateProduct
// @Tags Products
// @Summary Update product
// @Description Update existing product
// @Accept json
// @Produce json
// @Param UpdateProductRequestDto body dtos.UpdateProductRequestDto true "Product data"
// @Param id path string true "Product ID"
// @Success 204
// @Router /api/v1/products/{id} [put]

// handler maneja la solicitud para actualizar un producto
func (ep *updateProductEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		// Obtener el DTO de la solicitud
		request := &dtos.UpdateProductRequestDto{}
		// Vincular la solicitud al DTO
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		// Convert satori UUID to google UUID
		googleUUID, err := uuid.Parse(request.ProductID.String())
		if err != nil {
			return customErrors.NewBadRequestErrorWrap(err, "invalid product ID format")
		}

		// Crear la estructura para actualizar un producto con validación
		command, err := NewUpdateProductWithValidation(
			googleUUID,
			request.Name,
			request.Description,
			request.Price,
		)
		if err != nil {
			return err
		}

		// Enviar la estructura para actualizar un producto
		_, err = mediatr.Send[*UpdateProduct, *mediatr.Unit](
			ctx,
			command,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending UpdateProduct",
			)
		}

		// Devolver un resultado vacío
		return c.NoContent(http.StatusNoContent)
	}
}
