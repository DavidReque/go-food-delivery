package serializer

import (
	"reflect"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/domain"
)

type EventSerializer interface {
	// Serialize a domain event to a byte array
	Serialize(event domain.IDomainEvent) (*EventSerializationResult, error) 
	// Serialize an object to a byte array
	SerializeObject(event interface{}) (*EventSerializationResult, error)
	// Deserialize a byte array to a domain event
	Deserialize(data []byte, eventType string, contentType string) (interface{}, error)
	// Deserialize a byte array to an object
	DeserializeType(data []byte, eventType reflect.Type, contentType string) (domain.IDomainEvent, error)
	// ContentType returns the content type of the event serializer
	ContentType() string
	// Serializer returns the serializer used by the event serializer
	Serializer() Serializer
}