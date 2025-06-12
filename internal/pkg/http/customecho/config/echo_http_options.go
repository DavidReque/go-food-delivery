package config

import (
	"fmt"
	"net/url"
	"time"
	//"github.com/iancoleman/strcase"
)

//var optionName = strcase.ToLowerCamel("EchoHttpOptions")
// EchoHttpOptions contiene la configuración del servidor
type EchoHttpOptions struct {
	// Port es el puerto en el que se iniciará el servidor
	Port string `mapstructure:"port" validate:"required" env:"Port"`
	// ReadTimeout es el tiempo máximo de espera para leer las peticiones
	ReadTimeout time.Duration `mapstructure:"readTimeout" validate:"required" env:"ReadTimeout"`
	// WriteTimeout es el tiempo máximo de espera para escribir las respuestas
	WriteTimeout time.Duration `mapstructure:"writeTimeout" validate:"required" env:"WriteTimeout"`
	// MaxHeaderBytes es el tamaño máximo de los encabezados de las peticiones
	MaxHeaderBytes int `mapstructure:"maxHeaderBytes" env:"MaxHeaderBytes"`
	// BodyLimit es el tamaño máximo del cuerpo de las peticiones
	BodyLimit string `mapstructure:"bodyLimit" env:"BodyLimit"`
	// EnableGzip habilita la compresión gzip
	EnableGzip bool `mapstructure:"enableGzip" env:"EnableGzip"`
	// EnableRateLimit habilita el límite de tasa de las peticiones
	EnableRateLimit bool `mapstructure:"enableRateLimit" env:"EnableRateLimit"`
	// Development indica si el servidor está en modo desarrollo
	Development bool `mapstructure:"development" env:"Development"`
	// BasePath es la ruta base para todas las peticiones
	BasePath string `mapstructure:"basePath" validate:"required" env:"BasePath"`
	// Host es el host del servidor
	Host string `mapstructure:"host" env:"Host"`
	// Name es el nombre del servicio
	Name string `mapstructure:"name" env:"ServiceName"`
}

// DefaultConfig retorna una configuración por defecto
func DefaultConfig() *EchoHttpOptions {
	return &EchoHttpOptions{
		Port:            ":8080",
		ReadTimeout:     time.Second * 30,
		WriteTimeout:    time.Second * 30,
		MaxHeaderBytes:  1 << 20, // 1 MB
		BodyLimit:       "2M",
		EnableGzip:      true,
		EnableRateLimit: true,
		Development:     false,
		BasePath:        "/api",
		Host:            "localhost",
		Name:            "echo-server",
	}
}

// Address retorna la dirección completa del servidor
func (c *EchoHttpOptions) Address() string {
	return fmt.Sprintf("%s%s", c.Host, c.Port)
}

// BasePathAddress retorna la dirección completa incluyendo la ruta base
func (c *EchoHttpOptions) BasePathAddress() string {
	path, err := url.JoinPath(c.Address(), c.BasePath)
	if err != nil {
		return ""
	}
	return path
}
