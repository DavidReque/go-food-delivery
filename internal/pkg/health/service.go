package health

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/health/contracts"
)

// healthService is the service that checks the health of the application
type healthService struct {
	healthParams contracts.HealthParams
}

// NewHealthService creates a new health service
func NewHealthService(
	healthParams contracts.HealthParams,
) contracts.HealthService {
	return &healthService{
		healthParams: healthParams,
	}
}

// CheckHealth checks the health of the application
func (service *healthService) CheckHealth(ctx context.Context) contracts.Check {
	checks := make(contracts.Check)

	for _, health := range service.healthParams.Healths {
		checks[health.GetHealthName()] = contracts.NewStatus(
			health.CheckHealth(ctx),
		)
	}

	return checks
}
