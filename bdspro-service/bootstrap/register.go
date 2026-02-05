package bootstrap

import (
	"bdspro/internal/job"
	"core/kernel/container"
	"core/kernel/db"
	"core/kernel/queue"
	_redis "core/kernel/queue/redis"
)

func Register(c *container.Container) {
	// DB manager (singleton)
	if _, err := db.NewManager(c); err != nil {
		panic(err)
	}

	// Redis adapter (singleton)
	adapter := _redis.NewAdapter(c)

	// Queue dispatcher (singleton, thread-safe)
	c.BindSingleton("queue.dispatcher", func() any {
		return _redis.NewClient(adapter)
	})

	// Job handlers
	queue.RegisterJob("SendBalance", job.HandleSendBalanceJob)
}
