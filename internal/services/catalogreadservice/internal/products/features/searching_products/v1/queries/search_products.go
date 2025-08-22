package queries

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation"
)

type SearchProducts struct {
	SearchText string
	*utils.ListQuery
}

func (s *SearchProducts) Validate() error {
	return validation.ValidateStruct(s, validation.Field(&s.SearchText, validation.Required))
}
