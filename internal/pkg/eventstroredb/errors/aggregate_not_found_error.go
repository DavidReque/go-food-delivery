package errors

import (
	"fmt"

	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"

	"emperror.dev/errors"
	satoriUUID "github.com/satori/go.uuid"
)

// https://klotzandrew.com/blog/error-handling-in-golang/
// https://banzaicloud.com/blog/error-handling-go/

type aggregateNotFoundError struct {
	customErrors.NotFoundError
}

type AggregateNotFoundError interface {
	customErrors.NotFoundError
	IsAggregateNotFoundError() bool
}

func NewAggregateNotFoundError(err error, id satoriUUID.UUID) error {
	notFound := customErrors.NewNotFoundErrorWrap(err, fmt.Sprintf("aggregate with id %s not found", id.String()))
	customErr := customErrors.GetCustomError(notFound)
	br := &aggregateNotFoundError{
		NotFoundError: customErr.(customErrors.NotFoundError),
	}

	return errors.WithStackIf(br)
}

func (err *aggregateNotFoundError) IsAggregateNotFoundError() bool {
	return true
}

func IsAggregateNotFoundError(err error) bool {
	var an AggregateNotFoundError
	if errors.As(err, &an) {
		return an.IsAggregateNotFoundError()
	}

	return false
}
