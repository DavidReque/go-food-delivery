package queries

import "github.com/DavidReque/go-food-delivery/internal/pkg/utils"

type GetOrders struct {
	*utils.ListQuery
}

func NewGetOrders(query *utils.ListQuery) *GetOrders {
	return &GetOrders{ListQuery: query}
}
