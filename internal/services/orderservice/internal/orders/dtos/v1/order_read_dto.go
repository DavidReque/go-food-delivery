package dtosV1

import "time"

// OrderReadDto DTO for reading orders
// @Description DTO for reading orders from the system
type OrderReadDto struct {
	// @Description Unique ID of the order
	Id string `json:"id"`

	// @Description ID of the order for external reference
	OrderId string `json:"orderId"`

	// @Description List of items from the shop in the order
	ShopItems []*ShopItemReadDto `json:"shopItems"`

	// @Description Email of the user's account
	AccountEmail string `json:"accountEmail"`

	// @Description Delivery address
	DeliveryAddress string `json:"deliveryAddress"`

	// @Description Reason for cancellation (if applicable)
	CancelReason string `json:"cancelReason"`

	// @Description Total price of the order
	TotalPrice float64 `json:"totalPrice"`

	// @Description Delivery time
	// @Format date-time
	DeliveredTime time.Time `json:"deliveredTime"`

	// @Description Indicates if the order has been paid
	Paid bool `json:"paid"`

	// @Description Indicates if the order has been sent
	Submitted bool `json:"submitted"`

	// @Description Indicates if the order has been completed
	Completed bool `json:"completed"`

	// @Description Indicates if the order has been cancelled
	Canceled bool `json:"canceled"`

	// @Description ID of the associated payment
	PaymentId string `json:"paymentId"`

	// @Description Creation date
	// @Format date-time
	CreatedAt time.Time `json:"createdAt"`

	// @Description Last update date
	// @Format date-time
	UpdatedAt time.Time `json:"updatedAt"`
}
