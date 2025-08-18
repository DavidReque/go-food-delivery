package v1

import (
	"net/http"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/web/route"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/gettingproducts/v1/dtos"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type getProductsEndpoint struct {
	fxparams.ProductRouteParams
}

func NewGetProductsEndpoint(
	params fxparams.ProductRouteParams,
) route.Endpoint {
	return &getProductsEndpoint{ProductRouteParams: params}
}

// MapEndpoint mapea el endpoint para obtener productos
func (ep *getProductsEndpoint) MapEndpoint() {
	ep.ProductsGroup.GET("", ep.handler())
}

// GetAllProducts
// @Tags Products
// @Summary Get all product
// @Description Get all products
// @Accfxparams
// @Produce json
// @Param getProductsRequestDto query dtos.GetProductsRequestDto false "GetProductsRequestDto"
// @Success 200 {object} dtos.GetProductsResponseDto
// @Router /api/v1/products [get]

// handler maneja la solicitud para obtener productos
func (ep *getProductsEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Obtener el contexto de la solicitud
		ctx := c.Request().Context()

		// Obtener los parámetros de paginación desde el contexto de Echo
		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			// Si hay un error, devolver un error de validación
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in getting data from query string",
			)

			return badRequestErr
		}

		// Crear el DTO de la solicitud
		request := &dtos.GetProductsRequestDto{ListQuery: listQuery}
		// Vincular la solicitud al DTO
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		// Crear la consulta con validación
		query, err := NewGetProducts(request.ListQuery)
		if err != nil {
			return err
		}

		// Enviar la consulta al mediador
		queryResult, err := mediatr.Send[*GetProducts, *dtos.GetProductsResponseDto](
			ctx,
			query,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending GetProducts",
			)
		}

		// Devolver el resultado de la consulta
		return c.JSON(http.StatusOK, queryResult)
	}
}
