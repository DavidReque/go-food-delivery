package dtos

import "github.com/DavidReque/go-food-delivery/internal/pkg/utils"

type SearchProductsRequestDto struct {
	SearchText       string `query:"search" json:"search"`
	*utils.ListQuery `                      json:"listQuery"`
}
