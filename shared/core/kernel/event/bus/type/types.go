package bus_type

import "time"

// Marker interfaces & optional behavior interfaces
type Command interface{ IsCommand() }
type Event interface{ IsEvent() }
type Job interface{ IsJob() }

// Behavior interfaces
type ShouldQueue interface {
	QueueName() string
}
type HasRetry interface {
	Retries() int
}
type HasDelay interface {
	Delay() time.Duration
}
