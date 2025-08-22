package dtos

import uuid "github.com/satori/go.uuid"

// GetProductByIdRequestDto is the request dto for the get product by id endpoint
type GetProductByIdRequestDto struct {
	Id uuid.UUID `param:"id" json:"-"`
}
