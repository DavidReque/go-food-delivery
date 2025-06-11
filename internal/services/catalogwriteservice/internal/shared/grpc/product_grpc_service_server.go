package grpc

import (
	"context"
	"fmt"

	customErrors "github.com/DavidReque/go-food-delivery/internal/pkg/http/httperrors/customerrors"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/otel/tracing/attribute"
	createProductCommandV1 "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1"
	createProductDtosV1 "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/dtos"
	"github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/shared/contracts"
	productsService "github.com/DavidReque/go-food-delivery/internal/services/catalogwriteservice/internal/shared/grpc/genproto"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/pkg/errors"

	attribute2 "go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var grpcMetricsAttr = api.WithAttributes(
	attribute2.Key("MetricsType").String("Http"),
)

type ProductGrpcServiceServer struct {
	catalogsMetrics *contracts.CatalogsMetrics
	logger          logger.Logger
}

// CreateProduct es un método del servidor gRPC que maneja la creación de nuevos productos.
// Recibe un contexto y una solicitud de creación de producto, y devuelve una respuesta con el ID del producto creado.
func (s *ProductGrpcServiceServer) CreateProduct(
	ctx context.Context,
	req *productsService.CreateProductReq,
) (*productsService.CreateProductRes, error) {
	// Obtiene el span de tracing del contexto para monitoreo y observabilidad
	span := trace.SpanFromContext(ctx)
	// Agrega los datos de la solicitud como atributos del span
	span.SetAttributes(attribute.Object("Request", req))
	// Incrementa el contador de métricas para las solicitudes gRPC de creación de productos
	s.catalogsMetrics.CreateProductGrpcRequests.Add(ctx, 1, grpcMetricsAttr)

	// Crea y valida un nuevo comando de creación de producto con los datos de la solicitud
	command, err := createProductCommandV1.NewCreateProductWithValidation(
		req.GetName(),
		req.GetDescription(),
		req.GetPrice(),
	)
	if err != nil {
		// Si la validación falla, crea un error personalizado con contexto adicional
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"[ProductGrpcServiceServer_CreateProduct.StructCtx] command validation failed",
		)
		// Registra el error de validación en los logs
		s.logger.Errorf(
			fmt.Sprintf(
				"[ProductGrpcServiceServer_CreateProduct.StructCtx] err: %v",
				validationErr,
			),
		)
		return nil, validationErr
	}

	// Envía el comando validado al mediador para su procesamiento
	// mediatr.Send es un patrón mediador que maneja la lógica de negocio
	result, err := mediatr.Send[*createProductCommandV1.CreateProduct, *createProductDtosV1.CreateProductResponseDto](
		ctx,
		command,
	)
	if err != nil {
		// Si ocurre un error durante el procesamiento, agrega contexto adicional
		err = errors.WithMessage(
			err,
			"[ProductGrpcServiceServer_CreateProduct.Send] error in sending CreateProduct",
		)
		// Registra el error con información adicional incluyendo el ID del producto
		s.logger.Errorw(
			fmt.Sprintf(
				"[ProductGrpcServiceServer_CreateProduct.Send] id: {%s}, err: %v",
				command.ProductID,
				err,
			),
			logger.Fields{"Id": command.ProductID},
		)
		return nil, err
	}

	// Si todo es exitoso, devuelve una respuesta con el ID del producto creado
	return &productsService.CreateProductRes{
		ProductId: result.ProductID.String(),
	}, nil
}
