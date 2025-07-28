package options

import "github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/types"

type RabbitMQExchangeOptions struct {
	Name       string
	Type       types.ExchangeType
	AutoDelete bool           // if true, the exchange will be deleted when the last queue is unbound from it
	Durable    bool           // if true, the exchange will survive a broker restart
	Args       map[string]any // additional properties for the exchange
}
