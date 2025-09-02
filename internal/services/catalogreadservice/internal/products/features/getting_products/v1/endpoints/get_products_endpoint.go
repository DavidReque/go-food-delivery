package endpoints

import (
	"net/http"

	"emperror.dev/errors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/web/route"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/contracts/params"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/features/getting_products/v1/dtos"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/features/getting_products/v1/queries"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type getProductsEndpoint struct {
	params.ProductRouteParams
}

func NewGetProductsEndpoint(
	params params.ProductRouteParams,
) route.Endpoint {
	return &getProductsEndpoint{
		ProductRouteParams: params,
	}
}

func (ep *getProductsEndpoint) MapEndpoint() {
	ep.ProductsGroup.GET("", ep.handler())
}

// GetAllProducts
// @Tags Products
// @Summary Get all products
// @Description Get all products with pagination, filtering and sorting capabilities
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1)
// @Param size query int false "Page size" default(10) minimum(1) maximum(100)
// @Param orderBy query string false "Field to order by" example("createdAt")
// @Param filters query string false "Applied filters" example("field=name&value=pizza&comparison=contains")
// @Success 200 {object} dtos.GetProductsResponseDto "Products retrieved successfully"
// @Success 206 {object} dtos.GetProductsResponseDto "Partial content - Paginated results"
// @Failure 400 {object} object "Bad request - Invalid query parameters"
// @Failure 401 {object} object "Unauthorized - Missing or invalid authentication"
// @Failure 403 {object} object "Forbidden - Insufficient permissions"
// @Failure 404 {object} object "Not found - No products match the criteria"
// @Failure 429 {object} object "Too many requests - Rate limit exceeded"
// @Failure 500 {object} object "Internal server error - Something went wrong"
// @Router /api/v1/products [get]
func (ep *getProductsEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in getting data from query string",
			)

			return badRequestErr
		}

		// create query
		request := queries.NewGetProducts(listQuery)
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		// send query to get products
		query := &queries.GetProducts{ListQuery: request.ListQuery}

		queryResult, err := mediatr.Send[*queries.GetProducts, *dtos.GetProductsResponseDto](
			ctx,
			query,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending GetProducts",
			)
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
