package dtos

// Request específico - solo lo que el cliente envía
type CreateProductRequestDto struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
