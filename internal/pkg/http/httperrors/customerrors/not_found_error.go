package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

// NotFoundError representa errores de recurso no encontrado (404)
type NotFoundError interface {
	CustomError
	isNotFoundError()
}

// notFoundError implementa NotFoundError
type notFoundError struct {
	CustomError
}

func (n *notFoundError) isNotFoundError() {
	// Método marcador para identificar errores de not found
}

// NewNotFoundError crea un nuevo error de recurso no encontrado
func NewNotFoundError(message string) NotFoundError {
	notFoundErrMessage := errors.NewPlain("not found error")
	stackErr := errors.WrapIf(notFoundErrMessage, message)

	notFoundError := &notFoundError{
		CustomError: NewCustomError(stackErr, http.StatusNotFound, message),
	}

	return notFoundError
}

// NewNotFoundErrorWrap crea un error de not found wrapeando un error existente
func NewNotFoundErrorWrap(err error, message string) NotFoundError {
	if err == nil {
		return NewNotFoundError(message)
	}

	notFoundErrMessage := errors.WithMessage(err, "not found error")
	stackErr := errors.WrapIf(notFoundErrMessage, message)

	notFoundError := &notFoundError{
		CustomError: NewCustomError(stackErr, http.StatusNotFound, message),
	}

	return notFoundError
}

// IsNotFoundError verifica si un error es un NotFoundError
func IsNotFoundError(err error) bool {
	var notFoundError NotFoundError

	if _, ok := err.(NotFoundError); ok {
		return true
	}

	if errors.As(err, &notFoundError) {
		return true
	}

	return false
}
