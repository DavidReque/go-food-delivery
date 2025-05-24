package types

import (
	"time"
)

type IMessage interface {
	GeMessageId() string
	GetCreated() time.Time
	GetMessageTypeName() string
	GetMessageFullTypeName() string
}

type Message struct {
	MessageId string    `json:"messageId,omitempty"`
	Created   time.Time `json:"created"`
	EventType string    `json:"eventType"`
	//isMessage bool
}

func NewMessage(messageId string) *Message {
	return &Message{MessageId: messageId, Created: time.Now()}
}

func NewMessageWithTypeName(messageId string, eventTypeName string) *Message {
	return &Message{MessageId: messageId, Created: time.Now(), EventType: eventTypeName}
}

func (m *Message) GeMessageId() string {
	return m.MessageId
}

func (m *Message) GetCreated() time.Time {
	return m.Created
}

/*func (m *Message) GetMessageTypeName() string {
	return typeMapper.GetTypeName(m)
}

func (m *Message) GetMessageFullTypeName() string {
	return typeMapper.GetFullTypeName(m)
}*/