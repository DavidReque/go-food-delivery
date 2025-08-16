package dbcontext

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/contracts"
	"github.com/DavidReque/go-food-delivery/internal/pkg/postgresgorm/gormdbcontext"
	"gorm.io/gorm"
	//"gorm.io/gorm"
)

type CatalogsGormDBContext struct {
	// nuestro dbcontext
	contracts.GormDBContext
}

func NewCatalogsDBContext(db *gorm.DB) *CatalogsGormDBContext {
	// initialize base GormContext
	c := &CatalogsGormDBContext{GormDBContext: gormdbcontext.NewGormDBContext(db)}

	return c
}
