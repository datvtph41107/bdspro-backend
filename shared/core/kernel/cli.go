package kernel

import (
	"context"
	"core/kernel/container"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunCLI() {
	if len(os.Args) < 2 {
		fmt.Println("need a command: queue:work | igf:serve | ws:serve")
		return
	}
	cmd := os.Args[1]
	switch cmd {
	case "queue:work":
		// resolve queue and start worker
		q, err := container.ResolveTyped[interface{}](appC, "queue")
		if err != nil {
			panic(err)
		}
		switch qv := q.(type) {
		case interface{ StartWorker(ctx context.Context) }:
			ctx, cancel := context.WithCancel(context.Background())
			go qv.StartWorker(ctx)
			waitShutdown(cancel)
		default:
			fmt.Println("queue binding does not implement StartWorker(ctx)")
		}
	case "igf:serve":
		igf, err := container.ResolveTyped[interface{}](appC, "igf")
		if err != nil {
			panic(err)
		}
		if srv, ok := igf.(interface {
			Serve(ctx context.Context) error
		}); ok {
			ctx, cancel := context.WithCancel(context.Background())
			go srv.Serve(ctx)
			waitShutdown(cancel)
		} else {
			fmt.Println("igf binding invalid")
		}
	case "ws:serve":
		ws, err := container.ResolveTyped[interface{}](appC, "ws")
		if err != nil {
			panic(err)
		}
		if b, ok := ws.(interface {
			Serve(ctx context.Context) error
		}); ok {
			ctx, cancel := context.WithCancel(context.Background())
			go b.Serve(ctx)
			waitShutdown(cancel)
		} else {
			fmt.Println("ws binding invalid")
		}
	default:
		fmt.Println("unknown command", cmd)
	}
}

func waitShutdown(cancel context.CancelFunc) {
	// wait signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	fmt.Println("shutting down...")
	cancel()
	time.Sleep(500 * time.Millisecond)
}
