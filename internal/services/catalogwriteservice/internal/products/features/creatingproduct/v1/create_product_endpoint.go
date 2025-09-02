package v1

import (
	"net/http"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/web/route"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/dtos"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

// createProductEndpoint es una estructura que representa el endpoint para crear un producto.
type createProductEndpoint struct {
	fxparams.ProductRouteParams
}

func NewCreteProductEndpoint(
	params fxparams.ProductRouteParams,
) route.Endpoint {
	return &createProductEndpoint{ProductRouteParams: params}
}

func (ep *createProductEndpoint) MapEndpoint() {
	ep.ProductsGroup.POST("", ep.handler())
}

// CreateProduct
// @Tags Products
// @Summary Create new product
// @Description Create a new product with name, description and price
// @Accept json
// @Produce json
// @Param CreateProductRequestDto body dtos.CreateProductRequestDto true "Product data including name, description and price"
// @Success 201 {object} dtos.CreateProductResponseDto "Product created successfully"
// @Failure 400 {object} object "Bad request - Invalid product data"
// @Failure 401 {object} object "Unauthorized - Authentication required"
// @Failure 422 {object} object "Validation error - Invalid input data"
// @Failure 500 {object} object "Internal server error - Something went wrong"
// @Router /api/v1/products [post]

// handler es una función que maneja la solicitud HTTP para crear un producto.
func (ep *createProductEndpoint) handler() echo.HandlerFunc {
	// handler es una función que maneja la solicitud HTTP para crear un producto.
	return func(c echo.Context) error {
		// Obtiene el contexto de la solicitud HTTP.
		ctx := c.Request().Context()

		// Crea un nuevo objeto CreateProductRequestDto.
		request := dtos.CreateProductRequestDto{}

		// Vincula los datos del request con el objeto request.
		if err := c.Bind(&request); err != nil {
			// Si hay un error al vincular los datos, se retorna un error.
			badRequestErr := customErrors.NewBadRequestError(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		// Crea un nuevo objeto CreateProduct.
		command, err := NewCreateProductWithValidation(
			request.Name,
			request.Description,
			request.Price,
		)
		if err != nil {
			return err
		}

		// Envía el comando a la cola de mensajes.
		result, err := mediatr.Send[*CreateProduct, *dtos.CreateProductResponseDto](
			ctx,
			command,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending CreateProduct",
			)
		}

		// Retorna el resultado de la creación del producto.
		return c.JSON(http.StatusCreated, result)
	}
}
