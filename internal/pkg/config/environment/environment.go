package environment

import (
	"log"
	"os"
)

// GetEnv obtiene el valor de una variable de entorno de manera segura, con un valor por defecto.
// key: nombre de la variable de entorno que queremos obtener
// fallback: valor por defecto que se retornará si la variable no existe
func GetEnv(key, fallback string) string {
    if v, found := os.LookupEnv(key); found {
        return v
    }
    log.Printf("ADVERTENCIA: la variable de entorno %q no está definida, usando valor por defecto %q", key, fallback)
    return fallback
}

// RequireEnv obtiene una variable de entorno obligatoria.
// Si la variable no está definida, muestra un error crítico y termina la ejecución del programa.
func RequireEnv(key string) string {
    v := os.Getenv(key)
    if v == "" {
        log.Fatalf("ERROR: la variable de entorno %q es obligatoria pero no está definida", key)
    }
    return v
}