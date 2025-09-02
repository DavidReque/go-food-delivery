package v1

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// ProductDto represents data for external transfer (APIs, JSON)
// @Description DTO for representing a product with comprehensive details
type ProductDto struct {
	// @Description Unique identifier of the product
	// @Example "550e8400-e29b-41d4-a716-446655440000"
	Id uuid.UUID `json:"id"`

	// @Description Product name or title
	// @Required
	// @MinLength 1
	// @MaxLength 100
	// @Example "Margherita Pizza"
	Name string `json:"name"`

	// @Description Detailed product description
	// @MaxLength 500
	// @Example "Traditional Italian pizza with tomato sauce, mozzarella and basil"
	Description string `json:"description"`

	// @Description Product price in the local currency
	// @Required
	// @Minimum 0.01
	// @Example 12.99
	Price float64 `json:"price"`

	// @Description Timestamp when the product was created
	// @Format date-time
	// @Example "2023-12-01T10:00:00Z07:00"
	CreatedAt time.Time `json:"createdAt"`

	// @Description Timestamp of the last product update
	// @Format date-time
	// @Example "2023-12-01T10:00:00Z07:00"
	UpdatedAt time.Time `json:"updatedAt"`
}
