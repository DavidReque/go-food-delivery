package contracts

import (
	"context"

	"gorm.io/gorm"
)

// GormContext combina una transacción de GORM con un contexto de Go
// Esta estructura permite pasar tanto la conexión a la base de datos como el contexto
// en una sola entidad, facilitando el manejo de transacciones y timeouts
//
// Casos de uso típicos:
// - Operaciones que requieren transacciones (múltiples inserts/updates)
type GormContext struct {
	// Ejemplo: gormCtx.Tx.Create(&product)
	Tx *gorm.DB

	context.Context
}
