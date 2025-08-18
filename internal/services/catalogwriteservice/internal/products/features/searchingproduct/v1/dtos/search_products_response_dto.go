package dtos

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"
	dtoV1 "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1"
)

// SearchProductsResponseDto es el DTO para la respuesta de la búsqueda de productos
type SearchProductsResponseDto struct {
	Products *utils.ListResult[*dtoV1.ProductDto] // Lista de productos encontrados
}
