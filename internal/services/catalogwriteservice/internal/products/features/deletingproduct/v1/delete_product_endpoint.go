package v1

import (
	"net/http"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/web/route"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/deletingproduct/v1/dtos"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type deleteProductEndpoint struct {
	fxparams.ProductRouteParams // Inyección de parámetros de ruta
}

func NewDeleteProductEndpoint(
	params fxparams.ProductRouteParams,
) route.Endpoint {
	return &deleteProductEndpoint{ProductRouteParams: params}
}

// MapEndpoint maps the endpoint for the delete product command
func (ep *deleteProductEndpoint) MapEndpoint() {
	ep.ProductsGroup.DELETE("/:id", ep.handler())
}

// DeleteProduct
// @Tags Products
// @Summary Delete product
// @Description Delete existing product
// @Accept json
// @Produce json
// @Success 204
// @Param id path string true "Product ID"
// @Router /api/v1/products/{id} [delete]

// handler handles the delete product command
func (ep *deleteProductEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the context from the request
		ctx := c.Request().Context()

		// Bind the request to the delete product request dto
		request := &dtos.DeleteProductRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		// Create the delete product command
		command, err := NewDeleteProductWithValidation(request.ProductID)
		if err != nil {
			return err
		}

		// Send the delete product command to the mediator
		_, err = mediatr.Send[*DeleteProduct, *mediatr.Unit](
			ctx,
			command,
		)
		// Return the result
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending DeleteProduct",
			)
		}

		// Return the result
		return c.NoContent(http.StatusNoContent)
	}
}
