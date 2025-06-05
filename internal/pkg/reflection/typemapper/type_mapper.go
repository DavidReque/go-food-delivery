package typemapper

import (
	"reflect"
	"strings"
)

// GetTypeName devuelve el nombre del tipo sin el paquete
func GetTypeName(obj interface{}) string {
	if obj == nil {
		return ""
	}
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

// GetFullTypeName devuelve el nombre completo del tipo incluyendo el paquete
func GetFullTypeName(obj interface{}) string {
	if obj == nil {
		return ""
	}
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return strings.Replace(t.String(), ".", "_", -1)
}
