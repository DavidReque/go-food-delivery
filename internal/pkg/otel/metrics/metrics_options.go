package metrics

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/config"
	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	typeMapper "github.com/DavidReque/go-food-delivery/internal/pkg/reflection/typemapper"
	"github.com/iancoleman/strcase"
)

// OTLPProvider define la configuración para un exportador de métricas genérico
// que utiliza el protocolo OpenTelemetry (OTLP).
type OTLPProvider struct {
	Name         string            `mapstructure:"name"`         // Nombre del proveedor (ej. "jaeger", "signoz").
	Enabled      bool              `mapstructure:"enabled"`      // Indica si este proveedor está habilitado.
	OTLPEndpoint string            `mapstructure:"otlpEndpoint"` // La URL del colector OTLP.
	OTLPHeaders  map[string]string `mapstructure:"otlpHeaders"`  // Cabeceras adicionales, como las de autenticación.
}

// MetricsOptions contiene la configuración para el servicio de métricas
type MetricsOptions struct {
	Host                string `mapstructure:"host"`                // Host donde se expondrá el endpoint de métricas (ej. para Prometheus).
	Port                string `mapstructure:"port"`                // Puerto para el endpoint de métricas.
	ServiceName         string `mapstructure:"serviceName"`         // Nombre del servicio, usado para identificar las métricas en el backend.
	Version             string `mapstructure:"version"`             // Versión del servicio.
	MetricsRoutePath    string `mapstructure:"metricsRoutePath"`    // Ruta del endpoint de métricas (ej. "/metrics").
	EnableHostMetrics   bool   `mapstructure:"enableHostMetrics"`   // Habilita la recolección de métricas del host (CPU, memoria, etc.).
	UseStdout           bool   `mapstructure:"useStdout"`           // Si es true, las métricas se imprimirán en la consola (útil para debug).
	InstrumentationName string `mapstructure:"instrumentationName"` // Nombre de la librería de instrumentación.
	UseOTLP             bool   `mapstructure:"useOTLP"`             // Un interruptor maestro para habilitar la exportación a través de OTLP.

	// OTLPProviders permite configurar múltiples exportadores OTLP de forma genérica.
	OTLPProviders []OTLPProvider `mapstructure:"otlpProviders"`

	// Las siguientes son configuraciones específicas para proveedores OTLP populares.
	ElasticApmExporterOptions *OTLPProvider `mapstructure:"elasticApmExporterOptions"`
	UptraceExporterOptions    *OTLPProvider `mapstructure:"uptraceExporterOptions"`
	SignozExporterOptions     *OTLPProvider `mapstructure:"signozExporterOptions"`
}

// ProvideMetricsConfig es una función "provider" para el framework de inyección de dependencias (Uber FX).
// Su responsabilidad es leer la configuración de métricas desde el entorno de la aplicación
// y proveerla como una dependencia al resto de los componentes que la necesiten.
func ProvideMetricsConfig(
	environment environment.Environment, // Recibe el entorno para acceder a la configuración.
) (*MetricsOptions, error) {
	// Determina el nombre de la clave de configuración de forma dinámica.
	// 1. Obtiene el nombre del tipo `MetricsOptions`.
	// 2. Lo convierte a lowerCamelCase (ej. "metricsOptions").
	// Este será el nombre de la sección en tu archivo config.yaml.
	optionName := strcase.ToLowerCamel(
		typeMapper.GetGenericTypeNameByT[MetricsOptions](),
	)

	// Usa el helper BindConfigKey para leer la sección de configuración y
	// volcarla en una estructura `MetricsOptions`.
	return config.BindConfigKey[MetricsOptions](optionName)
}
