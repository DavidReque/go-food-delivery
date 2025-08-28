package domainExceptions

import (
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"

	"emperror.dev/errors"
)

type invalidEmailAddressError struct {
	customErrors.BadRequestError
}

type InvalidEmailAddressError interface {
	customErrors.BadRequestError
}

func NewInvalidEmailAddressError(message string) error {
	originalErr := errors.New(message)

	bad := customErrors.NewBadRequestError(originalErr, message)
	customErr := customErrors.GetCustomError(bad).(customErrors.BadRequestError)
	br := &invalidEmailAddressError{
		BadRequestError: customErr,
	}

	return errors.WithStackIf(br)
}

func (i *invalidEmailAddressError) isInvalidEmailAddressError() bool {
	return true
}

func IsInvalidEmailAddressError(err error) bool {
	var ie *invalidEmailAddressError

	if errors.As(err, &ie) {
		return ie.isInvalidEmailAddressError()
	}

	return false
}
