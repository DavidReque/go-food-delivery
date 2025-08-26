package endpoints

import (
	"net/http"
	"context"

	"emperror.dev/errors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/web/route"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/contracts/params"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/features/searching_products/v1/dtos"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/features/searching_products/v1/queries"

	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type searchProductsEndpoint struct {
	params.ProductRouteParams
}

func NewSearchProductsEndpoint(
	params params.ProductRouteParams,
) route.Endpoint {
	return &searchProductsEndpoint{
		ProductRouteParams: params,
	}
}

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
func (ep *searchProductsEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Use context.Background() instead of c.Request().Context() to avoid premature cancellation
		ctx := context.Background()

		// Start log
		c.Logger().Info("SearchProducts endpoint called")

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			c.Logger().Errorf("Error getting list query: %v", err)
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in getting data from query string",
			)

			return badRequestErr
		}

		c.Logger().Infof("ListQuery: %+v", listQuery)

		request := &dtos.SearchProductsRequestDto{ListQuery: listQuery}

		if err := c.Bind(request); err != nil {
			c.Logger().Errorf("Error binding request: %v", err)
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		c.Logger().Infof("Request: %+v", request)

		query := &queries.SearchProducts{
			SearchText: request.SearchText,
			ListQuery:  request.ListQuery,
		}

		if err := query.Validate(); err != nil {
			c.Logger().Errorf("Error validating query: %v", err)
			validationErr := customErrors.NewValidationErrorWrap(
				err,
				"error in the binding request",
			)

			return validationErr
		}

		c.Logger().Infof("Query: %+v", query)

		queryResult, err := mediatr.Send[*queries.SearchProducts, *dtos.SearchProductsResponseDto](
			ctx,
			query,
		)
		if err != nil {
			c.Logger().Errorf("Error in mediatr.Send: %v", err)
			return errors.WithMessage(
				err,
				"error in sending SearchProducts",
			)
		}

		c.Logger().Infof("QueryResult: %+v", queryResult)

		return c.JSON(http.StatusOK, queryResult)
	}
}
