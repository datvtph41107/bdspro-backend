package contract_runtime

import "context"

type runnerStatus int

const (
	StatusRunnerIdle runnerStatus = iota
	StatusRunnerStarting
	StatusRunnerRunning
	StatusRunnerStopping
	StatusRunnerStopped
)

type Runtime interface {
	RegisterComponent(srv Service) error
	HandlerComponent(srv Service)
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
	GetService(t ServiceType) (Service, error)
	GetListServices() []ServiceDependency
}
