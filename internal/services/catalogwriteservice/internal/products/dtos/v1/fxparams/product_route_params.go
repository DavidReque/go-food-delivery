package fxparams

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/shared/contracts"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type ProductRouteParams struct {
	fx.In

	CatalogsMetrics *contracts.CatalogsMetrics
	Logger          logger.Logger
	ProductsGroup   *echo.Group `name:"product-echo-group"`
	Validator       *validator.Validate
}
