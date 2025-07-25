package producer

import "go.opentelemetry.io/otel/attribute"

// ProducerTracingOptions is a struct that contains the options for the producer tracing
type ProducerTracingOptions struct {
	MessagingSystem string               // The messaging system used to send the message
	DestinationKind string               // The kind of destination of the message
	Destination     string               // The destination of the message
	OtherAttributes []attribute.KeyValue // Other attributes to be added to the message
}
