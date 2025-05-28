package contracts

import (
	"go.opentelemetry.io/otel/metric"
)

type CatalogsMetrics struct {
	CreateProductGrpcRequests     metric.Float64Counter // Cuántas veces se llama crear producto
	UpdateProductGrpcRequests     metric.Float64Counter
	DeleteProductGrpcRequests     metric.Float64Counter
	GetProductByIdGrpcRequests    metric.Float64Counter
	SearchProductGrpcRequests     metric.Float64Counter
	SuccessRabbitMQMessages       metric.Float64Counter // Mensajes procesados exitosamente
	ErrorRabbitMQMessages         metric.Float64Counter // Mensajes procesados con error
	CreateProductRabbitMQMessages metric.Float64Counter
	UpdateProductRabbitMQMessages metric.Float64Counter
	DeleteProductRabbitMQMessages metric.Float64Counter
}
