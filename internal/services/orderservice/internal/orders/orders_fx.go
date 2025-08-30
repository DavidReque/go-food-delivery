package orders

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/web/route"
	"github.com/DavidReque/go-food-delivery/internal/pkg/es"
	"github.com/DavidReque/go-food-delivery/internal/pkg/eventstroredb"
	echocontracts "github.com/DavidReque/go-food-delivery/internal/pkg/http/customecho/contracts"
	"github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/data/repositories"
	createOrderV1 "github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/features/creating_order/v1/endpoints"
	getOrderByIdV1 "github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/features/getting_order_by_id/v1/endpoints"
	getOrdersV1 "github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/features/getting_orders/v1/endpoints"
	"github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/models/orders/aggregate"
	"github.com/DavidReque/go-food-delivery/internal/services/orderservice/internal/orders/projections"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"ordersfx",

	// Other provides
	fx.Provide(fx.Annotate(repositories.NewMongoOrderReadRepository)),
	fx.Provide(repositories.NewElasticOrderReadRepository),

	fx.Provide(eventstroredb.NewEventStoreAggregateStore[*aggregate.Order]),
	fx.Provide(fx.Annotate(func(catalogsServer echocontracts.EchoHttpServer) *echo.Group {
		var g *echo.Group
		catalogsServer.RouteBuilder().RegisterGroupFunc("/api/v1", func(v1 *echo.Group) {
			group := v1.Group("/orders")
			g = group
		})

		return g
	}, fx.ResultTags(`name:"order-echo-group"`))),

	fx.Provide(
		route.AsRoute(createOrderV1.NewCreteOrderEndpoint, "order-routes"),
		route.AsRoute(getOrderByIdV1.NewGetOrderByIdEndpoint, "order-routes"),
		route.AsRoute(getOrdersV1.NewGetOrdersEndpoint, "order-routes"),
	),

	fx.Provide(
		es.AsProjection(projections.NewElasticOrderProjection),
		es.AsProjection(projections.NewMongoOrderProjection),
	),
)
