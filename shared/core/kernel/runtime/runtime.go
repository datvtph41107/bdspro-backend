package runtime

import (
	"context"
	"core/kernel/container"
	"log"
	"sync"
	"time"
)

type Server interface {
	Start(ctx context.Context, c *container.Container) error
	Stop(ctx context.Context) error
	Name() string
}

type Runtime struct {
	servers []Server
	jobs    []func(ctx context.Context, c *container.Container)
	wg      sync.WaitGroup
	stop    context.CancelFunc
}

func NewRuntime() *Runtime {
	return &Runtime{}
}

func (r *Runtime) AddServer(s Server) {
	r.servers = append(r.servers, s)
}

func (r *Runtime) AddJob(job func(ctx context.Context, c *container.Container)) {
	r.jobs = append(r.jobs, job)
}

func (r *Runtime) Start(c *container.Container) error {
	ctx, cancel := context.WithCancel(context.Background())
	r.stop = cancel

	// start servers
	for _, s := range r.servers {
		r.wg.Add(1)
		go func(srv Server) {
			defer r.wg.Done()
			log.Printf("runtime: starting server %s", srv.Name())
			if err := srv.Start(ctx, c); err != nil {
				log.Printf("runtime: server %s exited with error: %v", srv.Name(), err)
			} else {
				log.Printf("runtime: server %s stopped", srv.Name())
			}
		}(s)
	}

	// start background jobs
	for i, job := range r.jobs {
		r.wg.Add(1)
		go func(idx int, jb func(context.Context, *container.Container)) {
			defer r.wg.Done()
			log.Printf("runtime: starting job %d", idx)
			jb(ctx, c)
			log.Printf("runtime: job %d stopped", idx)
		}(i, job)
	}

	<-ctx.Done()

	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, s := range r.servers {
		_ = s.Stop(timeout)
	}
	r.wg.Wait()
	return nil
}

func (r *Runtime) Stop() {
	if r.stop != nil {
		r.stop()
	}
}
