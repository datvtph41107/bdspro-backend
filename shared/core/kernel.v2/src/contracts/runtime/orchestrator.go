package contract_runtime

import (
	"context"
	"reflect"
)

type OrchestratorStatus int

const (
	RUNNER_STATUS_IDLE OrchestratorStatus = iota
	RUNNER_STATUS_STARTING
	RUNNER_STATUS_RUNNING
	RUNNER_STATUS_STOPPING
	RUNNER_STATUS_STOPPED
)

type Orchestrator interface {
	RegisterComponent(comp Component) error
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
	GetComponent(typ reflect.Type) (Component, error)
	ListComponents() []reflect.Type
}

func (s OrchestratorStatus) String() string {
	return [...]string{
		"idle",
		"starting",
		"running",
		"stopping",
		"stopped",
	}[s]
}
