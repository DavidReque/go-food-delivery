package dtos

import "github.com/DavidReque/go-food-delivery/internal/pkg/utils"

type SearchProductsRequestDto struct {
	SearchText       string                                   `query:"search" json:"search"` // Texto de búsqueda
	*utils.ListQuery `                      json:"listQuery"` // Parámetros de paginación
}
