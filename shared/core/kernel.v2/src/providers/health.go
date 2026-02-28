package providers

import (
	"core/kernel.v2/src/contracts/health"
	"core/kernel.v2/src/pkg/simplehealth"
)

func ProvideHealthRegistry() health.Registry {
	return simplehealth.NewRegistry()
}
