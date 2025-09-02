package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
)

// Constants para el tamaño de página y el número de página por defecto
const (
	defaultSize = 10
	defaultPage = 1
)

// FilterModel define the filtering criteria for queries
// @Description Model to define filters in queries
type FilterModel struct {
		// @Description Field to filter
	// @Required
	Field string `query:"field" json:"field"`

	// @Description Filter value
	// @Required
	Value string `query:"value" json:"value"`

	// @Description Comparison operator (eq, ne, gt, lt, gte, lte, contains)
	// @Required
	Comparison string `query:"comparison" json:"comparison"`
}

// ListResult contains the paginated result of a query
// @Description Paginated result of a query
type ListResult[T any] struct {
		// @Description Current page size
	Size int `json:"size,omitempty" bson:"size"`

	// @Description Current page number
	Page int `json:"page,omitempty" bson:"page"`

	// @Description Total available items
	TotalItems int64 `json:"totalItems,omitempty" bson:"totalItems"`

	// @Description Total available pages
	TotalPage int `json:"totalPage,omitempty" bson:"totalPage"`

	// @Description Current page items
	Items []T `json:"items,omitempty" bson:"items"`
}

// NewListResult creates a new instance of ListResult with the provided data
func NewListResult[T any](items []T, size int, page int, totalItems int64) *ListResult[T] {
	listResult := &ListResult[T]{
		Items:      items,
		Size:       size,
		Page:       page,
		TotalItems: totalItems,
	}

	listResult.TotalPage = getTotalPages(totalItems, size)

	return listResult
}

// String convierte el ListResult a una cadena JSON
func (p *ListResult[T]) String() string {
	j, _ := json.Marshal(p)
	return string(j)
}

// getTotalPages calcula el número total de páginas
func getTotalPages(totalCount int64, size int) int {
	d := float64(totalCount) / float64(size)
	return int(math.Ceil(d))
}

// ListQuery contains the parameters for paginated queries
// @Description Parameters for paginated queries
type ListQuery struct {
		// @Description Current page size
	// @Minimum 1
	// @Maximum 100
	// @Default 10
	Size int `query:"size" json:"size,omitempty"`

	// @Description Current page number
	// @Minimum 1
	// @Default 1
	Page int `query:"page" json:"page,omitempty"`

	// @Description Field to order the results
	OrderBy string `query:"orderBy" json:"orderBy,omitempty"`

	// @Description Filters applied to the query
	Filters []*FilterModel `query:"filters" json:"filters,omitempty"`
}

// NewListQuery creates a new instance of ListQuery
func NewListQuery(size int, page int) *ListQuery {
	return &ListQuery{Size: size, Page: page}
}

// NewListQueryFromQueryParams crea una nueva instancia de ListQuery desde parámetros de consulta
func NewListQueryFromQueryParams(size string, page string) *ListQuery {
	p := &ListQuery{Size: defaultSize, Page: defaultPage}

	if sizeNum, err := strconv.Atoi(size); err == nil && sizeNum != 0 {
		p.Size = sizeNum
	}

	if pageNum, err := strconv.Atoi(page); err == nil && pageNum != 0 {
		p.Page = pageNum
	}

	return p
}

// GetListQueryFromCtx obtiene los parámetros de paginación desde el contexto de Echo
func GetListQueryFromCtx(c echo.Context) (*ListQuery, error) {
	q := &ListQuery{}
	var page, size, orderBy string

	err := echo.QueryParamsBinder(c).
		CustomFunc("filters", func(values []string) []error {
			for _, v := range values {
				if v == "" {
					continue
				}
				f := &FilterModel{}
				if err := c.Bind(f); err != nil {
					return []error{err}
				}
				q.Filters = append(q.Filters, f)
			}
			return nil
		}).
		String("size", &size).
		String("page", &page).
		String("orderBy", &orderBy).
		BindError()

	if err = q.SetPage(page); err != nil {
		return nil, err
	}
	if err = q.SetSize(size); err != nil {
		return nil, err
	}
	q.SetOrderBy(orderBy)

	return q, nil
}

// SetSize establece el tamaño de página
func (q *ListQuery) SetSize(sizeQuery string) error {
	if sizeQuery == "" {
		q.Size = defaultSize
		return nil
	}
	n, err := strconv.Atoi(sizeQuery)
	if err != nil {
		return err
	}
	q.Size = n
	return nil
}

// SetPage establece el número de página
func (q *ListQuery) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Page = defaultPage
		return nil
	}
	n, err := strconv.Atoi(pageQuery)
	if err != nil {
		return err
	}
	q.Page = n
	return nil
}

// SetOrderBy establece el campo de ordenamiento
func (q *ListQuery) SetOrderBy(orderByQuery string) {
	q.OrderBy = orderByQuery
}

// GetOffset obtiene el desplazamiento para la consulta
func (q *ListQuery) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

// GetLimit obtiene el límite de elementos
func (q *ListQuery) GetLimit() int {
	return q.Size
}

// GetOrderBy obtiene el campo de ordenamiento
func (q *ListQuery) GetOrderBy() string {
	return q.OrderBy
}

// GetPage obtiene el número de página actual
func (q *ListQuery) GetPage() int {
	return q.Page
}

// GetSize obtiene el tamaño de página
func (q *ListQuery) GetSize() int {
	return q.Size
}

// GetQueryString obtiene la cadena de consulta
func (q *ListQuery) GetQueryString() string {
	return fmt.Sprintf("page=%v&size=%v&orderBy=%s", q.GetPage(), q.GetSize(), q.GetOrderBy())
}

// ListResultToListResultDto convierte un ListResult de un tipo a otro
func ListResultToListResultDto[TDto any, TModel any](
	listResult *ListResult[TModel],
) (*ListResult[TDto], error) {
	if listResult == nil {
		return nil, errors.New("listResult is nil")
	}

	// Nota: Necesitarás implementar tu propia función de mapeo o usar una biblioteca de mapeo
	items, err := MapItems[[]TDto](listResult.Items)
	if err != nil {
		return nil, err
	}

	return &ListResult[TDto]{
		Items:      items,
		Size:       listResult.Size,
		Page:       listResult.Page,
		TotalItems: listResult.TotalItems,
		TotalPage:  listResult.TotalPage,
	}, nil
}

// MapItems mapea elementos entre diferentes tipos usando reflection
func MapItems[TDto any, TModel any](items TModel) (TDto, error) {
	var dto TDto

	// Convertir el modelo a JSON
	jsonData, err := json.Marshal(items)
	if err != nil {
		return dto, fmt.Errorf("error al serializar el modelo: %w", err)
	}

	// Convertir JSON al tipo DTO
	if err := json.Unmarshal(jsonData, &dto); err != nil {
		return dto, fmt.Errorf("error al deserializar al DTO: %w", err)
	}

	return dto, nil
}

func MapItemsManual[TDto any, TModel any](items TModel) (TDto, error) {
	var dto TDto
	modelValue := reflect.ValueOf(items)
	dtoValue := reflect.ValueOf(&dto).Elem()

	// Si es un slice, mapear cada elemento
	if modelValue.Kind() == reflect.Slice {
		// Crear un nuevo slice del tipo DTO
		sliceType := reflect.SliceOf(dtoValue.Type().Elem())
		newSlice := reflect.MakeSlice(sliceType, modelValue.Len(), modelValue.Cap())

		// Mapear cada elemento del slice
		for i := 0; i < modelValue.Len(); i++ {
			modelItem := modelValue.Index(i)
			dtoItem := reflect.New(dtoValue.Type().Elem()).Elem()

			// Mapear campos
			if err := mapFields(modelItem, dtoItem); err != nil {
				return dto, fmt.Errorf("error mapeando elemento %d: %w", i, err)
			}

			newSlice.Index(i).Set(dtoItem)
		}

		dtoValue.Set(newSlice)
		return dto, nil
	}

	// Si no es un slice, mapear directamente
	if err := mapFields(modelValue, dtoValue); err != nil {
		return dto, fmt.Errorf("error mapeando campos: %w", err)
	}

	return dto, nil
}

// mapFields mapea los campos entre dos estructuras
func mapFields(src, dst reflect.Value) error {
	if src.Kind() == reflect.Ptr {
		src = src.Elem()
	}
	if dst.Kind() == reflect.Ptr {
		dst = dst.Elem()
	}

	// Verificar que ambos sean estructuras
	if src.Kind() != reflect.Struct || dst.Kind() != reflect.Struct {
		return fmt.Errorf("tanto la fuente como el destino deben ser estructuras")
	}

	dstType := dst.Type()

	// Iterar sobre los campos del destino
	for i := 0; i < dst.NumField(); i++ {
		dstField := dst.Field(i)
		dstFieldType := dstType.Field(i)

		// Buscar campo correspondiente en la fuente
		srcField := src.FieldByName(dstFieldType.Name)
		if !srcField.IsValid() {
			continue // Skip si el campo no existe en la fuente
		}

		// Si los tipos son compatibles, copiar el valor
		if srcField.Type() == dstField.Type() && dstField.CanSet() {
			dstField.Set(srcField)
			continue
		}

		// Manejar casos especiales de conversión de tipos
		if err := handleSpecialCases(srcField, dstField); err != nil {
			return fmt.Errorf("error en conversión especial para campo %s: %w", dstFieldType.Name, err)
		}
	}

	return nil
}

// handleSpecialCases maneja conversiones especiales entre tipos
func handleSpecialCases(src, dst reflect.Value) error {
	if !dst.CanSet() {
		return nil
	}

	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch dst.Kind() {
		case reflect.String:
			dst.SetString(strconv.FormatInt(src.Int(), 10))
		case reflect.Float32, reflect.Float64:
			dst.SetFloat(float64(src.Int()))
		}
	case reflect.Float32, reflect.Float64:
		switch dst.Kind() {
		case reflect.String:
			dst.SetString(strconv.FormatFloat(src.Float(), 'f', -1, 64))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			dst.SetInt(int64(src.Float()))
		}
	case reflect.String:
		switch dst.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v, err := strconv.ParseInt(src.String(), 10, 64); err == nil {
				dst.SetInt(v)
			}
		case reflect.Float32, reflect.Float64:
			if v, err := strconv.ParseFloat(src.String(), 64); err == nil {
				dst.SetFloat(v)
			}
		}
	}

	return nil
}
