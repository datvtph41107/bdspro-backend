package worker

import (
	"log"
	"os"

	"shared/core/bootstrap"
	"shared/core/container"
	"shared/core/queue"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: worker <queueName>")
	}
	queueName := os.Args[1]

	// Boot kernel with service register
	k := bootstrap.Boot("config/develop.yml", func(c *container.Container) {
		bdsbootstrap.Register(c)
	})

	// Create worker and register to runtime
	w := queue.NewWorker(k.Container, queueName)
	k.Runtime.Register(w)

	// Run runtime (blocks until signal)
	k.Runtime.Run()
}
