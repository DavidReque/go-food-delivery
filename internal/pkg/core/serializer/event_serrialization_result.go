package serializer

// EventSerializationResult is a struct that contains the data and content type of an event
type EventSerializationResult struct {
	Data        []byte
	ContentType string
}
