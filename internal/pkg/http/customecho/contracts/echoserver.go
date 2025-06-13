package contracts

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/http/customecho/config"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/labstack/echo/v4"
)

// EchoGHttpServer es la interfaz que define los métodos del servidor HTTP
type EchoHttpServer interface {
	RunHttpServer(configEcho ...func(echo *echo.Echo)) error
	GracefulShutdown(ctx context.Context) error
	ApplyVersioningFromHeader()
	GetEchoInstance() *echo.Echo
	Logger() logger.Logger
	Cfg() *config.EchoHttpOptions
	SetupDefaultMiddlewares()
	RouteBuilder() *RouteBuilder
	AddMiddlewares(middlewares ...echo.MiddlewareFunc)
	ConfigGroup(groupName string, groupFunc func(group *echo.Group))
}
