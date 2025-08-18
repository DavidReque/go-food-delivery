package v1

import "github.com/DavidReque/go-food-delivery/internal/pkg/utils"

// Ref: https://golangbot.com/inheritance/

type GetProducts struct {
	*utils.ListQuery
}

// NewGetProducts crea una nueva consulta para obtener productos
func NewGetProducts(query *utils.ListQuery) (*GetProducts, error) {
	return &GetProducts{ListQuery: query}, nil
}
