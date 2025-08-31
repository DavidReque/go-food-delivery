package utils

import (
	"github.com/gofrs/uuid"
	"github.com/google/uuid"
	satoriUUID "github.com/satori/go.uuid"
)

// UUIDConverter proporciona funciones para convertir entre diferentes tipos de UUID
// utilizados en el proyecto (google/uuid, gofrs/uuid, satori/go.uuid)

// ConvertGoogleUUIDToSatoriUUID convierte github.com/google/uuid.UUID a github.com/satori/go.uuid.UUID
func ConvertGoogleUUIDToSatoriUUID(id uuid.UUID) satoriUUID.UUID {
	u, err := satoriUUID.FromBytes(id[:])
	if err != nil {
		return satoriUUID.UUID{}
	}
	return u
}

// ConvertSatoriUUIDToGoogleUUID convierte github.com/satori/go.uuid.UUID a github.com/google/uuid.UUID
func ConvertSatoriUUIDToGoogleUUID(id satoriUUID.UUID) uuid.UUID {
	u, err := uuid.FromBytes(id[:])
	if err != nil {
		return uuid.Nil
	}
	return u
}

// ConvertGoogleUUIDToGofrsUUID convierte github.com/google/uuid.UUID a github.com/gofrs/uuid.UUID
func ConvertGoogleUUIDToGofrsUUID(id uuid.UUID) uuid.UUID {
	return uuid.Must(uuid.FromString(id.String()))
}

// ConvertGofrsUUIDToGoogleUUID convierte github.com/gofrs/uuid.UUID a github.com/google/uuid.UUID
func ConvertGofrsUUIDToGoogleUUID(id uuid.UUID) uuid.UUID {
	return uuid.Must(uuid.Parse(id.String()))
}

// ConvertGofrsUUIDToSatoriUUID convierte github.com/gofrs/uuid.UUID a github.com/satori/go.uuid.UUID
func ConvertGofrsUUIDToSatoriUUID(id uuid.UUID) satoriUUID.UUID {
	u, err := satoriUUID.FromBytes(id.Bytes())
	if err != nil {
		return satoriUUID.UUID{}
	}
	return u
}

// ConvertSatoriUUIDToGofrsUUID convierte github.com/satori/go.uuid.UUID a github.com/gofrs/uuid.UUID
func ConvertSatoriUUIDToGofrsUUID(id satoriUUID.UUID) uuid.UUID {
	u, err := uuid.FromBytes(id[:])
	if err != nil {
		return uuid.Nil
	}
	return u
}

// ConvertGoogleUUIDToSatoriUUIDSimple convierte github.com/google/uuid.UUID a github.com/satori/go.uuid.UUID
// usando conversión directa de tipo (más eficiente cuando es posible)
func ConvertGoogleUUIDToSatoriUUIDSimple(id uuid.UUID) satoriUUID.UUID {
	return satoriUUID.UUID(id)
}
