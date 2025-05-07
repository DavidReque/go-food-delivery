package environment

import (
	"fmt"
	"net/url"
	"strconv"
)

// MustAtoi convierte un string a entero y panic si falla.
func MustAtoi(s, name string) int {
    i, err := strconv.Atoi(s)
    if err != nil {
        panic(fmt.Sprintf("valor %q para %s no es un entero válido: %v", s, name, err))
    }
    return i
}

// ValidateURL comprueba que el string es una URL válida (tiene esquema y host).
func ValidateURL(raw, name string) string {
    u, err := url.Parse(raw)
    if err != nil || u.Scheme == "" || u.Host == "" {
        panic(fmt.Sprintf("valor %q para %s no es una URL válida", raw, name))
    }
    return raw
}