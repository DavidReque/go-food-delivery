package integrationevents

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/types"

	uuid "github.com/satori/go.uuid"
)

type ProductDeletedV1 struct {
	*types.Message        // Mensaje base
	ProductId      string `json:"productId,omitempty"` // ID del producto eliminado
}

func NewProductDeletedV1(productId string) *ProductDeletedV1 {
	return &ProductDeletedV1{ProductId: productId, Message: types.NewMessage(uuid.NewV4().String())}
}
