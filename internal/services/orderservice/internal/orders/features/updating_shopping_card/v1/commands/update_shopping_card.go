package commands

import (
	dtosV1 "github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/dtos/v1"

	"github.com/google/uuid"
)

type UpdateShoppingCart struct {
	OrderId   uuid.UUID             `validate:"required"`
	ShopItems []*dtosV1.ShopItemDto `validate:"required"`
}

func NewUpdateShoppingCart(orderId uuid.UUID, shopItems []*dtosV1.ShopItemDto) *UpdateShoppingCart {
	return &UpdateShoppingCart{OrderId: orderId, ShopItems: shopItems}
}
