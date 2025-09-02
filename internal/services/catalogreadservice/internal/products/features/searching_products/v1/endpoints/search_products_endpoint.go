package endpoints

import (
	"context"
	"net/http"

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
// @Description Search products by text with pagination and filtering capabilities
// @Accept json
// @Produce json
// @Param search query string false "Search text to find products" example("pizza")
// @Param page query int false "Page number" default(1) minimum(1)
// @Param size query int false "Page size" default(10) minimum(1) maximum(100)
// @Param orderBy query string false "Field to order by" example("name")
// @Param filters query string false "Applied filters" example("field=category&value=italian&comparison=eq")
// @Success 200 {object} dtos.SearchProductsResponseDto "Products found successfully"
// @Success 206 {object} dtos.SearchProductsResponseDto "Partial content - Paginated results"
// @Failure 400 {object} object "Bad request - Invalid search parameters"
// @Failure 401 {object} object "Unauthorized - Missing or invalid authentication"
// @Failure 403 {object} object "Forbidden - Insufficient permissions"
// @Failure 404 {object} object "Not found - No products match the search criteria"
// @Failure 429 {object} object "Too many requests - Rate limit exceeded"
// @Failure 500 {object} object "Internal server error - Something went wrong"
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
