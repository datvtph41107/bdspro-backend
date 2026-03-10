package http

import (
	contract_runtime "core/kernel.v2/src/contracts/runtime"
	"core/kernel.v2/src/modules/logger"
	"fmt"
	"reflect"
)

type Processor struct{}

func (p *Processor) Type() reflect.Type {
	return reflect.TypeOf((*HTTPServer)(nil)).Elem()
}

func (p *Processor) Name() string {
	return "http-server"
}

func (p *Processor) Dependencies() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*logger.Component)(nil)),
		reflect.TypeOf((*tracer.Component)(nil)),
		reflect.TypeOf((*metrics.Component)(nil)),
	}
}

func (p *Processor) Priority() int {
	return 20
}

func (p *Processor) Validate(comp contract_runtime.Component) error {
	if _, ok := comp.(HTTPServer); !ok {
		return fmt.Errorf("component does not implement HTTPServer interface")
	}
	return nil
}
