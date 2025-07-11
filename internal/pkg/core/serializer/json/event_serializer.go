package json

import (
	"reflect"

	"github.com/DavidReque/go-food-delivery/internal/pkg/core/domain"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/serializer"
	typeMapper "github.com/DavidReque/go-food-delivery/internal/pkg/reflection/typemapper"

	"emperror.dev/errors"
)

type DefaultEventJsonSerializer struct {
	serializer serializer.Serializer
}

func NewDefaultEventJsonSerializer(serializer serializer.Serializer) serializer.EventSerializer {
	return &DefaultEventJsonSerializer{serializer: serializer}
}

func (s *DefaultEventJsonSerializer) Serialize(event domain.IDomainEvent) (*serializer.EventSerializationResult, error) {
	return s.SerializeObject(event)
}

func (s *DefaultEventJsonSerializer) SerializeObject(event interface{}) (*serializer.EventSerializationResult, error) {
	if event == nil {
		return &serializer.EventSerializationResult{Data: nil, ContentType: s.ContentType()}, nil
	}

	// we use event short type name instead of full type name because this event in other receiver packages could have different package name
	eventType := typeMapper.GetTypeName(event)

	data, err := s.serializer.Marshal(event)
	if err != nil {
		return nil, errors.WrapIff(err, "error in Marshaling: `%s`", eventType)
	}

	result := &serializer.EventSerializationResult{Data: data, ContentType: s.ContentType()}

	return result, nil
}

func (s *DefaultEventJsonSerializer) ContentType() string {
	return "application/json"
}

func (s *DefaultEventJsonSerializer) Deserialize(
	data []byte,
	eventType string,
	contentType string,
) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	targetEventPointer := typeMapper.EmptyInstanceByTypeNameAndImplementedInterface[domain.IDomainEvent](
		eventType,
	)

	if targetEventPointer == nil {
		return nil, errors.Errorf("event type `%s` is not impelemted IDomainEvent or can't be instansiated", eventType)
	}

	if contentType != s.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := s.serializer.Unmarshal(data, targetEventPointer); err != nil {
		return nil, errors.WrapIff(err, "error in Unmarshaling: `%s`", eventType)
	}

	return targetEventPointer, nil
}

func (s *DefaultEventJsonSerializer) DeserializeObject(
	data []byte,
	eventType string,
	contentType string,
) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	targetEventPointer := typeMapper.InstanceByTypeName(eventType)

	if targetEventPointer == nil {
		return nil, errors.Errorf("event type `%s` can't be instansiated", eventType)
	}

	if contentType != s.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := s.serializer.Unmarshal(data, targetEventPointer); err != nil {
		return nil, errors.WrapIff(err, "error in Unmarshaling: `%s`", eventType)
	}

	return targetEventPointer, nil
}

func (s *DefaultEventJsonSerializer) DeserializeType(
	data []byte,
	eventType reflect.Type,
	contentType string,
) (domain.IDomainEvent, error) {
	if data == nil {
		return nil, nil
	}

	// we use event short type name instead of full type name because this event in other receiver packages could have different package name
	eventTypeName := typeMapper.GetTypeName(eventType)

	result, err := s.Deserialize(data, eventTypeName, contentType)
	if err != nil {
		return nil, err
	}

	event, ok := result.(domain.IDomainEvent)
	if !ok {
		return nil, errors.New("result is not a domain event")
	}

	return event, nil
}

func (s *DefaultEventJsonSerializer) Serializer() serializer.Serializer {
	return s.serializer
}
