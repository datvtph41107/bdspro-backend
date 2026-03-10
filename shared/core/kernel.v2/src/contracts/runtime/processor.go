package contract_runtime

import "reflect"

type Processor interface {
	Type() reflect.Type
	Name() string
	Dependencies() []reflect.Type
	Priority() int
	Validate(comp Component) error
}

type ProcessorRegistry interface {
	Register(p Processor) error
	Get(t reflect.Type) (Processor, error)
	GetAll() []Processor
}
