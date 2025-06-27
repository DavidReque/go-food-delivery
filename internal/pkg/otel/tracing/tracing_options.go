package tracing

import (
	"github.com/DavidReque/go-food-delivery/internal/pkg/config"
	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	"github.com/DavidReque/go-food-delivery/internal/pkg/reflection/typemapper"
	"github.com/iancoleman/strcase"
)

// OTLPProvider es una estructura que contiene la configuración para un proveedor OTLP.
// Se utiliza para configurar los proveedores de OpenTelemetry.
type OTLPProvider struct {
	Name         string            `mapstructure:"name"`
	Enabled      bool              `mapstructure:"enabled"`
	OTLPEndpoint string            `mapstructure:"otlpEndpoint"`
	OTLPHeaders  map[string]string `mapstructure:"otlpHeaders"`
}

// TracingOptions es una estructura que contiene la configuración para el servicio de trazado.
type TracingOptions struct {
	Enabled                   bool                   `mapstructure:"enabled"`
	ServiceName               string                 `mapstructure:"serviceName"`
	Version                   string                 `mapstructure:"version"`
	InstrumentationName       string                 `mapstructure:"instrumentationName"`
	Id                        int64                  `mapstructure:"id"`
	AlwaysOnSampler           bool                   `mapstructure:"alwaysOnSampler"`
	ZipkinExporterOptions     *ZipkinExporterOptions `mapstructure:"zipkinExporterOptions"`
	JaegerExporterOptions     *OTLPProvider          `mapstructure:"jaegerExporterOptions"`
	ElasticApmExporterOptions *OTLPProvider          `mapstructure:"elasticApmExporterOptions"`
	UptraceExporterOptions    *OTLPProvider          `mapstructure:"uptraceExporterOptions"`
	SignozExporterOptions     *OTLPProvider          `mapstructure:"signozExporterOptions"`
	TempoExporterOptions      *OTLPProvider          `mapstructure:"tempoExporterOptions"`
	UseStdout                 bool                   `mapstructure:"useStdout"`
	UseOTLP                   bool                   `mapstructure:"useOTLP"`
	OTLPProviders             []OTLPProvider         `mapstructure:"otlpProviders"`
}

// ZipkinExporterOptions es una estructura que contiene la configuración para el exportador Zipkin.
type ZipkinExporterOptions struct {
	Url string `mapstructure:"url"`
}

func ProvideTracingConfig(
	environment environment.Environment,
) (*TracingOptions, error) {
	optionName := strcase.ToLowerCamel(
		typemapper.GetTypeName(TracingOptions{}),
	)

	return config.BindConfigKey[TracingOptions](optionName, config.WithEnvironment(environment))
}
