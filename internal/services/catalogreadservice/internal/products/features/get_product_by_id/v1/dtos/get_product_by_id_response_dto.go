package dtos

import "github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/dto"

type GetProductByIdResponseDto struct {
	Product *dto.ProductDto `json:"product"`
}
