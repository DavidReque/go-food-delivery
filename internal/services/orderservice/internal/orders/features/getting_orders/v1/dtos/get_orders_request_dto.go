package dtos

import "github.com/DavidReque/go-food-delivery/internal/pkg/utils"

// GetOrdersRequestDto DTO to query orders with pagination and filters
// @Description DTO to query orders with pagination and filters
type GetOrdersRequestDto struct {
	// @Description Pagination and filters parameters
	*utils.ListQuery
}
