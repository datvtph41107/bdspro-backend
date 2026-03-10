package kernel

import "context"

type Component interface {
	Meta() Metadata
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Observer interface {
	OnStartBegin(ctx context.Context, meta Metadata)
	OnStartDone(ctx context.Context, meta Metadata, d DurationInfo, err error)
	OnStopBegin(ctx context.Context, meta Metadata)
	OnStopDone(ctx context.Context, meta Metadata, d DurationInfo, err error)
}

type DurationInfo struct {
	Milliseconds int64
}
