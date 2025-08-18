package v1

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/cqrs"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/helpers/gormextensions"
	reflectionHelper "github.com/DavidReque/go-food-delivery/internal/pkg/reflection/reflectionhelper"
	"github.com/DavidReque/go-food-delivery/internal/pkg/reflection/typemapper"
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"
	datamodel "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/data/datamodels"
	dto "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/searchingproduct/v1/dtos"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/models"

	"github.com/iancoleman/strcase"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

type searchProductsHandler struct {
	fxparams.ProductHandlerParams
}

func NewSearchProductsHandler(
	params fxparams.ProductHandlerParams,
) cqrs.RequestHandlerWithRegisterer[*SearchProducts, *dtos.SearchProductsResponseDto] {
	return &searchProductsHandler{
		ProductHandlerParams: params,
	}
}

// RegisterHandler registra el manejador para buscar productos
func (c *searchProductsHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*SearchProducts, *dtos.SearchProductsResponseDto](
		c,
	)
}

// Handle maneja la solicitud para buscar productos
func (c *searchProductsHandler) Handle(
	ctx context.Context,
	query *SearchProducts, // Consulta para buscar productos
) (*dtos.SearchProductsResponseDto, error) {
	// Preparar la consulta a la base de datos
	dbQuery := c.prepareSearchDBQuery(query)

	// Obtener los productos de la base de datos
	products, err := gormextensions.Paginate[*datamodel.ProductDataModel, *models.Product](
		ctx,
		query.ListQuery, // Parámetros de paginación
		dbQuery,         // Consulta a la base de datos
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in searching products in the repository",
		)
	}

	// Mapear los productos a un DTO
	listResultDto, err := utils.ListResultToListResultDto[*dto.ProductDto](
		products,
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping ListResultToListResultDto",
		)
	}

	c.Log.Info("products fetched")

	return &dtos.SearchProductsResponseDto{Products: listResultDto}, nil
}

// prepareSearchDBQuery prepara la consulta a la base de datos
func (c *searchProductsHandler) prepareSearchDBQuery(
	query *SearchProducts, // Consulta para buscar productos
) *gorm.DB {
	// Obtener todos los campos de la estructura de datos
	fields := reflectionHelper.GetAllFields(
		typemapper.GetGenericTypeByT[*datamodel.ProductDataModel](),
	)

	// Inicializar la consulta a la base de datos
	dbQuery := c.CatalogsDBContext.DB()

	// Iterar sobre los campos de la estructura de datos
	for _, field := range fields {
		if field.Type.Kind() != reflect.String {
			continue
		}

		// Agregar una condición de búsqueda para el campo actual
		dbQuery = dbQuery.Or(
			fmt.Sprintf("%s LIKE ?", strcase.ToSnake(field.Name)),
			"%"+strings.ToLower(query.SearchText)+"%",
		)
	}

	return dbQuery
}
