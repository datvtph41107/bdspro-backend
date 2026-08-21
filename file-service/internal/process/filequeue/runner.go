package filequeue

import (
	"context"
	"sync"
)

const DefaultWorkers = 10

// Runner owns the legacy FileService in-memory work queue. The queue is kept
// for compatibility, but worker lifetime belongs to process composition.
type Runner struct {
	queue   <-chan func()
	workers int
}

func New(queue <-chan func(), workers int) *Runner {
	if workers <= 0 {
		workers = DefaultWorkers
	}
	return &Runner{queue: queue, workers: workers}
}

// Run blocks until cancellation and all in-flight jobs have returned.
func (r *Runner) Run(ctx context.Context) {
	if r == nil || r.queue == nil {
		return
	}
	var wg sync.WaitGroup
	for i := 0; i < r.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					r.drainAccepted()
					return
				case job, ok := <-r.queue:
					if !ok {
						return
					}
					run(job)
				}
			}
		}()
	}
	wg.Wait()
}

func (r *Runner) drainAccepted() {
	for {
		select {
		case job, ok := <-r.queue:
			if !ok {
				return
			}
			run(job)
		default:
			return
		}
	}
}

func run(job func()) {
	if job != nil {
		job()
	}
}
