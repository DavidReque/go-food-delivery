package v1

import (
	"time"

	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

// UpdateProduct es la estructura para actualizar un producto
type UpdateProduct struct {
	ProductID   uuid.UUID
	Name        string
	Description string
	Price       float64
	UpdatedAt   time.Time
}

// NewUpdateProduct crea una nueva estructura para actualizar un producto
func NewUpdateProduct(
	productID uuid.UUID,
	name string,
	description string,
	price float64,
) *UpdateProduct {
	command := &UpdateProduct{
		ProductID:   productID,
		Name:        name,
		Description: description,
		Price:       price,
		UpdatedAt:   time.Now(),
	}

	return command
}

// NewUpdateProductWithValidation crea una nueva estructura para actualizar un producto con validación
func NewUpdateProductWithValidation(
	productID uuid.UUID,
	name string,
	description string,
	price float64,
) (*UpdateProduct, error) {
	// Crear la estructura para actualizar un producto
	command := NewUpdateProduct(productID, name, description, price)
	// Validar la estructura
	err := command.Validate()

	return command, err
}

// IsTxRequest for enabling transactions on the mediatr pipeline
func (c *UpdateProduct) isTxRequest() {
}

// Validate valida la estructura para actualizar un producto
func (c *UpdateProduct) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ProductID, validation.Required),
		validation.Field(
			&c.Name,
			validation.Required,
			validation.Length(0, 255),
		),
		validation.Field(
			&c.Description,
			validation.Required,
			validation.Length(0, 5000),
		),
		validation.Field(&c.Price, validation.Required, validation.Min(0.0)),
		validation.Field(&c.UpdatedAt, validation.Required),
	)
	if err != nil {
		return customErrors.NewValidationErrorWrap(err, "validation error")
	}

	return nil
}
