package environment

import (
	"fmt"
	"os"
)

// GetEnv obtiene el valor de una variable de entorno de manera segura, con un valor por defecto.
// key: nombre de la variable de entorno que queremos obtener
// fallback: valor por defecto que se retornará si la variable no existe
func GetEnv(key, fallback string) string {
    if v, found := os.LookupEnv(key); found {
        return v
    }
    return fallback
}

// Propósito: Obtener una variable de entorno que es obligatoria para el funcionamiento del programa.
func RequireEnv(key string) string {
    v := os.Getenv(key)
    if v == "" {
        fmt.Printf("ERROR: la variable de entorno %q es obligatoria pero no está definida\n", key)
        os.Exit(1)
    }
    return v
}