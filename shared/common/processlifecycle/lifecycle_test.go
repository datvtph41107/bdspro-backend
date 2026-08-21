package processlifecycle

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	_process "common/process"
)

type compatServer struct{}

func (compatServer) Serve(net.Listener) error { return errors.New("serve failed") }
func (compatServer) GracefulStop()            {}
func (compatServer) Stop()                    {}

func TestCompatibilityTypesAliasCanonicalProcess(t *testing.T) {
	var _ Server = compatServer{}
	var _ _process.Server = compatServer{}

	if ErrServerMissing != _process.ErrServerMissing ||
		ErrCancelMissing != _process.ErrCancelMissing ||
		ErrServeReturnedWithoutError != _process.ErrServeReturnedWithoutError {
		t.Fatal("compatibility errors do not share canonical process authority")
	}

	cfg := Config{GracefulStopTimeout: time.Second}
	canonical := _process.Config(cfg)
	if canonical.GracefulStopTimeout != time.Second {
		t.Fatalf("canonical config timeout = %s", canonical.GracefulStopTimeout)
	}
}

func TestCompatibilityRunDelegatesCanonicalProcess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := Run(ctx, cancel, compatServer{}, nil, nil, Config{GracefulStopTimeout: time.Millisecond})
	if err == nil || err.Error() != "serve: serve failed" {
		t.Fatalf("Run error = %v", err)
	}
}
