package v1

import (
	"net/http"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/web/route"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/gettingproductbyid/v1/dtos"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type getProductByIdEndpoint struct {
	fxparams.ProductRouteParams
}

func NewGetProductByIdEndpoint(
	params fxparams.ProductRouteParams,
) route.Endpoint {
	return &getProductByIdEndpoint{ProductRouteParams: params}
}

// MapEndpoint mapea el endpoint para obtener un producto por su ID
func (ep *getProductByIdEndpoint) MapEndpoint() {
	ep.ProductsGroup.GET("/:id", ep.handler())
}

// GetProductByID
// @Tags Products
// @Summary Get product by id
// @Description Get product by id
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dtos.GetProductByIdResponseDto
// @Router /api/v1/products/{id} [get]

// handler maneja la solicitud para obtener un producto por su ID
func (ep *getProductByIdEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Obtener el contexto de la solicitud
		ctx := c.Request().Context()

		// Crear el DTO de la solicitud
		request := &dtos.GetProductByIdRequestDto{}
		// Vincular la solicitud al DTO
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		// Crear la consulta con validación
		query, err := NewGetProductByIdWithValidation(request.ProductId)
		if err != nil {
			return err
		}

		// Enviar la consulta al mediador
		queryResult, err := mediatr.Send[*GetProductById, *dtos.GetProductByIdResponseDto](
			ctx,
			query,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending GetProductById",
			)
		}

		// Devolver el resultado de la consulta
		return c.JSON(http.StatusOK, queryResult)
	}
}
