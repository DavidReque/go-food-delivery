package dtos

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/dto"
)

type SearchProductsResponseDto struct {
	Products *utils.ListResult[*dto.ProductDto]
}
