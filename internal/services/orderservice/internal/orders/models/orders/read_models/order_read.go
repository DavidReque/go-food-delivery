package read_models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// OrderReadModel es una estructura de datos optimizada para lectura (una proyección).
// Representa el estado de una orden tal como se almacena en una base de datos de solo lectura (probablemente MongoDB).
// Este modelo se construye y actualiza escuchando eventos del dominio de órdenes.
type OrderReadModel struct {
	// Id es el identificador único del documento en la colección de MongoDB.
	// Se genera un UUID v4 propio para tener control sobre el formato y evitar el ObjectID de MongoDB.
	// La etiqueta `bson:"_id"` lo mapea como la clave primaria en MongoDB.
	Id string `json:"id" bson:"_id,omitempty"` // https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/write-operations/insert/#the-_id-field

	// OrderId es el identificador de negocio de la orden, correspondiente al ID del agregado de dominio.
	OrderId string `json:"orderId" bson:"orderId,omitempty"`

	// ShopItems es una lista de los productos incluidos en la orden.
	ShopItems []*ShopItemReadModel `json:"shopItems,omitempty" bson:"shopItems,omitempty"`

	AccountEmail string `json:"accountEmail,omitempty" bson:"accountEmail,omitempty"`

	// DeliveryAddress es la dirección donde se debe entregar la orden.
	DeliveryAddress string `json:"deliveryAddress,omitempty" bson:"deliveryAddress,omitempty"`

	// CancelReason guarda el motivo por el cual una orden fue cancelada.
	CancelReason string `json:"cancelReason,omitempty" bson:"cancelReason,omitempty"`

	TotalPrice float64 `json:"totalPrice,omitempty" bson:"totalPrice,omitempty"`

	DeliveredTime time.Time `json:"deliveredTime,omitempty" bson:"deliveredTime,omitempty"`

	Paid bool `json:"paid,omitempty" bson:"paid,omitempty"`

	Submitted bool `json:"submitted,omitempty" bson:"submitted,omitempty"`

	Completed bool `json:"completed,omitempty" bson:"completed,omitempty"`

	Canceled bool `json:"canceled,omitempty" bson:"canceled,omitempty"`

	PaymentId string `json:"paymentId" bson:"paymentId,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`

	// UpdatedAt es la fecha y hora de la última actualización del registro de lectura.
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

// Se utiliza típicamente cuando se procesa un evento de creación de orden.
func NewOrderReadModel(
	orderId uuid.UUID,
	items []*ShopItemReadModel,
	accountEmail string,
	deliveryAddress string,
	deliveryTime time.Time,
) *OrderReadModel {
	return &OrderReadModel{
		// Genera un nuevo UUID para el campo `Id` del documento.
		Id: uuid.NewV4().
			String(),
		// Asigna los valores iniciales de la orden.
		OrderId:         orderId.String(),
		ShopItems:       items,
		AccountEmail:    accountEmail,
		DeliveryAddress: deliveryAddress,
		// Calcula el precio total de la orden para almacenarlo directamente.
		TotalPrice:    getShopItemsTotalPrice(items),
		DeliveredTime: deliveryTime,
		CreatedAt:     time.Now(),
	}
}

// getShopItemsTotalPrice es una función de ayuda para calcular el precio total
// sumando el precio de cada artículo multiplicado por su cantidad.
func getShopItemsTotalPrice(shopItems []*ShopItemReadModel) float64 {
	var totalPrice float64 = 0
	for _, item := range shopItems {
		totalPrice += item.Price * float64(item.Quantity)
	}

	return totalPrice
}
