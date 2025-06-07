package attribute

import (
	"encoding/json"

	"go.opentelemetry.io/otel/attribute"
)

// Object creates a KeyValue with a interface{} Value type.
func Object(k string, v interface{}) attribute.KeyValue {
	marshal, err := json.Marshal(&v)
	if err != nil {
		return attribute.KeyValue{}
	}
	return attribute.KeyValue{Key: attribute.Key(k), Value: attribute.StringValue(string(marshal))}
}
