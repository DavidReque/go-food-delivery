package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

// ConflictError representa errores de conflicto (409)
type ConflictError interface {
	CustomError
	isConflictError()
}

// conflictError implementa ConflictError
type conflictError struct {
	CustomError
}

func (c *conflictError) isConflictError() {
	// Método marcador para identificar errores de conflicto
}

// NewConflictError crea un nuevo error de conflicto
func NewConflictError(message string) ConflictError {
	conflictErrMessage := errors.NewPlain("conflict error")
	stackErr := errors.WrapIf(conflictErrMessage, message)

	conflictError := &conflictError{
		CustomError: NewCustomError(stackErr, http.StatusConflict, message),
	}

	return conflictError
}

// NewConflictErrorWrap crea un error de conflicto wrapeando un error existente
func NewConflictErrorWrap(err error, message string) ConflictError {
	if err == nil {
		return NewConflictError(message)
	}

	conflictErrMessage := errors.WithMessage(err, "conflict error")
	stackErr := errors.WrapIf(conflictErrMessage, message)

	conflictError := &conflictError{
		CustomError: NewCustomError(stackErr, http.StatusConflict, message),
	}

	return conflictError
}

// IsConflictError verifica si un error es un ConflictError
func IsConflictError(err error) bool {
	var conflictError ConflictError

	if _, ok := err.(ConflictError); ok {
		return true
	}

	if errors.As(err, &conflictError) {
		return true
	}

	return false
}
