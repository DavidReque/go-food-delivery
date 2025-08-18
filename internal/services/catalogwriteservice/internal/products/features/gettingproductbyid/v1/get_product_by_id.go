package v1

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/cqrs"
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/go-ozzo/ozzo-validation/is"
	uuid "github.com/satori/go.uuid"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

type GetProductById struct {
	cqrs.Query // Query es una interfaz que define los métodos para obtener un producto por su ID
	ProductID  uuid.UUID
}

func NewGetProductById(productId uuid.UUID) *GetProductById {
	query := &GetProductById{
		Query:     cqrs.NewQueryByT[GetProductById](),
		ProductID: productId,
	}

	return query
}

// NewGetProductByIdWithValidation crea una nueva consulta con validación
func NewGetProductByIdWithValidation(productId uuid.UUID) (*GetProductById, error) {
	// Crear la consulta
	query := NewGetProductById(productId)
	// Validar la consulta
	err := query.Validate()

	return query, err
}

// Validate valida la consulta
func (p *GetProductById) Validate() error {
	// Validar el ID del producto
	err := validation.ValidateStruct(
		p,
		validation.Field(&p.ProductID, validation.Required, is.UUIDv4),
	)
	// Si hay un error, devolver un error de validación
	if err != nil {
		return customErrors.NewValidationErrorWrap(err, "validation error")
	}

	return nil
}
