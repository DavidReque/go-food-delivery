package specification

import (
	"fmt"
	"strings"
)

// Specification es la interfaz principal que define los métodos necesarios
// para construir consultas SQL dinámicas
type Specification interface {
	// GetQuery devuelve la parte SQL de la especificación
	GetQuery() string
	// GetValues devuelve los valores que se usarán como parámetros en la consulta
	GetValues() []any
}

// joinSpecification permite combinar múltiples especificaciones con un operador
// (AND u OR) entre ellas
type joinSpecification struct {
	specifications []Specification
	separator      string // El separador puede ser "AND" u "OR"
}

// GetQuery combina las consultas de todas las especificaciones usando el separador
func (s joinSpecification) GetQuery() string {
	queries := make([]string, 0, len(s.specifications))

	for _, spec := range s.specifications {
		queries = append(queries, spec.GetQuery())
	}

	return strings.Join(queries, fmt.Sprintf(" %s ", s.separator))
}

// GetValues recopila todos los valores de las especificaciones combinadas
func (s joinSpecification) GetValues() []any {
	values := make([]any, 0)

	for _, spec := range s.specifications {
		values = append(values, spec.GetValues()...)
	}

	return values
}

// And crea una especificación que combina múltiples especificaciones con el operador AND
func And(specifications ...Specification) Specification {
	return joinSpecification{
		specifications: specifications,
		separator:      "AND",
	}
}

// Or crea una especificación que combina múltiples especificaciones con el operador OR
func Or(specifications ...Specification) Specification {
	return joinSpecification{
		specifications: specifications,
		separator:      "OR",
	}
}

// notSpecification representa una negación de otra especificación
type notSpecification struct {
	Specification
}

// GetQuery envuelve la consulta original en una negación SQL
func (s notSpecification) GetQuery() string {
	return fmt.Sprintf(" NOT (%s)", s.Specification.GetQuery())
}

// Not crea una especificación que niega otra especificación
func Not(specification Specification) Specification {
	return notSpecification{
		specification,
	}
}

// binaryOperatorSpecification representa una comparación entre un campo y un valor
// usando un operador binario (=, >, <, >=, <=)
type binaryOperatorSpecification[T any] struct {
	field    string // Nombre del campo en la base de datos
	operator string // Operador SQL (=, >, <, etc.)
	value    T      // Valor contra el que se compara
}

// GetQuery genera la parte SQL de la comparación
func (s binaryOperatorSpecification[T]) GetQuery() string {
	return fmt.Sprintf("%s %s ?", s.field, s.operator)
}

// GetValues devuelve el valor que se usará en la comparación
func (s binaryOperatorSpecification[T]) GetValues() []any {
	return []any{s.value}
}

// Equal crea una especificación para comparar si un campo es igual a un valor
func Equal[T any](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: "=",
		value:    value,
	}
}

// GreaterThan crea una especificación para comparar si un campo es mayor que un valor
func GreaterThan[T comparable](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: ">",
		value:    value,
	}
}

// GreaterOrEqual crea una especificación para comparar si un campo es mayor o igual que un valor
func GreaterOrEqual[T comparable](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: ">=",
		value:    value,
	}
}

// LessThan crea una especificación para comparar si un campo es menor que un valor
func LessThan[T comparable](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: "<",
		value:    value,
	}
}

// LessOrEqual crea una especificación para comparar si un campo es menor o igual que un valor
func LessOrEqual[T comparable](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: ">=",
		value:    value,
	}
}

// stringSpecification es una especificación simple que representa una consulta SQL literal
type stringSpecification string

func (s stringSpecification) GetQuery() string {
	return string(s)
}

func (s stringSpecification) GetValues() []any {
	return nil
}

// IsNull crea una especificación para verificar si un campo es NULL
func IsNull(field string) Specification {
	return stringSpecification(fmt.Sprintf("%s IS NULL", field))
}
