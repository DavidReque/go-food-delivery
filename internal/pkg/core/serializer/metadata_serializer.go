package serializer

import "github.com/DavidReque/go-food-delivery/internal/pkg/core/metadata"

// MetadataSerializer is an interface that defines the methods for serializing and deserializing metadata
type MetadataSerializer interface {
	// Serialize serializes a metadata object to a byte array
	Serialize(meta metadata.Metadata) ([]byte, error)
	// Deserialize deserializes a byte array to a metadata object
	Deserialize(bytes []byte) (metadata.Metadata, error)
}
