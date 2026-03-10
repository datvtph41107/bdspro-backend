package contract_runtime

import "context"

type ServiceStatus int

type ServiceDependency struct {
	Name         string
	Type         ServiceType
	Dependencies []ServiceType
	Critical     bool
	AutoStart    bool
	Timeout      int
}

// const (
// 	StatusIdle ServiceStatus = iota
// 	StatusRunning
// 	StatusStopping
// 	StatusStopped
// 	StatusFailed
// )

type Service interface {
	Get() ServiceDependency
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Status() ServiceStatus
}
