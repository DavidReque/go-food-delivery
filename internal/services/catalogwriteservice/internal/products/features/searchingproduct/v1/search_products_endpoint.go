package v1

import (
	"net/http"

	"emperror.dev/errors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/web/route"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/searchingproduct/v1/dtos"

	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type searchProductsEndpoint struct {
	fxparams.ProductRouteParams
}

// NewSearchProductsEndpoint crea un nuevo endpoint para buscar productos
func NewSearchProductsEndpoint(
	params fxparams.ProductRouteParams,
) route.Endpoint {
	return &searchProductsEndpoint{ProductRouteParams: params}
}

// MapEndpoint mapea el endpoint para buscar productos
func (ep *searchProductsEndpoint) MapEndpoint() {
	ep.ProductsGroup.GET("/search", ep.handler())
}

// SearchProducts
// @Tags Products
// @Summary Search products
// @Description Search products
// @Accept json
// @Produce json
// @Param searchProductsRequestDto query dtos.SearchProductsRequestDto false "SearchProductsRequestDto"
// @Success 200 {object} dtos.SearchProductsResponseDto
// @Router /api/v1/products/search [get]

// handler maneja la solicitud para buscar productos
func (ep *searchProductsEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Obtener el contexto de la solicitud
		ctx := c.Request().Context()

		// Obtener los parámetros de paginación desde el contexto de Echo
		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in getting data from query string",
			)

			return badRequestErr
		}

		// Crear el DTO de la solicitud
		request := &dtos.SearchProductsRequestDto{ListQuery: listQuery}
		// Vincular la solicitud al DTO
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		// Crear la consulta con validación
		query, err := NewSearchProductsWithValidation(
			request.SearchText, // Texto de búsqueda
			request.ListQuery,  // Parámetros de paginación
		)
		if err != nil {
			return err
		}

		// Enviar la consulta al mediador
		queryResult, err := mediatr.Send[*SearchProducts, *dtos.SearchProductsResponseDto](
			ctx,
			query, // Consulta para buscar productos
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending SearchProducts",
			)
		}

		// Devolver el resultado de la consulta
		return c.JSON(http.StatusOK, queryResult)
	}
}
