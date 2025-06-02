package gormextensions

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/constants"
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/contracts"

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
