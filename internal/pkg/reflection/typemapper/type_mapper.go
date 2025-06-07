package typemapper

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/iancoleman/strcase"
)

var (
	types    map[string][]reflect.Type
	packages map[string]map[string][]reflect.Type
)

// emptyInterface es una representación de una interfaz vacía
type emptyInterface struct {
	typ  unsafe.Pointer
	data unsafe.Pointer
}

// discoverTypes inicializa types y packages
func init() {
	types = make(map[string][]reflect.Type)
	packages = make(map[string]map[string][]reflect.Type)

	discoverTypes()
}

func discoverTypes() {
	typ := reflect.TypeOf(0)
	sections, offset := typelinks2()
	for i, offs := range offset {
		rodata := sections[i]
		for _, off := range offs {
			emptyInterface := (*emptyInterface)(unsafe.Pointer(&typ))
			emptyInterface.data = resolveTypeOff(rodata, off)
			if typ.Kind() == reflect.Ptr &&
				typ.Elem().Kind() == reflect.Struct {
				// just discover pointer types, but we also register this pointer type actual struct type to the registry
				loadedTypePtr := typ
				loadedType := typ.Elem()

				pkgTypes := packages[loadedType.PkgPath()]
				pkgTypesPtr := packages[loadedTypePtr.PkgPath()]

				if pkgTypes == nil {
					pkgTypes = map[string][]reflect.Type{}
					packages[loadedType.PkgPath()] = pkgTypes
				}
				if pkgTypesPtr == nil {
					pkgTypesPtr = map[string][]reflect.Type{}
					packages[loadedTypePtr.PkgPath()] = pkgTypesPtr
				}

				types[GetFullTypeNameByType(loadedType)] = append(
					types[GetFullTypeNameByType(loadedType)],
					loadedType,
				)
				types[GetFullTypeNameByType(loadedTypePtr)] = append(
					types[GetFullTypeNameByType(loadedTypePtr)],
					loadedTypePtr,
				)
				types[GetTypeNameByType(loadedType)] = append(
					types[GetTypeNameByType(loadedType)],
					loadedType,
				)
				types[GetTypeNameByType(loadedTypePtr)] = append(
					types[GetTypeNameByType(loadedTypePtr)],
					loadedTypePtr,
				)
			}
		}
	}
}

func RegisterType(typ reflect.Type) {
	types[GetFullTypeName(typ)] = append(types[GetFullTypeName(typ)], typ)
	types[GetTypeName(typ)] = append(types[GetTypeName(typ)], typ)
}

func RegisterTypeWithKey(key string, typ reflect.Type) {
	types[key] = append(types[key], typ)
}

func GetAllRegisteredTypes() map[string][]reflect.Type {
	return types
}

// TypeByName retorna el tipo por su nombre
func TypeByName(typeName string) reflect.Type {
	if typ, ok := types[typeName]; ok {
		return typ[0]
	}
	return nil
}

func TypesByName(typeName string) []reflect.Type {
	if types, ok := types[typeName]; ok {
		return types
	}
	return nil
}

func TypeByNameAndImplementedInterface[TInterface interface{}](
	typeName string,
) reflect.Type {
	implementedInterface := GetGenericTypeByT[TInterface]()
	if types, ok := types[typeName]; ok {
		for _, t := range types {
			if t.Implements(implementedInterface) {
				return t
			}
		}
	}

	return nil
}

func TypesImplementedInterfaceWithFilterTypes[TInterface interface{}](
	types []reflect.Type,
) []reflect.Type {
	implementedInterface := GetGenericTypeByT[TInterface]()

	var res []reflect.Type
	for _, t := range types {
		if t.Implements(implementedInterface) {
			res = append(res, t)
		}
	}

	return res
}

func TypesImplementedInterface[TInterface interface{}]() []reflect.Type {
	implementedInterface := GetGenericTypeByT[TInterface]()

	var res []reflect.Type
	for _, t := range types {
		for _, v := range t {
			if v.Implements(implementedInterface) {
				res = append(res, v)
			}
		}
	}

	return res
}

// GetFullTypeName retorna el nombre completo del tipo incluyendo el paquete
func GetFullTypeName(input interface{}) string {
	if input == nil {
		return ""
	}

	t := reflect.TypeOf(input)
	return t.String()
}

func GetGenericFullTypeNameByT[T any]() string {
	t := reflect.TypeOf((*T)(nil)).Elem()
	return t.String()
}

func GetFullTypeNameByType(typ reflect.Type) string {
	return typ.String()
}

// GetTypeName retorna el nombre del tipo sin el paquete
func GetTypeName(input interface{}) string {
	if input == nil {
		return ""
	}

	t := reflect.TypeOf(input)
	if t.Kind() != reflect.Ptr {
		return t.Name()
	}

	return fmt.Sprintf("*%s", t.Elem().Name())
}

func GetSnakeTypeName(input interface{}) string {
	if input == nil {
		return ""
	}

	t := reflect.TypeOf(input)
	if t.Kind() != reflect.Ptr {
		return t.Name()
	}

	return strcase.ToSnake(t.Elem().Name())
}

func GetKebabTypeName(input interface{}) string {
	if input == nil {
		return ""
	}

	t := reflect.TypeOf(input)
	if t.Kind() != reflect.Ptr {
		return t.Name()
	}

	return strcase.ToKebab(t.Elem().Name())
}

func GetGenericTypeNameByT[T any]() string {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() != reflect.Ptr {
		return t.Name()
	}

	return fmt.Sprintf("*%s", t.Elem().Name())
}

func GetGenericNonePointerTypeNameByT[T any]() string {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() != reflect.Ptr {
		return t.Name()
	}

	return t.Elem().Name()
}

// GetNonePointerTypeName retorna el nombre del tipo sin el paquete y sin el puntero
func GetNonePointerTypeName(input interface{}) string {
	if input == nil {
		return ""
	}

	t := reflect.TypeOf(input)
	if t.Kind() != reflect.Ptr {
		return t.Name()
	}

	return t.Elem().Name()
}

func GetTypeNameByType(typ reflect.Type) string {
	if typ == nil {
		return ""
	}

	if typ.Kind() != reflect.Ptr {
		return typ.Name()
	}

	return fmt.Sprintf("*%s", typ.Elem().Name())
}

// TypeByPackageName retorna el tipo por su paquete y nombre
func TypeByPackageName(pkgPath string, name string) reflect.Type {
	if pkgTypes, ok := packages[pkgPath]; ok {
		return pkgTypes[name][0]
	}
	return nil
}

func GetPackageName(value interface{}) string {
	inputType := reflect.TypeOf(value)
	if inputType.Kind() == reflect.Ptr {
		inputType = inputType.Elem()
	}

	packagePath := inputType.PkgPath()
	parts := strings.Split(packagePath, "/")

	return parts[len(parts)-1]
}

func TypesByPackageName(pkgPath string, name string) []reflect.Type {
	if pkgTypes, ok := packages[pkgPath]; ok {
		return pkgTypes[name]
	}
	return nil
}

func GetGenericTypeByT[T interface{}]() reflect.Type {
	res := reflect.TypeOf((*T)(nil)).Elem()
	return res
}

func GetBaseType(value interface{}) interface{} {
	if reflect.ValueOf(value).Kind() == reflect.Pointer {
		return reflect.ValueOf(value).Elem().Interface()
	}

	return value
}

func GetReflectType(value interface{}) reflect.Type {
	if reflect.TypeOf(value).Kind() == reflect.Pointer &&
		reflect.TypeOf(value).Elem().Kind() == reflect.Interface {
		return reflect.TypeOf(value).Elem()
	}

	res := reflect.TypeOf(value)
	return res
}

func GetBaseReflectType(value interface{}) reflect.Type {
	if reflect.ValueOf(value).Kind() == reflect.Pointer {
		return reflect.TypeOf(reflect.ValueOf(value).Elem().Interface())
	}

	return reflect.TypeOf(value)
}

func GenericInstanceByT[T any]() T {
	typ := GetGenericTypeByT[T]()
	return getInstanceFromType(typ).(T)
}

func InstanceByType(typ reflect.Type) interface{} {
	return getInstanceFromType(typ)
}

// InstanceByTypeName retorna una instancia vacía del tipo por su nombre
// Si el tipo es un puntero, retornará una instancia de puntero del tipo y
// si el tipo es una estructura, retornará una estructura vacía
func InstanceByTypeName(name string) interface{} {
	typ := TypeByName(name)
	return getInstanceFromType(typ)
}

func EmptyInstanceByTypeNameAndImplementedInterface[TInterface interface{}](
	name string,
) interface{} {
	typ := TypeByNameAndImplementedInterface[TInterface](name)
	return getInstanceFromType(typ)
}

func EmptyInstanceByTypeAndImplementedInterface[TInterface interface{}](
	typ reflect.Type,
) interface{} {
	typeName := GetTypeName(typ)
	return EmptyInstanceByTypeNameAndImplementedInterface[TInterface](typeName)
}

// InstancePointerByTypeName retorna una instancia de puntero vacía del tipo por su nombre
// Si el tipo es un puntero, retornará una instancia de puntero del tipo y
// si el tipo es una estructura, retornará un puntero a la estructura
func InstancePointerByTypeName(name string) interface{} {
	typ := TypeByName(name)
	if typ.Kind() == reflect.Ptr {
		res := reflect.New(typ.Elem()).Interface()
		return res
	}

	return reflect.New(typ).Interface()
}

// InstanceByPackageName retorna una instancia vacía del tipo por su nombre y paquete
// Si el tipo es un puntero, retornará una instancia de puntero del tipo y
// si el tipo es una estructura, retornará una estructura vacía
func InstanceByPackageName(pkgPath string, name string) interface{} {
	typ := TypeByPackageName(pkgPath, name)
	return getInstanceFromType(typ)
}

func getInstanceFromType(typ reflect.Type) interface{} {
	if typ.Kind() == reflect.Ptr {
		res := reflect.New(typ.Elem()).Interface()
		return res
	}

	return reflect.Zero(typ).Interface()
}

func GetGenericImplementInterfaceTypesT[T any]() map[string][]reflect.Type {
	result := make(map[string][]reflect.Type)

	// Get the interface type
	interfaceType := reflect.TypeOf((*T)(nil)).Elem()

	// Iterate over the types in the map
	for groupName, typeList := range types {
		var implementingTypes []reflect.Type

		// Check each type in the list
		for _, t := range typeList {
			// Check if the type implements the interface
			if t.Implements(interfaceType) {
				implementingTypes = append(implementingTypes, t)
			}
		}

		if len(implementingTypes) > 0 {
			result[groupName] = implementingTypes
		}
	}

	return result
}

func ImplementedInterfaceT[T any](obj interface{}) bool {
	// Get the interface type
	interfaceType := reflect.TypeOf((*T)(nil)).Elem()

	typ := GetReflectType(obj)

	implemented := typ.Implements(interfaceType)

	return implemented
}
