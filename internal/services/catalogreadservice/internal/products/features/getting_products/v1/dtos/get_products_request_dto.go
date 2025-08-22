package dtos

import "github.com/DavidReque/go-food-delivery/internal/pkg/utils"

type GetProductsRequestDto struct {
	*utils.ListQuery
}
