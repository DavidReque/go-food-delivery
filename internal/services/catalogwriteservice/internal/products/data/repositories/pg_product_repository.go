package repositories

import (
	"context"
	"fmt"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/data"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"
	attribute "github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing/attribute"
	utils2 "github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing/utils"
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/repository"
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"
	data2 "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/contracts"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/models"

	"emperror.dev/errors"
	uuid "github.com/satori/go.uuid"
	attribute2 "go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm"
)

// postgresProductRepository implementa el repositorio de productos usando PostgreSQL como almacenamiento
// Esta estructura sigue el patrón Repository y utiliza GORM como ORM para el acceso a datos
type postgresProductRepository struct {
	// log es el logger para registrar operaciones y errores del repositorio
	log logger.Logger

	// gormGenericRepository es una implementación genérica del repositorio usando GORM
	// Proporciona operaciones CRUD básicas para la entidad Product
	gormGenericRepository data.GenericRepository[*models.Product]

	// tracer es el rastreador de OpenTelemetry para monitoreo y observabilidad
	// Permite rastrear las operaciones del repositorio en un sistema distribuido
	tracer tracing.AppTracer
}

func NewPostgresProductRepository(
	log logger.Logger,
	db *gorm.DB,
	tracer tracing.AppTracer,
) data2.ProductRepository {
	gormRepository := repository.NewGenericGormRepository[*models.Product](db)
	return &postgresProductRepository{
		log:                   log,
		gormGenericRepository: gormRepository,
		tracer:                tracer,
	}
}

// GetAllProducts implementa la interfaz ProductRepository
// Obtiene todos los productos de la base de datos
// Utiliza el repositorio genérico de GORM para obtener los datos
// Aplica el tracing para monitorear la operación
// Retorna la lista de productos y cualquier error que ocurra
func (p *postgresProductRepository) GetAllProducts(
	ctx context.Context,
	listQuery *utils.ListQuery,
) (*utils.ListResult[*models.Product], error) {
	// Inicia el span de OpenTelemetry para rastrear la operación
	ctx, span := p.tracer.Start(ctx, "postgresProductRepository.GetAllProducts")
	// Finaliza el span al finalizar la función
	defer span.End()

	// Obtiene todos los productos de la base de datos
	result, err := p.gormGenericRepository.GetAll(ctx, listQuery)
	// Establece el estado del trace basado en el error
	err = utils2.TraceStatusFromContext(
		ctx,
		errors.WrapIf(
			err,
			"error in the paginate",
		),
	)

	if err != nil {
		return nil, err
	}

	// Registra el resultado de la operación
	p.log.Infow(
		"products loaded",
		logger.Fields{"ProductsResult": result},
	)

	// Establece los atributos del span con el resultado de la operación
	span.SetAttributes(attribute.Object("ProductsResult", result))

	// Retorna la lista de productos y cualquier error que ocurra
	return result, nil
}

// SearchProducts implementa la interfaz ProductRepository
// Busca productos en la base de datos
// Utiliza el repositorio genérico de GORM para obtener los datos
// Aplica el tracing para monitorear la operación
// Retorna la lista de productos y cualquier error que ocurra
func (p *postgresProductRepository) SearchProducts(
	ctx context.Context,
	searchText string,
	listQuery *utils.ListQuery,
) (*utils.ListResult[*models.Product], error) {
	// Inicia el span de OpenTelemetry para rastrear la operación
	ctx, span := p.tracer.Start(ctx, "postgresProductRepository.SearchProducts")
	// Establece el atributo del span con el texto de búsqueda
	span.SetAttributes(attribute2.String("SearchText", searchText))
	defer span.End()

	// Busca productos en la base de datos
	result, err := p.gormGenericRepository.Search(ctx, searchText, listQuery)
	// Establece el estado del trace basado en el error
	err = utils2.TraceStatusFromContext(
		ctx,
		errors.WrapIf(
			err,
			"error in the paginate",
		),
	)
	if err != nil {
		return nil, err
	}

	// Registra el resultado de la operación
	p.log.Infow(
		fmt.Sprintf(
			"products loaded for search term '%s'",
			searchText,
		),
		logger.Fields{"ProductsResult": result},
	)

	// Establece los atributos del span con el resultado de la operación
	span.SetAttributes(attribute.Object("ProductsResult", result))

	// Retorna la lista de productos y cualquier error que ocurra
	return result, nil
}

// GetProductById implementa la interfaz ProductRepository
// Obtiene un producto por su ID
// Utiliza el repositorio genérico de GORM para obtener los datos
// Aplica el tracing para monitorear la operación
// Retorna el producto y cualquier error que ocurra
func (p *postgresProductRepository) GetProductById(
	ctx context.Context,
	uuid uuid.UUID,
) (*models.Product, error) {
	// Inicia el span de OpenTelemetry para rastrear la operación
	ctx, span := p.tracer.Start(ctx, "postgresProductRepository.GetProductById")
	// Establece el atributo del span con el ID del producto
	span.SetAttributes(attribute2.String("Id", uuid.String()))
	defer span.End()

	// Obtiene el producto por su ID
	product, err := p.gormGenericRepository.GetById(ctx, uuid)
	// Establece el estado del trace basado en el error
	err = utils2.TraceStatusFromSpan(
		span,
		errors.WrapIf(
			err,
			fmt.Sprintf(
				"can't find the product with id %s into the database.",
				uuid,
			),
		),
	)

	if err != nil {
		return nil, err
	}

	// Establece los atributos del span con el resultado de la operación
	span.SetAttributes(attribute.Object("Product", product))

	// Registra el resultado de la operación
	p.log.Infow(
		fmt.Sprintf(
			"product with id %s laoded",
			uuid.String(),
		),
		logger.Fields{"Product": product, "Id": uuid},
	)

	// Retorna el producto y cualquier error que ocurra
	return product, nil
}

// CreateProduct implementa la interfaz ProductRepository
// Crea un nuevo producto en la base de datos
// Utiliza el repositorio genérico de GORM para crear el producto
// Aplica el tracing para monitorear la operación
// Retorna el producto y cualquier error que ocurra
func (p *postgresProductRepository) CreateProduct(
	ctx context.Context,
	product *models.Product,
) (*models.Product, error) {
	// Inicia el span de OpenTelemetry para rastrear la operación
	ctx, span := p.tracer.Start(ctx, "postgresProductRepository.CreateProduct")
	// Finaliza el span al finalizar la función
	defer span.End()

	// Crea el nuevo producto en la base de datos
	err := p.gormGenericRepository.Add(ctx, product)
	err = utils2.TraceStatusFromSpan(
		span,
		errors.WrapIf(
			err,
			"error in the inserting product into the database.",
		),
	)
	if err != nil {
		return nil, err
	}

	// Establece los atributos del span con el resultado de la operación
	span.SetAttributes(attribute.Object("Product", product))
	p.log.Infow(
		fmt.Sprintf(
			"product with id '%s' created",
			product.Id,
		),
		logger.Fields{"Product": product, "Id": product.Id},
	)

	// Retorna el producto y cualquier error que ocurra
	return product, nil
}

// UpdateProduct implementa la interfaz ProductRepository
// Actualiza un producto en la base de datos
// Utiliza el repositorio genérico de GORM para actualizar el producto
// Aplica el tracing para monitorear la operación
// Retorna el producto y cualquier error que ocurra
func (p *postgresProductRepository) UpdateProduct(
	ctx context.Context,
	updateProduct *models.Product,
) (*models.Product, error) {
	ctx, span := p.tracer.Start(ctx, "postgresProductRepository.UpdateProduct")
	defer span.End()

	err := p.gormGenericRepository.Update(ctx, updateProduct)
	err = utils2.TraceStatusFromSpan(
		span,
		errors.WrapIf(
			err,
			fmt.Sprintf(
				"error in updating product with id %s into the database.",
				updateProduct.Id,
			),
		),
	)

	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Object("Product", updateProduct))
	p.log.Infow(
		fmt.Sprintf(
			"product with id '%s' updated",
			updateProduct.Id,
		),
		logger.Fields{
			"Product": updateProduct,
			"Id":      updateProduct.Id,
		},
	)

	return updateProduct, nil
}

// DeleteProductByID implementa la interfaz ProductRepository
// Elimina un producto por su ID
// Utiliza el repositorio genérico de GORM para eliminar el producto
// Aplica el tracing para monitorear la operación
// Retorna cualquier error que ocurra
func (p *postgresProductRepository) DeleteProductByID(
	ctx context.Context,
	uuid uuid.UUID,
) error {
	// Inicia el span de OpenTelemetry para rastrear la operación
	ctx, span := p.tracer.Start(ctx, "postgresProductRepository.UpdateProduct")
	// Establece el atributo del span con el ID del producto
	span.SetAttributes(attribute2.String("Id", uuid.String()))
	defer span.End()

	// Elimina el producto por su ID
	err := p.gormGenericRepository.Delete(ctx, uuid)
	// Establece el estado del trace basado en el error
	err = utils2.TraceStatusFromSpan(span, errors.WrapIf(err, fmt.Sprintf(
		"error in the deleting product with id %s into the database.",
		uuid,
	)))

	if err != nil {
		return err
	}

	// Registra el resultado de la operación
	p.log.Infow(
		fmt.Sprintf(
			"product with id %s deleted",
			uuid,
		),
		logger.Fields{"Product": uuid},
	)

	return nil
}
