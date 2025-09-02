package dtos

import dtosV1 "github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/dtos/v1"

// GetOrderByIdResponseDto DTO for response to get order by id
// @Description DTO for response to get order by id
type GetOrderByIdResponseDto struct {
	// @Description Found order
	Order *dtosV1.OrderReadDto `json:"order"`
}
