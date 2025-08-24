package contracts

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"go.uber.org/fx"
)

// ApplicationBuilder es una interfaz que define cómo se construye una aplicación.
// Proporciona métodos para registrar módulos, dependencias y configurar la aplicación.
type ApplicationBuilder interface {
	// AddModule agrega módulos a la aplicación (versión variádica).
	// Permite agregar múltiples módulos de una vez.
	// Parámetros:
	//   - module: Lista variable de opciones de fx que representan módulos
	AddModule(module ...fx.Option)

	// ProvideModule registra módulos directamente en la aplicación.
	// Los módulos son conjuntos de funcionalidades que se pueden agregar de forma independiente.
	// Parámetros:
	//   - module: Una opción de fx que representa un módulo completo
	ProvideModule(module fx.Option)

	// Provide registra constructores de funciones para resolver dependencias.
	// Estos constructores son utilizados por el contenedor de inyección de dependencias
	// para crear instancias de servicios cuando sean necesarios.
	// Parámetros:
	//   - constructors: Lista variable de constructores que crean servicios
	Provide(constructors ...interface{})

	// Decorate registra funciones que modifican o envuelven servicios existentes.
	// Útil para agregar funcionalidad adicional a servicios sin modificar su implementación original.
	// Parámetros:
	//   - constructors: Lista variable de funciones decoradoras
	Decorate(constructors ...interface{})

	// Build construye y devuelve la instancia final de la aplicación.
	// Este método debe llamarse después de registrar todos los módulos y dependencias.
	// Retorna:
	//   - Una instancia de Application lista para ejecutarse
	Build() Application

	// GetProvides devuelve la lista de todos los constructores registrados.
	// Útil para inspeccionar qué servicios están configurados en la aplicación.
	GetProvides() []interface{}

	// GetInvokes devuelve la lista de todos los invokes registrados.
	// Permite verificar qué funciones se ejecutarán al iniciar la aplicación.
	GetInvokes() []interface{}

	// GetDecorates devuelve la lista de todos los decoradores registrados.
	// Permite verificar qué modificaciones se aplicarán a los servicios.
	GetDecorates() []interface{}

	// Options devuelve la lista de opciones de configuración de fx.
	// Estas opciones determinan el comportamiento general de la aplicación.
	Options() []fx.Option

	// Logger proporciona acceso al sistema de registro de la aplicación.
	// Esencial para el seguimiento y depuración de la aplicación.
	Logger() logger.Logger

	// Environment da acceso a la configuración del entorno actual.
	// Permite que la aplicación se adapte según el entorno (desarrollo, producción, etc.).
	Environment() environment.Environment
}
