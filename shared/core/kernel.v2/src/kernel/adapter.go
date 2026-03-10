package kernel

import "context"

type HookComponent struct {
	meta    Metadata
	startFn func(context.Context) error
	stopFn  func(context.Context) error
}

func NewHookComponent(meta Metadata, startFn, stopFn func(context.Context) error) *HookComponent {
	return &HookComponent{meta: meta, startFn: startFn, stopFn: stopFn}
}

func (c *HookComponent) Meta() Metadata {
	return c.meta
}

func (c *HookComponent) Start(ctx context.Context) error {
	if c.startFn == nil {
		return nil
	}
	return c.startFn(ctx)
}

func (c *HookComponent) Stop(ctx context.Context) error {
	if c.stopFn == nil {
		return nil
	}
	return c.stopFn(ctx)
}
