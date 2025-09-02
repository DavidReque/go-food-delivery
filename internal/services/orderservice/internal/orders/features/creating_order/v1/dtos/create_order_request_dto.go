package dtos

import (
	"time"

	dtosV1 "github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/dtos/v1"
)

// https://echo.labstack.com/guide/binding/
// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

// CreateOrderRequestDto validation will handle in command level
// @Description DTO to create a new order
type CreateOrderRequestDto struct {
	// @Description List of items from the shop
	// @Required
	ShopItems []*dtosV1.ShopItemDto `json:"shopItems"`

	// @Description Email of the user's account
	// @Required
	// @Pattern ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
	AccountEmail string `json:"accountEmail" example:"user@example.com"`

	// @Description Delivery address
	// @Required
	DeliveryAddress string `json:"deliveryAddress" example:"Main Street 123, City"`

	// @Description Delivery time in RFC3339 format
	// @Required
	// @Format date-time
	DeliveryTime time.Time `json:"deliveryTime"`
}
