package dtos

import uuid "github.com/satori/go.uuid"

// https://echo.labstack.com/guide/response/
// CreateOrderResponseDto DTO for response to create orders
// @Description DTO for response to create orders
type CreateOrderResponseDto struct {
	// @Description Unique ID of the created order
	// @Required
	OrderId uuid.UUID `json:"Id"`
}
