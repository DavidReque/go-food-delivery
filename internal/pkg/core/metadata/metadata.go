package metadata

// Metadata es un tipo map[string]interface{} que representa los metadatos de un mensaje
// Se utiliza para transportar información adicional entre componentes de la aplicación
type Metadata map[string]interface{}

// ExistsKey verifica si una clave específica existe en los metadatos
func (m Metadata) ExistsKey(key string) bool {
	_, exists := m[key]
	return exists
}

// Get obtiene el valor de una clave específica en los metadatos
func (m Metadata) Get(key string) interface{} {
	val, exists := m[key]
	if !exists {
		return nil
	}

	return val
}

// Set establece un valor para una clave específica en los metadatos
func (m Metadata) Set(key string, value interface{}) {
	m[key] = value
}

// Keys devuelve una lista de todas las claves en los metadatos
func (m Metadata) Keys() []string {
	i := 0
	r := make([]string, len(m))

	for k := range m {
		r[i] = k
		i++
	}

	return r
}

// MapToMetadata convierte un map[string]interface{} a un tipo Metadata
func MapToMetadata(data map[string]interface{}) Metadata {
	m := Metadata(data)
	return m
}

// MetadataToMap convierte un tipo Metadata a un map[string]interface{}
func MetadataToMap(meta Metadata) map[string]interface{} {
	return meta
}

// FromMetadata devuelve un tipo Metadata a partir de un tipo Metadata
func FromMetadata(m Metadata) Metadata {
	if m == nil {
		return Metadata{}
	}
	return m
}
