package logger

import (
	contract_runtime "core/kernel.v2/src/contracts/runtime"
	"fmt"
	"reflect"
)

type Processor struct{}

func (p *Processor) Type() reflect.Type {
	return reflect.TypeOf((*Component)(nil))
}

func (p *Processor) Name() string {
	return "logger"
}

func (p *Processor) Dependencies() []reflect.Type {
	return nil
}

func (p *Processor) Priority() int {
	return 0
}

func (p *Processor) Validate(comp contract_runtime.Component) error {
	_, ok := comp.(*Component)
	if !ok {
		return fmt.Errorf("expected logger component")
	}
	return nil
}
