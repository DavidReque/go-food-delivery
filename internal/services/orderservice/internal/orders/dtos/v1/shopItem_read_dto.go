package dtosV1

// ShopItemReadDto DTO for reading items from the shop in an order
// @Description DTO for reading items from the shop in an order
type ShopItemReadDto struct {
	// @Description Title of the product
	Title string `json:"title"`

	// @Description Description of the product
	Description string `json:"description"`

	// @Description Quantity of the product
	Quantity uint64 `json:"quantity"`

	// @Description Unit price of the product
	Price float64 `json:"price"`
}
