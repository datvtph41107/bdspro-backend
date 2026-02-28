package providers

import (
	"core/kernel.v2/src/contracts/tracer"
	"core/kernel.v2/src/pkg/oteltracer"
)

type TracerConfig struct {
	ServiceName    string  `yaml:"service_name"`
	ServiceVersion string  `yaml:"service_version"`
	Environment    string  `yaml:"environment"`
	Endpoint       string  `yaml:"endpoint"`
	SamplingRate   float64 `yaml:"sampling_rate"`
}

func ProvideTracer(cfg TracerConfig) (tracer.Tracer, error) {
	return oteltracer.New(tracer.Config{
		ServiceName:    cfg.ServiceName,
		ServiceVersion: cfg.ServiceVersion,
		Environment:    cfg.Environment,
		Endpoint:       cfg.Endpoint,
		SamplingRate:   cfg.SamplingRate,
	})
}
