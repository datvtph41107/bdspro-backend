package config

type ObservabilityConfig struct {
	// Logging
	LogLevel       string `mapstructure:"log_level"      validate:"oneof=debug info warn error" reload:"hot"`
	LogFormat      string `mapstructure:"log_format"     validate:"oneof=json text"`
	LogRequestBody bool   `mapstructure:"log_request_body"` // false trên prod

	// Metrics
	MetricsEnabled bool   `mapstructure:"metrics_enabled"`
	MetricsPort    int    `mapstructure:"metrics_port"    validate:"min=1024,max=65535"`
	MetricsPath    string `mapstructure:"metrics_path"`

	// Tracing
	TracingEnabled bool    `mapstructure:"tracing_enabled"`
	OTELEndpoint   string  `mapstructure:"otel_endpoint"`
	SamplingRate   float64 `mapstructure:"sampling_rate"  validate:"min=0,max=1" reload:"hot"`
	ServiceName    string  `mapstructure:"service_name"   validate:"required_if=TracingEnabled true"`
}
