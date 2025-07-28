package consumercontracts

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/consumer"
	messagingTypes "github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/types"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/consumer/configurations"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/types"
)

// ConsumerFactory is a factory for creating consumers
type ConsumerFactory interface {
	CreateConsumer(
		consumerConfiguration *configurations.RabbitMQConsumerConfiguration,
		isConsumedNotifications ...func(message messagingTypes.IMessage),
	) (consumer.Consumer, error)

	Connection() types.IConnection
}
