package models

import "time"

type Product struct {
	// we generate id ourselves because auto generate mongo string id column with type _id is not an uuid
	// Id is the unique identifier in mongo db
	Id          string    `json:"id"                    bson:"_id,omitempty"` // https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/write-operations/insert/#the-_id-field
	ProductId   string    `json:"productId"             bson:"productId"`
	Name        string    `json:"name,omitempty"        bson:"name,omitempty"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"`
	Price       float64   `json:"price,omitempty"       bson:"price,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"   bson:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"   bson:"updatedAt,omitempty"`
}

type ProductsList struct {
	// TotalCount is the total number of products
	TotalCount int64 `json:"totalCount" bson:"totalCount"`
	// TotalPages is the total number of pages
	TotalPages int64 `json:"totalPages" bson:"totalPages"`
	// Page is the current page
	Page int64 `json:"page"       bson:"page"`
	// Size is the number of products per page
	Size int64 `json:"size"       bson:"size"`
	// Products is the list of products
	Products []*Product `json:"products"   bson:"products"`
}
