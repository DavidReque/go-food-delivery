package gormdbcontext

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/logger/defaultlogger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/contracts"
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/helpers/gormextensions"

	"gorm.io/gorm"
)

type gormDBContext struct {
	db *gorm.DB
}

func NewGormDBContext(db *gorm.DB) contracts.GormDBContext {
	c := &gormDBContext{db: db}

	return c
}

func (c *gormDBContext) DB() *gorm.DB {
	return c.db
}

// WithTx crea un nuevo contexto de transacción
// ctx: contexto de Go para timeouts y cancelaciones
// Retorna: GormDBContext que encapsula la nueva transacción
// Error: si no se puede iniciar la transacción
// Uso: Para operaciones que requieren una nueva transacción explícita
func (c *gormDBContext) WithTx(
	ctx context.Context,
) (contracts.GormContext, error) {
	tx, err := gormextensions.GetTxFromContext(ctx)
	if err != nil {
		return contracts.GormContext{}, err
	}

	return contracts.GormContext{
		Tx:      tx,
		Context: ctx,
	}, nil
}

// WithTxIfExists verifica si ya existe una transacción en el contexto
// Si existe una transacción, la reutiliza; si no, usa la conexión normal
// ctx: contexto que puede contener una transacción existente
// Retorna: GormDBContext que puede ser la misma instancia o una nueva
// Uso: Para operaciones que deben participar en transacciones existentes
func (c *gormDBContext) WithTxIfExists(
	ctx context.Context,
) contracts.GormDBContext {
	tx := gormextensions.GetTxFromContextIfExists(ctx)
	if tx == nil {
		return c
	}

	return NewGormDBContext(tx)
}

// RunInTx ejecuta una función dentro de una transacción
// ctx: contexto de Go para la operación
// action: función que contiene las operaciones a ejecutar en la transacción
// Retorna: error si la transacción falla o si action retorna error
// La transacción se hace commit automáticamente si action retorna nil
// La transacción se hace rollback automáticamente si action retorna error
func (c *gormDBContext) RunInTx(
	ctx context.Context,
	action contracts.ActionFunc,
) error {
	// https://gorm.io/docs/transactions.html#Transaction
	tx := c.DB().WithContext(ctx).Begin()

	defaultlogger.GetLogger().Info("beginning database transaction")

	gormContext := gormextensions.SetTxToContext(ctx, tx)
	ctx = gormContext

	defer func() {
		if r := recover(); r != nil {
			tx.WithContext(ctx).Rollback()

			if err, _ := r.(error); err != nil {
				defaultlogger.GetLogger().Errorf(
					"panic tn the transaction, rolling back transaction with panic err: %+v",
					err,
				)
			} else {
				defaultlogger.GetLogger().Errorf("panic tn the transaction, rolling back transaction with panic message: %+v", r)
			}
		}
	}()

	err := action(ctx, c)
	if err != nil {
		defaultlogger.GetLogger().Error("rolling back transaction")
		tx.WithContext(ctx).Rollback()

		return err
	}

	defaultlogger.GetLogger().Info("committing transaction")

	if err = tx.WithContext(ctx).Commit().Error; err != nil {
		defaultlogger.GetLogger().Errorf("transaction commit error: %+v", err)
	}

	return err
}
