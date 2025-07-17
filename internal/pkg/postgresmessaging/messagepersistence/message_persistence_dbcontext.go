package messagepersistence

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/contracts"
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/gormdbcontext"

	"gorm.io/gorm"
)

type PostgresMessagePersistenceDBContext struct {
	// our db context base
	contracts.GormDBContext
}

func NewPostgresMessagePersistenceDBContext(
	db *gorm.DB,
) *PostgresMessagePersistenceDBContext {
	// initialize base GormContext
	c := &PostgresMessagePersistenceDBContext{GormDBContext: gormdbcontext.NewGormDBContext(db)}

	return c
}
