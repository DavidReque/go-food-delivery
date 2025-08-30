package integrationEvents

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/types"
	dtosV1 "github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/dtos/v1"

	uuid "github.com/satori/go.uuid"
)

type OrderCreatedV1 struct {
	*types.Message
	*dtosV1.OrderReadDto
}

func NewOrderCreatedV1(orderReadDto *dtosV1.OrderReadDto) *OrderCreatedV1 {
	return &OrderCreatedV1{
		OrderReadDto: orderReadDto,
		Message:      types.NewMessage(uuid.NewV4().String()),
	}
}
