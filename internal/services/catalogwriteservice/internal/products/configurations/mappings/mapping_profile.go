package mappings

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/mapper"
	datamodel "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/data/datamodels"
	dtoV1 "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/models"
	productsService "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/shared/grpc/genproto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConfigureProductsMappings() error {
	// create product mappings
	// product to product dto
	err := mapper.CreateMap[*models.Product, *dtoV1.ProductDto]()
	if err != nil {
		return err
	}

	// create product dto to product mappings
	// product dto to product
	err = mapper.CreateMap[*dtoV1.ProductDto, *models.Product]()
	if err != nil {
		return err
	}

	// create product data model to product mappings
	// product data model to product
	err = mapper.CreateMap[*datamodel.ProductDataModel, *models.Product]()
	if err != nil {
		return err
	}

	// create product to product data model mappings
	// product to product data model
	err = mapper.CreateMap[*models.Product, *datamodel.ProductDataModel]()
	if err != nil {
		return err
	}

	// create product dto to product grpc mappings
	// product dto to product grpc
	err = mapper.CreateCustomMap[*dtoV1.ProductDto, *productsService.Product](
		func(product *dtoV1.ProductDto) *productsService.Product {
			if product == nil {
				return nil
			}
			return &productsService.Product{
				ProductId:   product.Id.String(),
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				CreatedAt:   timestamppb.New(product.CreatedAt),
				UpdatedAt:   timestamppb.New(product.UpdatedAt),
			}
		},
	)
	if err != nil {
		return err
	}

	// create product to product grpc mappings
	// product to product grpc
	err = mapper.CreateCustomMap(
		func(product *models.Product) *productsService.Product {
			return &productsService.Product{
				ProductId:   product.Id.String(),
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				CreatedAt:   timestamppb.New(product.CreatedAt),
				UpdatedAt:   timestamppb.New(product.UpdatedAt),
			}
		},
	)

	return nil
}
