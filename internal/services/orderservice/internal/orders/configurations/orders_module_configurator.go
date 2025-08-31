package configurations

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/web/route"
	"github.com/DavidReque/go-food-delivery/internal/pkg/es/contracts/store"
	contracts2 "github.com/DavidReque/go-food-delivery/internal/pkg/fxapp/contracts"
	grpcServer "github.com/DavidReque/go-food-delivery/internal/pkg/grpc"
	echocontracts "github.com/DavidReque/go-food-delivery/internal/pkg/http/customecho/contracts"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing"
	"github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/configurations/mappings"
	"github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/configurations/mediatr"
	"github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/contracts/repositories"
	"github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/models/orders/aggregate"
	"github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/shared/contracts"
	"github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/shared/grpc"
	ordersservice "github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/shared/grpc/genproto"

	"github.com/go-playground/validator/v10"
	googleGrpc "google.golang.org/grpc"
)

type OrdersModuleConfigurator struct {
	contracts2.Application
}

func NewOrdersModuleConfigurator(
	app contracts2.Application,
) *OrdersModuleConfigurator {
	return &OrdersModuleConfigurator{
		Application: app,
	}
}

func (c *OrdersModuleConfigurator) ConfigureOrdersModule() {
	c.ResolveFunc(
		func(logger logger.Logger,
			server echocontracts.EchoHttpServer,
			orderRepository repositories.OrderMongoRepository,
			orderAggregateStore store.AggregateStore[*aggregate.Order],
			tracer tracing.AppTracer,
		) error {
			// config Orders Mappings
			err := mappings.ConfigureOrdersMappings()
			if err != nil {
				return err
			}

			// config Orders Mediators
			err = mediatr.ConfigOrdersMediator(logger, orderRepository, orderAggregateStore, tracer)
			if err != nil {
				return err
			}

			return nil
		},
	)
}

func (c *OrdersModuleConfigurator) MapOrdersEndpoints() {
	// config Orders Http Endpoints
	c.ResolveFuncWithParamTag(func(endpoints []route.Endpoint) {
		for _, endpoint := range endpoints {
			endpoint.MapEndpoint()
		}
	}, `group:"order-routes"`,
	)

	// config Orders Grpc Endpoints
	c.ResolveFunc(
		func(ordersGrpcServer grpcServer.GrpcServer, ordersMetrics *contracts.OrdersMetrics, logger logger.Logger, validator *validator.Validate) error {
			orderGrpcService := grpc.NewOrderGrpcService(logger, validator, ordersMetrics)
			ordersGrpcServer.GrpcServiceBuilder().RegisterRoutes(func(server *googleGrpc.Server) {
				ordersservice.RegisterOrdersServiceServer(server, orderGrpcService)
			})
			return nil
		},
	)
}
