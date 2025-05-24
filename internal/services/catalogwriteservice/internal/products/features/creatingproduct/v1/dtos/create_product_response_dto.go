package dtos

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type CreateProductResponseDto struct {
	ProductID uuid.UUID `json:"productId"`
}

func (c *CreateProductResponseDto) String() string {
	return fmt.Sprintf("Product created with ID: %s", c.ProductID)
}
