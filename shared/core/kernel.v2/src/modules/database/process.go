package database

import (
	"fmt"
	"reflect"

	"core/kernel.v2/src/contracts/database"
	contract_runtime "core/kernel.v2/src/contracts/runtime"
	"core/kernel.v2/src/modules/logger"
)

type Processor struct{}

func (p *Processor) Type() reflect.Type {
	return reflect.TypeOf((*database.Database)(nil)).Elem()
}

func (p *Processor) Name() string {
	return "database"
}

func (p *Processor) Dependencies() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*logger.Component)(nil))}
}

func (p *Processor) Priority() int {
	return 10
}

func (p *Processor) Validate(comp contract_runtime.Component) error {
	if _, ok := comp.(database.Database); !ok {
		return fmt.Errorf("component does not implement Database interface")
	}
	return nil
}
