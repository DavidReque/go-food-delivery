package domainExceptions

import (
	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"

	"emperror.dev/errors"
)

type orderShopItemsRequiredError struct {
	customErrors.BadRequestError
}

type OrderShopItemsRequiredError interface {
	customErrors.BadRequestError
}

func NewOrderShopItemsRequiredError(message string) error {
	originalErr := errors.New(message)
	bad := customErrors.NewBadRequestError(originalErr, message)
	customErr := customErrors.GetCustomError(bad).(customErrors.BadRequestError)
	br := &orderShopItemsRequiredError{
		BadRequestError: customErr,
	}

	return errors.WithStackIf(br)
}

func (i *orderShopItemsRequiredError) isOrderShopItemsRequiredError() bool {
	return true
}

func IsOrderShopItemsRequiredError(err error) bool {
	var os *orderShopItemsRequiredError
	if errors.As(err, &os) {
		return os.isOrderShopItemsRequiredError()
	}

	return false
}
