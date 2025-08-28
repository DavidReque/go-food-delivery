package domainExceptions

import (
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"

	"emperror.dev/errors"
)

type invalidDeliveryAddressError struct {
	customErrors.BadRequestError
}
type InvalidDeliveryAddressError interface {
	customErrors.BadRequestError
}

func NewInvalidDeliveryAddressError(message string) error {
	originalErr := errors.New(message)

	bad := customErrors.NewBadRequestError(originalErr, message)
	customErr := customErrors.GetCustomError(bad).(customErrors.BadRequestError)
	br := &invalidDeliveryAddressError{
		BadRequestError: customErr,
	}

	return errors.WithStackIf(br)
}

func (i *invalidDeliveryAddressError) isInvalidAddress() bool {
	return true
}

func IsInvalidDeliveryAddressError(err error) bool {
	var ia *invalidDeliveryAddressError
	if errors.As(err, &ia) {
		return ia.isInvalidAddress()
	}

	return false
}
