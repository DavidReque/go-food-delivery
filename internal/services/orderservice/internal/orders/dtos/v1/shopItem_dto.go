package dtosV1

// ShopItemDto DTO for representing an item from the shop in an order
// @Description DTO for representing an item from the shop in an order
type ShopItemDto struct {
	// @Description Title of the product
	// @Required
	Title string `json:"title"`

	// @Description Description of the product
	Description string `json:"description"`

	// @Description Quantity of the product
	// @Required
	// @Minimum 1
	Quantity uint64 `json:"quantity"`

	// @Description Unit price of the product
	// @Required
	// @Minimum 0
	Price float64 `json:"price"`
}
