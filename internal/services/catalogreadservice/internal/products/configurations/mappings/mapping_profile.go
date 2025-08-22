package mappings

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/mapper"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/dto"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/models"
)

func ConfigureProductsMappings() error {
	// create map for product to product dto
	err := mapper.CreateMap[*models.Product, *dto.ProductDto]()
	if err != nil {
		return err
	}

	// create map for product to product
	err = mapper.CreateMap[*models.Product, *models.Product]()
	if err != nil {
		return err
	}

	return nil
}
