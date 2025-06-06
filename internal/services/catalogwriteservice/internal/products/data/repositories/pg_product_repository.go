package repositories

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/data"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/models"
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
