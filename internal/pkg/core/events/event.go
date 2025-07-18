package events

import (
	"time"

	"github.com/DavidReque/go-food-delivery/internal/pkg/reflection/typemapper"
	uuid "github.com/satori/go.uuid"
)

type IEvent interface {
	GetEventId() uuid.UUID
	GetOcurredOn() time.Time
	GetEventTypeName() string
	GetEventFullTypeName() string
}

type Event struct {
	EventId    uuid.UUID `json:"event_id"`
	EventType  string    `json:"event_type"`
	OccurredOn time.Time `json:"ocurred_on"`
}

func NewEvent(eventType string) *Event {
	return &Event{
		EventId:    uuid.NewV4(),
		OccurredOn: time.Now(),
		EventType:  eventType,
	}
}

func (e *Event) GetEventId() uuid.UUID {
	return e.EventId
}

func (e *Event) GetEventType() string {
	return e.EventType
}

func (e *Event) GetOccurredOn() time.Time {
	return e.OccurredOn
}

func (e *Event) GetEventTypeName() string {
	return typemapper.GetTypeName(e)
}

func (e *Event) GetEventFullTypeName() string {
	return typemapper.GetFullTypeName(e)
}

func IsEvent(obj interface{}) bool {
	if _, ok := obj.(IEvent); ok {
		return true
	}

	return false
}
