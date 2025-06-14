package environment

import (
	"fmt"
	"net/url"
	"strconv"
)

// MustAtoi convierte un string a entero y hace panic si falla.
// s: valor string a convertir.
// name: nombre de la variable (para mensajes de error descriptivos).
func MustAtoi(s, name string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("valor %q para %s no es un entero válido: %v", s, name, err))
	}
	return i
}

// ValidateURL comprueba que el string es una URL válida (tiene esquema y host).
// raw: string a validar como URL.
// name: nombre de la variable (para mensajes de error descriptivos).
func ValidateURL(raw, name string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		panic(fmt.Sprintf("valor %q para %s no es una URL válida", raw, name))
	}
	return raw
}

// Explicación:
// - MustAtoi asegura que los valores numéricos requeridos en la configuración sean válidos y detiene la ejecución si no lo son.
// - ValidateURL garantiza que las URLs requeridas sean válidas antes de que la aplicación continúe.
// - Ambas funciones ayudan a detectar errores de configuración temprano y con mensajes claros.
