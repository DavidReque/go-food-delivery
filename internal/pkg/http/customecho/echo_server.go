package customecho

import (
	"context"
	"net/http"

	"github.com/DavidReque/go-food-delivery/internal/pkg/http/customecho/config"
	"github.com/DavidReque/go-food-delivery/internal/pkg/http/customecho/contracts"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// EchoServer representa la interfaz principal del servidor
type EchoServer interface {
	// Start inicia el servidor
	Start() error
	// Shutdown detiene el servidor
	Shutdown(ctx context.Context) error
	// GetInstance devuelve la instancia del servidor
	GetInstance() *echo.Echo
	// ConfigGroup configura un grupo de rutas
	ConfigGroup(groupName string, groupFunc func(group *echo.Group))
	// AddMiddlewares añade middlewares al servidor
	AddMiddlewares(middlewares ...echo.MiddlewareFunc)
	// SetupDefaultMiddlewares configura los middlewares por defecto
	SetupDefaultMiddlewares()
	// Config devuelve la configuración del servidor
	Config() *config.EchoHttpOptions
}

type echoHttpServer struct {
	echo   *echo.Echo
	config *config.EchoHttpOptions
	logger logger.Logger
	//meter metric.Meter
	routeBuilder *contracts.RouteBuilder
}

// NewEchoServer crea una nueva instancia del servidor
func NewEchoHttpServer(cfg *config.EchoHttpOptions, log logger.Logger) contracts.EchoHttpServer {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	e := echo.New()
	e.HideBanner = true

	server := &echoHttpServer{
		echo:   e,
		config: cfg,
		logger: log,
	}

	return server
}

func (s *echoHttpServer) ConfigGroup(groupName string, groupFunc func(group *echo.Group)) {
	groupFunc(s.echo.Group(groupName))
}

func (s *echoHttpServer) AddMiddlewares(middlewares ...echo.MiddlewareFunc) {
	if len(middlewares) > 0 {
		s.echo.Use(middlewares...)
	}
}

func (s *echoHttpServer) Config() *config.EchoHttpOptions {
	return s.config
}

// SetupDefaultMiddlewares configura los middlewares del servidor
func (s *echoHttpServer) SetupDefaultMiddlewares() {
	// Recuperación de pánico
	s.echo.Use(middleware.Recover())

	// Logger - disabled temporarily 
	// s.echo.Use(middleware.Logger())

	// CORS
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Límite de tamaño del body
	s.echo.Use(middleware.BodyLimit(s.config.BodyLimit))

	// Request ID
	s.echo.Use(middleware.RequestID())

	// Compresión gzip
	if s.config.EnableGzip {
		s.echo.Use(middleware.Gzip())
	}

	// Límite de tasa de las peticiones
	if s.config.EnableRateLimit {
		s.echo.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	}
}

func (s *echoHttpServer) Start() error {
	s.echo.Server.ReadTimeout = s.config.ReadTimeout
	s.echo.Server.WriteTimeout = s.config.WriteTimeout
	s.echo.Server.MaxHeaderBytes = s.config.MaxHeaderBytes

	return s.echo.Start(s.config.Port)
}

func (s *echoHttpServer) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}

func (s *echoHttpServer) GetInstance() *echo.Echo {
	return s.echo
}

// Implementación de los métodos de contracts.EchoHttpServer

func (s *echoHttpServer) ApplyVersioningFromHeader() {
	// TODO: Implementar versionado desde headers
}

func (s *echoHttpServer) RunHttpServer(configEcho ...func(echo *echo.Echo)) error {
	// Aplicar configuraciones opcionales
	for _, config := range configEcho {
		config(s.echo)
	}
	
	// Usar el método Start existente
	return s.Start()
}

func (s *echoHttpServer) GracefulShutdown(ctx context.Context) error {
	// Usar el método Shutdown existente
	return s.Shutdown(ctx)
}

func (s *echoHttpServer) GetEchoInstance() *echo.Echo {
	// Usar el método GetInstance existente
	return s.GetInstance()
}

func (s *echoHttpServer) Logger() logger.Logger {
	return s.logger
}

func (s *echoHttpServer) Cfg() *config.EchoHttpOptions {
	// Usar el método Config existente
	return s.Config()
}

func (s *echoHttpServer) RouteBuilder() *contracts.RouteBuilder {
	if s.routeBuilder == nil {
		s.routeBuilder = contracts.NewRouteBuilder(s.echo)
	}
	return s.routeBuilder
}
