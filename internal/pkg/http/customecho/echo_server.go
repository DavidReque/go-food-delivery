package customecho

import (
	"context"
	"net/http"
	"time"

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

// Config contiene la configuración del servidor
type Config struct {
	// Port es el puerto en el que se iniciará el servidor
	Port string
	// ReadTimeout es el tiempo máximo de espera para leer las peticiones
	ReadTimeout time.Duration
	// WriteTimeout es el tiempo máximo de espera para escribir las respuestas
	WriteTimeout time.Duration
	// MaxHeaderBytes es el tamaño máximo de los encabezados de las peticiones
	MaxHeaderBytes int
	// BodyLimit es el tamaño máximo del cuerpo de las peticiones
	BodyLimit string
	// EnableGzip habilita la compresión gzip
	EnableGzip bool
	// EnableRateLimit habilita el límite de tasa de las peticiones
	EnableRateLimit bool
}

// DefaultConfig retorna una configuración por defecto
func DefaultConfig() *Config {

	return &Config{
		Port:            ":8080",
		ReadTimeout:     time.Second * 30,
		WriteTimeout:    time.Second * 30,
		MaxHeaderBytes:  1 << 20, // 1 MB
		BodyLimit:       "2M",
		EnableGzip:      true,
		EnableRateLimit: true,
	}
}

type echoServer struct {
	echo   *echo.Echo
	config *Config
}

// NewEchoServer crea una nueva instancia del servidor
func NewEchoServer(config *Config) EchoServer {
	if config == nil {
		config = DefaultConfig()
	}

	// Creación de la instancia del servidor
	e := echo.New()
	e.HideBanner = true

	server := &echoServer{
		echo:   e,
		config: config,
	}

	// Configuración de los middlewares
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
		// Permite todas las origenes
		AllowOrigins: []string{"*"},
		// Permite los métodos HTTP permitidos
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
