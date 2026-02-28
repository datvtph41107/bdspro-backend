package simplehealth

import (
	"context"
	"core/kernel.v2/src/contracts/health"
)

type CheckerFunc func(ctx context.Context) health.Status

func (f CheckerFunc) Check(ctx context.Context) health.Status {
	return f(ctx)
}
