package dtos

// CreateProductRequestDto specific request - only what the client sends
// @Description DTO for creating a new product with essential details
type CreateProductRequestDto struct {
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
}
