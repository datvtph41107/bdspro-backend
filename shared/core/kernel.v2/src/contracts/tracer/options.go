package tracer

type Config struct {
	ServiceName    string  `yaml:"service_name"`
	ServiceVersion string  `yaml:"service_version"`
	Environment    string  `yaml:"environment"`
	Endpoint       string  `yaml:"endpoint"` // OTLP gRPC endpoint
	Insecure       bool    `yaml:"insecure" default:"true"`
	SamplingRate   float64 `yaml:"sampling_rate" default:"1.0"`
	BatchTimeout   int     `yaml:"batch_timeout" default:"5"` // seconds
	ExportTimeout  int     `yaml:"export_timeout" default:"30"`
	MaxQueueSize   int     `yaml:"max_queue_size" default:"2048"`
}
