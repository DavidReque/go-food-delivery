package v1

import (
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"

	validation "github.com/go-ozzo/ozzo-validation"
)

// SearchProducts es la consulta para buscar productos
type SearchProducts struct {
	SearchText       string // Texto de búsqueda
	*utils.ListQuery        // Parámetros de paginación
}

// NewSearchProducts crea una nueva consulta para buscar productos
func NewSearchProducts(searchText string, query *utils.ListQuery) *SearchProducts {
	searchProductQuery := &SearchProducts{
		SearchText: searchText,
		ListQuery:  query,
	}

	return searchProductQuery
}

// NewSearchProductsWithValidation crea una nueva consulta para buscar productos con validación
func NewSearchProductsWithValidation(searchText string, query *utils.ListQuery) (*SearchProducts, error) {
	// Crear la consulta con validación
	searchProductQuery := NewSearchProducts(searchText, query)

	// Validar la consulta
	err := searchProductQuery.Validate()

	return searchProductQuery, err
}

// Validate valida la consulta para buscar productos
func (p *SearchProducts) Validate() error {
	err := validation.ValidateStruct(p, validation.Field(&p.SearchText, validation.Required))
	if err != nil {
		return customErrors.NewValidationErrorWrap(err, "validation error")
	}

	return nil
}
