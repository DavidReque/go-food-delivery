package app

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/fxapp"
	"github.com/DavidReque/go-food-delivery/internal/pkg/fxapp/contracts"
)

type CatalogsReadApplicationBuilder struct {
	contracts.ApplicationBuilder
}

func NewCatalogsReadApplicationBuilder() *CatalogsReadApplicationBuilder {
	builder := &CatalogsReadApplicationBuilder{fxapp.NewApplicationBuilder()}

	return builder
}

func (a *CatalogsReadApplicationBuilder) Build() *CatalogsReadApplication {
	return NewCatalogsReadApplication(
		a.GetProvides(),
		a.GetInvokes(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
