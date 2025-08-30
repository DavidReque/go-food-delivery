package dtos

import "github.com/DavidReque/go-food-delivery/internal/pkg/utils"

type GetOrdersRequestDto struct {
	*utils.ListQuery
}
