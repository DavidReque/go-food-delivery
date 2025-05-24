package integrationevents

import dtoV1 "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1"

type ProductCreatedV1 struct {
	//*types.Message
	*dtoV1.ProductDto
}

func NewProductCreatedV1(productDto *dtoV1.ProductDto) *ProductCreatedV1 {
	return &ProductCreatedV1{
		ProductDto: productDto,
		//Message:    types.NewMessage(uuid.NewV4().String()),
	}
}
