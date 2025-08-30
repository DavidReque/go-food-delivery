package domainevEnts

import (
	"fmt"

	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"

	"github.com/google/uuid"
)

type OrderSubmittedV1 struct {
	OrderId uuid.UUID `json:"orderId" bson:"orderId,omitempty"`
}

func NewSubmitOrderV1(orderId uuid.UUID) (*OrderSubmittedV1, error) {
	if orderId == uuid.Nil {
		return nil, customErrors.NewDomainError(fmt.Sprintf("orderId {%s} is invalid", orderId))
	}

	event := OrderSubmittedV1{OrderId: orderId}

	return &event, nil
}
