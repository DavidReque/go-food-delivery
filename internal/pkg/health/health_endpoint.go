package health

import (
	"net/http"

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
}

// checkHealth is the endpoint that checks the health of the application
func (s *HealthCheckEndpoint) checkHealth(c echo.Context) error {
	check := s.service.CheckHealth(c.Request().Context())
	if !check.AllUp() {
		return c.JSON(http.StatusServiceUnavailable, check)
	}

	return c.JSON(http.StatusOK, check)
}
