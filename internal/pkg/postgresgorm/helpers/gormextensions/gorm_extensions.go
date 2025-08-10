package gormextensions

import (
	"context"
	"fmt"
	"strings"

	"github.com/DavidReque/go-food-delivery/internal/pkg/mapper"
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/constants"
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/contracts"
	"github.com/DavidReque/go-food-delivery/internal/pkg/reflection/typemapper"
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"

	"emperror.dev/errors"
	"gorm.io/gorm"
)

// GetTxFromContext obtiene la transacción de contexto
// Si la transacción está en el contexto, la devuelve
// Si no, busca en el contexto el valor de TxKey
// Si no encuentra la transacción, devuelve un error
func GetTxFromContext(ctx context.Context) (*gorm.DB, error) {
	gCtx, gCtxOk := ctx.(*contracts.GormContext)
	if gCtxOk {
		return gCtx.Tx, nil
	}

	tx, ok := ctx.Value(constants.TxKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("Transaction not found in context")
	}

	return tx, nil
}

// GetTxFromContextIfExists obtiene la transacción de contexto
// Si la transacción está en el contexto, la devuelve
// Si no, busca en el contexto el valor de TxKey
// Si no encuentra la transacción, devuelve nil
func GetTxFromContextIfExists(ctx context.Context) *gorm.DB {
	gCtx, gCtxOk := ctx.(*contracts.GormContext)
	if gCtxOk {
		return gCtx.Tx
	}

	tx, ok := ctx.Value(constants.TxKey).(*gorm.DB)
	if !ok {
		return nil
	}

	return tx
}

// SetTxToContext establece la transacción en el contexto
// ctx: contexto de Go para la operación
// tx: transacción de GORM
// Retorna: nuevo contexto con la transacción establecida
func SetTxToContext(ctx context.Context, tx *gorm.DB) *contracts.GormContext {
	newCtx := context.WithValue(ctx, constants.TxKey, tx)
	gormContext := &contracts.GormContext{Tx: tx, Context: newCtx}
	ctx = gormContext

	return gormContext
}

// Paginate implementa la paginación para consultas GORM
// TDataModel: tipo del modelo de datos
// TEntity: tipo de la entidad de retorno
func Paginate[TDataModel any, TEntity any](
	ctx context.Context,
	listQuery *utils.ListQuery,
	db *gorm.DB,
) (*utils.ListResult[TEntity], error) {
	var items []TDataModel
	var totalItems int64

	// Crear una copia de la consulta para el conteo
	countQuery := db.Model(typemapper.GenericInstanceByT[TDataModel]())
	query := db.Model(typemapper.GenericInstanceByT[TDataModel]())

	// Aplicar filtros si existen
	if listQuery.Filters != nil && len(listQuery.Filters) > 0 {
		for _, filter := range listQuery.Filters {
			if filter.Field != "" && filter.Value != "" {
				switch strings.ToLower(filter.Comparison) {
				case "eq", "":
					query = query.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
					countQuery = countQuery.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
				case "neq":
					query = query.Where(fmt.Sprintf("%s != ?", filter.Field), filter.Value)
					countQuery = countQuery.Where(fmt.Sprintf("%s != ?", filter.Field), filter.Value)
				case "gt":
					query = query.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
					countQuery = countQuery.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
				case "gte":
					query = query.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)
					countQuery = countQuery.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)
				case "lt":
					query = query.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
					countQuery = countQuery.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
				case "lte":
					query = query.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)
					countQuery = countQuery.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)
				case "like":
					query = query.Where(fmt.Sprintf("%s LIKE ?", filter.Field), "%"+filter.Value+"%")
					countQuery = countQuery.Where(fmt.Sprintf("%s LIKE ?", filter.Field), "%"+filter.Value+"%")
				case "ilike":
					query = query.Where(fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", filter.Field), "%"+filter.Value+"%")
					countQuery = countQuery.Where(fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", filter.Field), "%"+filter.Value+"%")
				}
			}
		}
	}

	// Obtener el total de items
	err := countQuery.Count(&totalItems).Error
	if err != nil {
		return nil, err
	}

	// Aplicar ordenamiento si existe
	if listQuery.OrderBy != "" {
		query = query.Order(listQuery.OrderBy)
	}

	// Aplicar paginación
	if listQuery.Page > 0 {
		query = query.Offset(listQuery.GetOffset())
	}
	if listQuery.Size > 0 {
		query = query.Limit(listQuery.GetLimit())
	}

	// Ejecutar la consulta
	err = query.Find(&items).Error
	if err != nil {
		return nil, err
	}

	// Mapear los resultados al tipo de entidad
	entities, err := mapper.MapSlice[TDataModel, TEntity](items)
	if err != nil {
		return nil, fmt.Errorf("error mapping items to entities: %w", err)
	}

	// Crear el resultado paginado
	result := utils.NewListResult(
		entities,
		listQuery.GetSize(),
		listQuery.GetPage(),
		totalItems,
	)

	return result, nil
}
