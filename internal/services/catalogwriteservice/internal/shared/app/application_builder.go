package app

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/fxapp"
	"github.com/DavidReque/go-food-delivery/internal/pkg/fxapp/contracts"
)

type CatalogsWriteApplicationBuilder struct {
	contracts.ApplicationBuilder
}

func NewCatalogsWriteApplicationBuilder() *CatalogsWriteApplicationBuilder {
	builder := &CatalogsWriteApplicationBuilder{fxapp.NewApplicationBuilder()}

	return builder
}

/*func (a *CatalogsWriteApplicationBuilder) Build() *CatalogsWriteApplication {
	return NewCatalogsWriteApplication 
}*/
