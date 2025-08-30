package dtos

import (
	dtosV1 "github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/dtos/v1"

	customTypes "github.com/DavidReque/go-food-delivery/internal/pkg/core/customtypes"
)

// https://echo.labstack.com/guide/binding/
// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

// CreateOrderRequestDto validation will handle in command level
type CreateOrderRequestDto struct {
	ShopItems       []*dtosV1.ShopItemDto  `json:"shopItems"`
	AccountEmail    string                 `json:"accountEmail"`
	DeliveryAddress string                 `json:"deliveryAddress"`
	DeliveryTime    customTypes.CustomTime `json:"deliveryTime"`
}
