package contracts

import "github.com/labstack/echo/v4"

// RouteBuilder es un constructor que facilita la configuración de rutas en Echo
// Implementa el patrón Builder para hacer más fluida la configuración del router
type RouteBuilder struct {
	echo *echo.Echo
}

// NewRouteBuilder crea una nueva instancia de RouteBuilder
// Recibe una instancia de echo.Echo como dependencia
func NewRouteBuilder(echo *echo.Echo) *RouteBuilder {
	return &RouteBuilder{echo: echo}
}

// RegisterRoutes permite registrar rutas directamente en la instancia de Echo
// Recibe una función builder que configura las rutas
// Retorna el builder para permitir encadenamiento de métodos
func (r *RouteBuilder) RegisterRoutes(builder func(e *echo.Echo)) *RouteBuilder {
	builder(r.echo)
	return r
}

// RegisterGroupFunc crea un nuevo grupo de rutas con un prefijo común
// Recibe el nombre del grupo y una función builder para configurar las rutas del grupo
// Retorna el builder para permitir encadenamiento de métodos
func (r *RouteBuilder) RegisterGroupFunc(groupName string, builder func(g *echo.Group)) *RouteBuilder {
	builder(r.echo.Group(groupName))
	return r
}

// RegisterGroup crea un nuevo grupo de rutas con un prefijo común
// Similar a RegisterGroupFunc pero sin la función de configuración
// Retorna el builder para permitir encadenamiento de métodos
func (r *RouteBuilder) RegisterGroup(groupName string) *RouteBuilder {
	r.echo.Group(groupName)
	return r
}

// Build finaliza la construcción y retorna la instancia configurada de Echo
// Este método debe llamarse al final de la configuración
func (r *RouteBuilder) Build() *echo.Echo {
	return r.echo
}
