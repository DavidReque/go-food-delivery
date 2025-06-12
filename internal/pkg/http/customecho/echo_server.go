package customecho

import (
	"context"
	"net/http"

	"github.com/DavidReque/go-food-delivery/internal/pkg/http/customecho/config"
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
}

type echoServer struct {
	echo   *echo.Echo
	config *config.EchoHttpOptions
}

// NewEchoServer crea una nueva instancia del servidor
func NewEchoServer(cfg *config.EchoHttpOptions) EchoServer {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	e := echo.New()
	e.HideBanner = true

	server := &echoServer{
		echo:   e,
		config: cfg,
	}

	server.setupMiddlewares()
	return server
}

// setupMiddlewares configura los middlewares del servidor
func (s *echoServer) setupMiddlewares() {
	// Recuperación de pánico
	s.echo.Use(middleware.Recover())

	// Logger
	s.echo.Use(middleware.Logger())

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

func (s *echoServer) Start() error {
	s.echo.Server.ReadTimeout = s.config.ReadTimeout
	s.echo.Server.WriteTimeout = s.config.WriteTimeout
	s.echo.Server.MaxHeaderBytes = s.config.MaxHeaderBytes

	return s.echo.Start(s.config.Port)
}

func (s *echoServer) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}

func (s *echoServer) GetInstance() *echo.Echo {
	return s.echo
}
