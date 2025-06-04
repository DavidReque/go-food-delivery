package cqrs

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// TypeInfo proporciona información de tipo para mensajes
type TypeInfo interface {
	GetTypeName() string
	GetFullTypeName() string
}

// Request representa información básica de una petición
type Request interface {
	GetRequestId() string
	GetTimestamp() time.Time
}

// typeInfo implementa TypeInfo
type typeInfo struct {
	typeName     string
	fullTypeName string
}

func (t *typeInfo) GetTypeName() string {
	return t.typeName
}

func (t *typeInfo) GetFullTypeName() string {
	return t.fullTypeName
}

// request implementa Request
type request struct {
	requestId string
	timestamp time.Time
}

func (r *request) GetRequestId() string {
	return r.requestId
}

func (r *request) GetTimestamp() time.Time {
	return r.timestamp
}

// command implementa Command
type command struct {
	TypeInfo
	Request
}

type Command interface {
	isCommand()
	Request
	TypeInfo
}

func NewCommandByT[T any]() Command {
	c := &command{
		TypeInfo: newTypeInfo[T](),
		Request:  newRequest(),
	}
	return c
}

func (c *command) isCommand() {
	// Método marcador para identificar comandos
}

func IsCommand(obj interface{}) bool {
	if _, ok := obj.(Command); ok {
		return true
	}
	return false
}

// Funciones auxiliares para crear las implementaciones
func newTypeInfo[T any]() TypeInfo {
	var t T
	typeName := getTypeName(t)
	return &typeInfo{
		typeName:     typeName,
		fullTypeName: getFullTypeName(t),
	}
}

func newRequest() Request {
	return &request{
		requestId: uuid.NewV4().String(),
		timestamp: time.Now(),
	}
}

// Funciones helper para obtener nombres de tipos
func getTypeName(v interface{}) string {
	// Implementación simple - se puede mejorar con reflection
	return "Command" // Por ahora placeholder
}

func getFullTypeName(v interface{}) string {
	// Implementación simple - se puede mejorar con reflection
	return "cqrs.Command" // Por ahora placeholder
}
