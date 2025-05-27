package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

// ValidationError representa errores de validación específicos
type ValidationError interface {
	BadRequestError
	isValidationError()
}

// validationError implementa ValidationError
type validationError struct {
	CustomError
}

func (v *validationError) isValidationError() {
	// Método marcador para identificar errores de validación
}

func (v *validationError) isBadRequestError() {
	// Método marcador para identificar errores de bad request
}

// NewValidationError crea un nuevo error de validación
func NewValidationError(message string) ValidationError {
	validationErrMessage := errors.NewPlain("validation error")

	stackErr := errors.WrapIf(validationErrMessage, message)

	validationError := &validationError{
		CustomError: NewCustomError(stackErr, http.StatusBadRequest, message),
	}

	return validationError
}

// NewValidationErrorWrap crea un error de validación wrapeando un error existente
func NewValidationErrorWrap(err error, message string) ValidationError {
	if err == nil {
		return NewValidationError(message)
	}

	validationErrMessage := errors.WithMessage(err, "validation error")

	stackErr := errors.WrapIf(validationErrMessage, message)

	validationError := &validationError{
		CustomError: NewCustomError(stackErr, http.StatusBadRequest, message),
	}

	return validationError
}

// IsValidationError verifica si un error es un ValidationError
func IsValidationError(err error) bool {
	var validationError ValidationError

	if _, ok := err.(ValidationError); ok {
		return true
	}

	if errors.As(err, &validationError) {
		return true
	}

	return false
}
