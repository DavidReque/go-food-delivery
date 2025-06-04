package datamodels

import (
	"encoding/json"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// ProductDataModel representa el modelo de datos para la tabla de productos
type ProductDataModel struct {
	Id          uuid.UUID `gorm:"primaryKey"`
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time `gorm:"default:current_timestamp"`
	UpdatedAt   time.Time
	// for soft delete - https://gorm.io/docs/delete.html#Soft-Delete
	gorm.DeletedAt
}

// TableName devuelve el nombre de la tabla en la base de datos
func (p *ProductDataModel) TableName() string {
	return "products"
}

// String devuelve una representación en formato JSON del modelo de datos
func (p *ProductDataModel) String() string {
	j, _ := json.Marshal(p)

	return string(j)
}
