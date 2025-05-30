package contracts

import (
	"context"

	"gorm.io/gorm"
)

// GormDBContext define la interfaz para manejo de contexto y transacciones de base de datos
// Esta interfaz abstrae las operaciones de transacciones proporcionando métodos para:
// - Crear nuevas transacciones
// - Reutilizar transacciones existentes
// - Ejecutar múltiples operaciones en una sola transacción
// - Acceder a la instancia de base de datos subyacente
type GormDBContext interface {
	// WithTx crea un nuevo contexto de transacción
	// ctx: contexto de Go para timeouts y cancelaciones
	// Retorna: GormContext que encapsula la nueva transacción
	// Error: si no se puede iniciar la transacción
	// Uso: Para operaciones que requieren una nueva transacción explícita
	WithTx(ctx context.Context) (GormContext, error)

	// WithTxIfExists verifica si ya existe una transacción en el contexto
	// Si existe una transacción, la reutiliza; si no, usa la conexión normal
	// ctx: contexto que puede contener una transacción existente
	// Retorna: GormDBContext que puede ser la misma instancia o una nueva
	// Uso: Para operaciones que deben participar en transacciones existentes
	WithTxIfExists(ctx context.Context) GormDBContext

	// RunInTx ejecuta una función dentro de una transacción
	// ctx: contexto de Go para la operación
	// action: función que contiene las operaciones a ejecutar en la transacción
	// Retorna: error si la transacción falla o si action retorna error
	// La transacción se hace commit automáticamente si action retorna nil
	// La transacción se hace rollback automáticamente si action retorna error
	RunInTx(ctx context.Context, action ActionFunc) error

	// DB retorna la instancia subyacente de GORM
	// Permite acceso directo a la base de datos cuando se necesita
	// Retorna: *gorm.DB instancia de la conexión principal
	// Uso: Para operaciones que no requieren el contexto extendido
	DB() *gorm.DB
}

// ActionFunc define el tipo de función que se ejecuta dentro de una transacción
// ctx: contexto de Go que puede contener timeouts, valores, etc.
// gormContext: contexto de base de datos que encapsula la transacción
// Retorna: error si la operación falla (causará rollback de la transacción)

type ActionFunc func(ctx context.Context, gormContext GormDBContext) error
