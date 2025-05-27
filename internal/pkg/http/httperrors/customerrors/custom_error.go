package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

// CustomError es la interfaz base para todos los errores personalizados
type CustomError interface {
	error
	GetStatusCode() int
	GetMessage() string
	GetError() error
	isCustomError()
}

// BadRequestError representa errores de petición incorrecta (400)
type BadRequestError interface {
	CustomError
	isBadRequestError()
}

// customError implementa CustomError
type customError struct {
	err        error
	statusCode int
	message    string
}

func (c *customError) Error() string {
	return c.message
}

func (c *customError) GetStatusCode() int {
	return c.statusCode
}

func (c *customError) GetMessage() string {
	return c.message
}

func (c *customError) GetError() error {
	return c.err
}

func (c *customError) isCustomError() {
	// Método marcador para identificar errores personalizados
}

// NewCustomError crea un nuevo error personalizado
func NewCustomError(err error, statusCode int, message string) CustomError {
	return &customError{
		err:        err,
		statusCode: statusCode,
		message:    message,
	}
}

// badRequestError implementa BadRequestError
type badRequestError struct {
	CustomError
}

func (b *badRequestError) isBadRequestError() {
	// Método marcador para identificar errores de bad request
}

// NewBadRequestError crea un nuevo error de bad request (400)
func NewBadRequestError(err error, message string) BadRequestError {
	return &badRequestError{
		CustomError: NewCustomError(err, http.StatusBadRequest, message),
	}
}

// IsCustomError verifica si un error es un CustomError
func IsCustomError(err error) bool {
	var customErr CustomError
	return errors.As(err, &customErr)
}

// IsBadRequestError verifica si un error es un BadRequestError
func IsBadRequestError(err error) bool {
	var badReqErr BadRequestError
	return errors.As(err, &badReqErr)
}
