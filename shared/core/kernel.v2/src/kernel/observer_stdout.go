package kernel

import (
	"context"
	"fmt"
)

type StdoutObserver struct{}

func (StdoutObserver) OnStartBegin(ctx context.Context, meta Metadata) {
	fmt.Printf("[kernel] start begin component=%s\n", meta.ID)
}

func (StdoutObserver) OnStartDone(ctx context.Context, meta Metadata, d DurationInfo, err error) {
	if err != nil {
		fmt.Printf("[kernel] start failed component=%s duration_ms=%d err=%v\n", meta.ID, d.Milliseconds, err)
		return
	}
	fmt.Printf("[kernel] start done component=%s duration_ms=%d\n", meta.ID, d.Milliseconds)
}

func (StdoutObserver) OnStopBegin(ctx context.Context, meta Metadata) {
	fmt.Printf("[kernel] stop begin component=%s\n", meta.ID)
}

func (StdoutObserver) OnStopDone(ctx context.Context, meta Metadata, d DurationInfo, err error) {
	if err != nil {
		fmt.Printf("[kernel] stop failed component=%s duration_ms=%d err=%v\n", meta.ID, d.Milliseconds, err)
		return
	}
	fmt.Printf("[kernel] stop done component=%s duration_ms=%d\n", meta.ID, d.Milliseconds)
}
