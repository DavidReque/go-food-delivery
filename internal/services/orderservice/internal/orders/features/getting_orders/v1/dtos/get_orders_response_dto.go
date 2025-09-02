package dtos

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"
	dtosV1 "github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/dtos/v1"
)

// GetOrdersResponseDto DTO for response to get orders with pagination
// @Description DTO for response to get orders with pagination
type GetOrdersResponseDto struct {
	// @Description Paginated list of orders
	Orders *OrdersListResult
}

// OrdersListResult Specific paginated result for orders
// @Description Paginated result for orders
type OrdersListResult struct {
	// @Description Current page size
	Size int `json:"size,omitempty"`

	// @Description Current page number
	Page int `json:"page,omitempty"`

	// @Description Total available items
	TotalItems int64 `json:"totalItems,omitempty"`

	// @Description Total available pages
	TotalPage int `json:"totalPage,omitempty"`

	// @Description Current page orders
	Items []*dtosV1.OrderReadDto `json:"items,omitempty"`
}
