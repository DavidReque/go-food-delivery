package externaleEvents

import "github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/types"

// ProductDeletedV1 is the event that is sent when a product is deleted
type ProductDeletedV1 struct {
	*types.Message
	ProductId string `json:"productId,omitempty"`
}
