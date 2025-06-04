package mapper

import "github.com/jinzhu/copier"

// Map convierte un objeto de un tipo a otro usando reflection
// TDestination: tipo al que se quiere convertir
// source: objeto que se quiere convertir
func Map[TDestination any](source interface{}) (TDestination, error) {
	var dest TDestination
	err := copier.Copy(&dest, source)
	if err != nil {
		return dest, err
	}
	return dest, nil
}

// MapSlice convierte un slice de objetos de un tipo a otro
// TDestination: tipo al que se quiere convertir cada elemento
// source: slice de objetos que se quiere convertir
func MapSlice[TDestination any](source interface{}) ([]TDestination, error) {
	var dest []TDestination
	err := copier.Copy(&dest, source)
	if err != nil {
		return nil, err
	}
	return dest, nil
}
