package state

type ProcessState int

const (
	Idle ProcessState = iota
	Booting
	Running
	Stopping
	Stopped
)
