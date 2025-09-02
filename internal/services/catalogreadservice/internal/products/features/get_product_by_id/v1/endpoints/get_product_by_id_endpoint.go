package endpoints

import (
	"net/http"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/web/route"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/contracts/params"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/features/get_product_by_id/v1/dtos"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/features/get_product_by_id/v1/queries"

	"emperror.dev/errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type getProductByIdEndpoint struct {
	params.ProductRouteParams
}

func NewGetProductByIdEndpoint(
	params params.ProductRouteParams,
) route.Endpoint {
	return &getProductByIdEndpoint{
		ProductRouteParams: params,
	}
}

func (ep *getProductByIdEndpoint) MapEndpoint() {
	ep.ProductsGroup.GET("/:id", ep.handler())
}

// GetProductByID
// @Tags Products
// @Summary Get product by ID
// @Description Get detailed information about a specific product by its unique identifier
// @Accept json
// @Produce json
// @Param id path string true "Product ID" format(uuid)
// @Success 200 {object} dtos.GetProductByIdResponseDto "Product details retrieved successfully"
// @Failure 400 {object} object "Bad request - Invalid product ID format"
// @Failure 401 {object} object "Unauthorized - Authentication required"
// @Failure 404 {object} object "Not found - Product not found"
// @Failure 500 {object} object "Internal server error - Something went wrong"
// @Router /api/v1/products/{id} [get]
func (ep *getProductByIdEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		// bind request
		request := &dtos.GetProductByIdRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		// parse product id
		productId, err := uuid.Parse(request.Id.String())
		if err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in parsing uuid",
			)

			return badRequestErr
		}

		// create query
		query, err := queries.NewGetProductById(productId)
		if err != nil {
			validationErr := customErrors.NewValidationErrorWrap(
				err,
				"query validation failed",
			)

			return validationErr
		}

		// send query to get product by id
		queryResult, err := mediatr.Send[*queries.GetProductById, *dtos.GetProductByIdResponseDto](
			ctx,
			query,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending GetProductById",
			)
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
