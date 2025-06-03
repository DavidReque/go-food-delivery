package cqrs

import "github.com/mehdihadeli/go-mediatr"

// HandlerRegisterer es una interfaz que define el método para registrar un handler
type HandlerRegisterer interface {
	RegisterHandler() error
}

// RequestHandlerWithRegisterer para registrar un RequestHandler en el registro de mediatr
type RequestHandlerWithRegisterer[TRequest any, TResponse any] interface {
	HandlerRegisterer
	mediatr.RequestHandler[TRequest, TResponse]
}
