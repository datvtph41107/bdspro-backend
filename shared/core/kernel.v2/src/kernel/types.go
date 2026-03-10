package kernel

import "time"

type ComponentID string

type State int

const (
	StateIdle State = iota
	StateStarting
	StateRunning
	StateStopping
	StateStopped
	StateFailed
)

func (s State) String() string {
	switch s {
	case StateIdle:
		return "idle"
	case StateStarting:
		return "starting"
	case StateRunning:
		return "running"
	case StateStopping:
		return "stopping"
	case StateStopped:
		return "stopped"
	case StateFailed:
		return "failed"
	default:
		return "unknown"
	}
}

type Metadata struct {
	ID              ComponentID
	Name            string
	DependsOn       []ComponentID
	Critical        bool
	StartupTimeout  time.Duration
	ShutdownTimeout time.Duration
	Priority        int
}

func (m Metadata) normalized() Metadata {
	out := m
	if out.Name == "" {
		out.Name = string(out.ID)
	}
	if out.StartupTimeout <= 0 {
		out.StartupTimeout = 10 * time.Second
	}
	if out.ShutdownTimeout <= 0 {
		out.ShutdownTimeout = 15 * time.Second
	}
	return out
}
