package contract_runtime

import "context"

type ComponentStatus int

const (
	COMP_STATUS_IDLE ComponentStatus = iota
	COMP_STATUS_STARTING
	COMP_STATUS_RUNNING
	COMP_STATUS_STOPPING
	COMP_STATUS_STOPPED
	COMP_STATUS_FAILED
)

type Component interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Status() ComponentStatus
}

func (s ComponentStatus) String() string {
	return [...]string{
		"idle",
		"starting",
		"running",
		"stopping",
		"stopped",
		"failed",
	}[s]
}
