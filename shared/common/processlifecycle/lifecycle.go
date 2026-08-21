// Package processlifecycle is a compatibility facade for common/process.
// New composition code must import common/process directly.
package processlifecycle

import (
	"context"
	"net"

	_process "common/process"
)

var (
	ErrServerMissing             = _process.ErrServerMissing
	ErrCancelMissing             = _process.ErrCancelMissing
	ErrServeReturnedWithoutError = _process.ErrServeReturnedWithoutError
)

type Server = _process.Server
type Actor = _process.Actor
type ActorFunc = _process.ActorFunc
type Config = _process.Config

func Run(
	ctx context.Context,
	cancel context.CancelFunc,
	server Server,
	listener net.Listener,
	actors []Actor,
	cfg Config,
) error {
	return _process.Run(ctx, cancel, server, listener, actors, cfg)
}
