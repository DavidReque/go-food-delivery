package logger

import (
	"time"
)

// Fields representa un mapa de campos estructurados para logging
// Permite agregar contexto adicional a los logs de forma organizada
// Ejemplo: Fields{"userId": 123, "action": "create_product"}
type Fields map[string]interface{}

type Logger interface {
	Configure(cfg func(internalLog interface{}))

	// Métodos de Debug - Para información detallada durante desarrollo
	Debug(args ...interface{})                   // Debug con argumentos variables
	Debugf(template string, args ...interface{}) // Debug con formato printf
	Debugw(msg string, fields Fields)            // Debug con campos estructurados

	// LogType() models.LogType - Comentado, posiblemente para obtener tipo de log

	// Métodos de Info - Para información general en producción
	Info(args ...interface{})                   // Info con argumentos variables
	Infof(template string, args ...interface{}) // Info con formato printf
	Infow(msg string, fields Fields)            // Info con campos estructurados

	// Métodos de Warning - Para advertencias y situaciones que requieren atención
	Warn(args ...interface{})                   // Warning con argumentos variables
	Warnf(template string, args ...interface{}) // Warning con formato printf
	WarnMsg(msg string, err error)              // Warning con mensaje y error específico

	// Métodos de Error - Para errores que requieren investigación
	Error(args ...interface{})                   // Error con argumentos variables
	Errorw(msg string, fields Fields)            // Error con campos estructurados
	Errorf(template string, args ...interface{}) // Error con formato printf
	Err(msg string, err error)                   // Error con mensaje descriptivo y error

	// Métodos de Fatal - Para errores críticos que terminan la aplicación
	Fatal(args ...interface{})                   // Fatal con argumentos variables
	Fatalf(template string, args ...interface{}) // Fatal con formato printf

	// Método de propósito general - Compatible con fmt.Printf
	Printf(template string, args ...interface{})
	WithName(name string)

	// GrpcMiddlewareAccessLogger registra información de acceso para middleware gRPC
	// method: nombre del método gRPC llamado (ej: "CreateProduct")
	// time: duración de la ejecución del método
	// metaData: headers y metadata del request gRPC
	// err: error ocurrido durante la ejecución (nil si fue exitoso)
	GrpcMiddlewareAccessLogger(
		method string,
		time time.Duration,
		metaData map[string][]string,
		err error,
	)

	// GrpcClientInterceptorLogger registra información para clientes gRPC
	// method: método gRPC llamado
	// req: objeto request enviado
	// reply: objeto response recibido
	// time: latencia de la llamada
	// metaData: metadata de la conexión
	// err: error de la comunicación (nil si fue exitoso)
	GrpcClientInterceptorLogger(
		method string,
		req interface{},
		reply interface{},
		time time.Duration,
		metaData map[string][]string,
		err error,
	)
}
