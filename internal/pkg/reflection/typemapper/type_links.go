package typemapper

import (
	"unsafe"
)

// Basado en la implementación de Go runtime
// https://github.com/golang/go/blob/master/src/runtime/type.go

//go:linkname typelinks2 reflect.typelinks2
func typelinks2() (sections []unsafe.Pointer, offset [][]int32)

//go:linkname resolveTypeOff reflect.resolveTypeOff
func resolveTypeOff(rtype unsafe.Pointer, off int32) unsafe.Pointer
