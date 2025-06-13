package cqrs

import (
	"fmt"

	"go.uber.org/fx"
)

// Cuando registramos múltiples manejadores (handlers) con el tipo de salida `mediatr.RequestHandler`,
// obtenemos una excepción que dice `type already provided` (tipo ya proporcionado).
// Para solucionar esto, debemos usar anotaciones con etiquetas (tags).

// AsHandler es una función auxiliar que anota un constructor de manejador (handler)
// para indicar que proporciona un manejador al grupo especificado.
//
// Parámetros:
// - handler: el constructor del manejador a anotar
// - handlerGroupName: nombre del grupo al que pertenecerá el manejador
//
// Retorna:
// - El constructor anotado listo para ser usado con fx
func AsHandler(handler interface{}, handlerGroupName string) interface{} {
	return fx.Annotate(
		handler,
		// Indica que el handler implementa la interfaz HandlerRegisterer
		fx.As(new(HandlerRegisterer)),
		// Agrega una etiqueta de grupo al resultado para poder agrupar handlers similares
		fx.ResultTags(fmt.Sprintf(
			`group:"%s"`,
			handlerGroupName,
		)),
	)
}
