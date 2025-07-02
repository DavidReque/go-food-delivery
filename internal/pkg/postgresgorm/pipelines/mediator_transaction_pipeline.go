package pipelines

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/cqrs"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/helpers/gormextensions"
	"github.com/DavidReque/go-food-delivery/internal/pkg/reflection/typemapper"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

type mediatorTransactionPipeline struct {
	logger logger.Logger
	db     *gorm.DB
}

func NewMediatorTransactionPipeline(
	l logger.Logger,
	db *gorm.DB,
) mediatr.PipelineBehavior {
	return &mediatorTransactionPipeline{
		logger: l,
		db:     db,
	}
}

func (m *mediatorTransactionPipeline) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {
	requestName := typemapper.GetSnakeTypeName(request)

	// Si el request no es un TxRequest, se ejecuta el siguiente pipeline
	_, ok := request.(cqrs.TxRequest)
	if !ok {
		return next(ctx)
	}

	var result interface{}

	// https://gorm.io/docs/transactions.html#Transaction
	// --- Inicio de la Transacción ---
	tx := m.db.WithContext(ctx).Begin()

	m.logger.Infof(
		"beginning database transaction for request %s",
		requestName,
	)

	// --- Inyección de Dependencia por Contexto ---
	// Se inyecta el objeto de transacción 'tx' en el contexto.
	// Esto permite que los Handlers o Repositorios que se ejecuten después en la cadena
	// puedan extraer y utilizar esta misma transacción, en lugar de la conexión a la DB principal.
	gormContext := gormextensions.SetTxToContext(ctx, tx)
	ctx = gormContext

	defer func() {
		if r := recover(); r != nil {
			// Si se captura un pánico, se hace un ROLLBACK inmediato para evitar
			// que se realicen cambios en la base de datos.
			tx.WithContext(ctx).Rollback()

			if err, _ := r.(error); err != nil {
				m.logger.Errorf(
					"panic in the transaction, rolling back transaction with with panic err: %+v",
					err,
				)
			} else {
				m.logger.Errorf("panic tn the transaction, rolling back transaction with panic message: %+v", r)
			}
		}
	}()

	// Se invoca al siguiente pipeline (o al handler final) con el nuevo contexto que contiene la transacción.
	middlewareResponse, err := next(ctx)
	result = middlewareResponse

	if err != nil {
		m.logger.Errorf(
			"rolling back transaction for request %s with error: %+v",
			requestName,
		)
		tx.WithContext(ctx).Rollback()

		return nil, err
	}

	// --- Commit de la Transacción ---
	m.logger.Infof("committing transaction for request %s", requestName)

	if err := tx.WithContext(ctx).Commit().Error; err != nil {
		m.logger.Errorf("transaction commit error: ", err)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}
