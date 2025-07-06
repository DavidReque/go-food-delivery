package serializer

type Serializer interface {
	// Marshal serializes an object to a byte array
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal deserializes a byte array to an object
	Unmarshal(data []byte, v interface{}) error
	// UnmarshalFromJson deserializes a JSON string to an object
	UnmarshalFromJson(data string, v interface{}) error
	// DecodeWithMapStructure decodes a byte array to a map structure
	DecodeWithMapStructure(
		input interface{},
		output interface{},
	) error
	// UnmarshalToMap deserializes a byte array to a map structure
	UnmarshalToMap(data []byte, v *map[string]interface{}) error
	// UnmarshalToMapFromJson deserializes a JSON string to a map structure
	UnmarshalToMapFromJson(data string, v *map[string]interface{}) error
	// PrettyPrint pretty prints an object
	PrettyPrint(data interface{}) string
	// ColoredPrint pretty prints an object with colors
	ColoredPrint(data interface{}) string
}
