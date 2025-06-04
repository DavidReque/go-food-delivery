package gormdbcontext

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/contracts"
)

// AddModel es una función genérica para crear registros en la base de datos
// TDataModel: tipo del modelo de datos
// TDomain: tipo del modelo de dominio
func AddModel[TDataModel any, TDomain any](
	ctx context.Context,
	dbContext contracts.GormDBContext,
	model TDomain,
) (TDomain, error) {
	if err := dbContext.DB().WithContext(ctx).Create(model).Error; err != nil {
		return model, err
	}
	return model, nil
}
