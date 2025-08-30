package dtos

import dtosV1 "github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/dtos/v1"

type GetOrderByIdResponseDto struct {
	Order *dtosV1.OrderReadDto `json:"order"`
}
