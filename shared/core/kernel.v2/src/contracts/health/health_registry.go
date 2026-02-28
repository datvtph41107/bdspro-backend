package health

import "context"

type Registry interface {
	Register(name string, checker Checker)
	RunChecks(ctx context.Context) map[string]Status
}
