package bus

type Meta struct {
	Async     bool
	QueueName string
	Retries   int
	DelayMs   int64
	TypeName  string
}
