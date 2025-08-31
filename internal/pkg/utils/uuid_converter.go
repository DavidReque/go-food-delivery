package utils

import (
	gofrsUUID "github.com/gofrs/uuid"
	googleUUID "github.com/google/uuid"
	satoriUUID "github.com/satori/go.uuid"
)

// UUIDConverter proporciona funciones para convertir entre diferentes tipos de UUID
// utilizados en el proyecto (google/uuid, gofrs/uuid, satori/go.uuid)

// ConvertGoogleUUIDToSatoriUUID convierte github.com/google/uuid.UUID a github.com/satori/go.uuid.UUID
func ConvertGoogleUUIDToSatoriUUID(id googleUUID.UUID) satoriUUID.UUID {
	u, err := satoriUUID.FromBytes(id[:])
	if err != nil {
		return satoriUUID.UUID{}
	}
	return u
}

// ConvertSatoriUUIDToGoogleUUID convierte github.com/satori/go.uuid.UUID a github.com/google/uuid.UUID
func ConvertSatoriUUIDToGoogleUUID(id satoriUUID.UUID) googleUUID.UUID {
	u, err := googleUUID.FromBytes(id[:])
	if err != nil {
		return googleUUID.Nil
	}
	return u
}

// ConvertGoogleUUIDToGofrsUUID convierte github.com/google/uuid.UUID a github.com/gofrs/uuid.UUID
func ConvertGoogleUUIDToGofrsUUID(id googleUUID.UUID) gofrsUUID.UUID {
	return gofrsUUID.Must(gofrsUUID.FromString(id.String()))
}

// ConvertGofrsUUIDToGoogleUUID convierte github.com/gofrs/uuid.UUID a github.com/google/uuid.UUID
func ConvertGofrsUUIDToGoogleUUID(id gofrsUUID.UUID) googleUUID.UUID {
	return googleUUID.Must(googleUUID.Parse(id.String()))
}

// ConvertGofrsUUIDToSatoriUUID convierte github.com/gofrs/uuid.UUID a github.com/satori/go.uuid.UUID
func ConvertGofrsUUIDToSatoriUUID(id gofrsUUID.UUID) satoriUUID.UUID {
	u, err := satoriUUID.FromBytes(id.Bytes())
	if err != nil {
		return satoriUUID.UUID{}
	}
	return u
}

// ConvertSatoriUUIDToGofrsUUID convierte github.com/satori/go.uuid.UUID a github.com/gofrs/uuid.UUID
func ConvertSatoriUUIDToGofrsUUID(id satoriUUID.UUID) gofrsUUID.UUID {
	u, err := gofrsUUID.FromBytes(id[:])
	if err != nil {
		return gofrsUUID.Nil
	}
	return u
}

// ConvertGoogleUUIDToSatoriUUIDSimple convierte github.com/google/uuid.UUID a github.com/satori/go.uuid.UUID
// usando conversión directa de tipo (más eficiente cuando es posible)
func ConvertGoogleUUIDToSatoriUUIDSimple(id googleUUID.UUID) satoriUUID.UUID {
	return satoriUUID.UUID(id)
}
