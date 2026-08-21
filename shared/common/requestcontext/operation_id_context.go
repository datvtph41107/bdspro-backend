package requestcontext

import (
	canonical "common/request"
	"context"
)

var ErrOperationIDContextConflict = canonical.ErrOperationIDContextConflict

func BindOperationID(ctx context.Context, operationID string) (context.Context, error) {
	return canonical.BindOperationID(ctx, operationID)
}

func OperationIDFromContext(ctx context.Context) (string, bool) {
	return canonical.OperationIDFromContext(ctx)
}
