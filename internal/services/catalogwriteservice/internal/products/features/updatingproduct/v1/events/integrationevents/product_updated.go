package integrationevents

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/types"
	dto "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1"

	uuid "github.com/satori/go.uuid"
)

// ProductUpdatedV1 es el evento de actualización de un producto
type ProductUpdatedV1 struct {
	*types.Message
	*dto.ProductDto
}

// NewProductUpdatedV1 crea un nuevo evento de actualización de un producto
func NewProductUpdatedV1(productDto *dto.ProductDto) *ProductUpdatedV1 {
	// Crear un nuevo mensaje con un UUID
	return &ProductUpdatedV1{
		Message: types.NewMessage(uuid.NewV4().String()),
		// Vincular el DTO del producto
		ProductDto: productDto, // DTO del producto
	}
}
