package health

import (
	"net/http"
	"context"

	contracts2 "github.com/DavidReque/go-food-delivery/internal/pkg/health/contracts"
	"github.com/DavidReque/go-food-delivery/internal/pkg/http/customecho/contracts"

	"github.com/labstack/echo/v4"
)

type HealthCheckEndpoint struct {
	service    contracts2.HealthService
	echoServer contracts.EchoHttpServer
}

func NewHealthCheckEndpoint(
	service contracts2.HealthService,
	server contracts.EchoHttpServer,
) *HealthCheckEndpoint {
	return &HealthCheckEndpoint{service: service, echoServer: server}
}

func (s *HealthCheckEndpoint) RegisterEndpoints() {
	s.echoServer.GetEchoInstance().GET("health", s.checkHealth)
	s.echoServer.GetEchoInstance().GET("test", s.testEndpoint) // test endpoint
}

// checkHealth is the endpoint that checks the health of the application
func (s *HealthCheckEndpoint) checkHealth(c echo.Context) error {
	// Use context.Background() instead of c.Request().Context() to avoid premature cancellation
	ctx := context.Background()
	check := s.service.CheckHealth(ctx)
	if !check.AllUp() {
		return c.JSON(http.StatusServiceUnavailable, check)
	}

	return c.JSON(http.StatusOK, check)
}

func (s *HealthCheckEndpoint) testEndpoint(c echo.Context) error {
	return c.String(http.StatusOK, "Test endpoint working!")
}
