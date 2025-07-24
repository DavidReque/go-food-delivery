package options

import "github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/types"

type RabbitMQExchangeOptions struct {
	// name of the exchange
	Name string

	// type of the exchange (fanout, direct, topic)
	Type types.ExchangeType

	// if the exchange is deleted when the last consumer unsubscribes
	AutoDelete bool

	// if the exchange is durable, it will survive a server restart
	// if it is true, will survive a server restart
	Durable bool

	// arguments for the exchange, aditional arguments for the exchange
	Args map[string]any
}
