package route

import (
	"fmt"

	"go.uber.org/fx"
)

// Cuando registramos múltiples manejadores HTTP con el tipo de salida `echo.HandlerFunc`,
// obtenemos una excepción que dice `type already provided` (tipo ya proporcionado).
// Para solucionar esto, debemos usar anotaciones con etiquetas (tags).

// AsRoute es una función auxiliar que anota un constructor de ruta HTTP
// para indicar que proporciona un endpoint al grupo de rutas especificado.
//
// Parámetros:
// - handler: el constructor del manejador HTTP a anotar
// - routeGroupName: nombre del grupo al que pertenecerá la ruta
//
// Retorna:
// - El constructor anotado listo para ser usado con fx
func AsRoute(handler interface{}, routeGroupName string) interface{} {
	return fx.Annotate(
		handler,
		// Indica que el handler implementa la interfaz Endpoint
		fx.As(new(Endpoint)),
		// Agrega una etiqueta de grupo al resultado para poder agrupar rutas similares
		fx.ResultTags(fmt.Sprintf(`group:"%s"`, routeGroupName)),
	)
}