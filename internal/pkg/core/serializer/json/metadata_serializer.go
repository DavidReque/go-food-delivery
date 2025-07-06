package json

import (
	"emperror.dev/errors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/metadata"
	"github.com/DavidReque/go-food-delivery/internal/pkg/core/serializer"
)

type DefaultMetadataJsonSerializer struct {
	serializer serializer.Serializer
}

func NewDefaultMetadataJsonSerializer(s serializer.Serializer) serializer.MetadataSerializer {
	return &DefaultMetadataJsonSerializer{serializer: s}
}

func (s *DefaultMetadataJsonSerializer) Serialize(meta metadata.Metadata) ([]byte, error) {
	if meta == nil {
		return nil, nil
	}

	marshal, err := s.serializer.Marshal(meta)
	if err != nil {
		return nil, errors.WrapIf(err, "failded to marshal metadata")
	}

	return marshal, nil
}

func (s *DefaultMetadataJsonSerializer) Deserialize(bytes []byte) (metadata.Metadata, error) {
	if bytes == nil {
		return nil, nil
	}

	var meta metadata.Metadata

	if err := s.serializer.Unmarshal(bytes, &meta); err != nil {
		return nil, errors.WrapIf(err, "failed to unmarshal metadata")
	}

	return meta, nil
}
