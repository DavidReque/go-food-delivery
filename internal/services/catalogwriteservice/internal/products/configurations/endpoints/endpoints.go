package endpoints

import "github.com/DavidReque/go-food-delivery/internal/pkg/core/web/route"

func RegisterEndpoints(endpoints []route.Endpoint) error {
	for _, endpoint := range endpoints {
		endpoint.MapEndpoint()
	}

	return nil
}
