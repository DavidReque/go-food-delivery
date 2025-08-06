package bus

import (
	consumer2 "github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/consumer"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/messaging/producer"
)

type Bus interface {
	producer.Producer
	consumer2.BusControl
	consumer2.ConsumerConnector
}
