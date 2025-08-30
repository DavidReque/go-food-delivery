package dtos

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"
	dtosV1 "github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/dtos/v1"
)

type GetOrdersResponseDto struct {
	Orders *utils.ListResult[*dtosV1.OrderReadDto]
}
