package fileauthorization

import (
	"context"
	"errors"
	"fmt"
)

// OrdinaryCheck evaluates the durable per-file access path. Global permission
// authority is consulted only when this path does not already allow the read.
type OrdinaryCheck func() (bool, error)

// DeniedAuthorizer is used for requests that do not carry an authenticated
// actor. Anonymous signed-file access remains supported, but anonymous callers
// never fan out to Auth to ask for a global bypass.
type DeniedAuthorizer struct{}

func (DeniedAuthorizer) Require(context.Context, Capability) error { return ErrDenied }

// AuthorizePrivateRead preserves the two independent authorization mechanisms:
// per-file durable access first, then the explicit global IAM capability.
//
// A global capability can intentionally bypass a broken/missing file_access row,
// preserving the established global operator fallback. If global authority is unavailable,
// the bypass fails closed rather than falling back to a role claim.
func AuthorizePrivateRead(
	ctx context.Context,
	ordinary OrdinaryCheck,
	global Authorizer,
) error {
	if ordinary == nil {
		return errors.New("ordinary private-file access checker is required")
	}
	allowed, ordinaryErr := ordinary()
	if allowed {
		return nil
	}
	if global == nil {
		global = UnavailableAuthorizer{}
	}
	if err := global.Require(ctx, PrivateReadAny); err != nil {
		if errors.Is(err, ErrDenied) {
			if ordinaryErr != nil {
				return fmt.Errorf("ordinary private-file access check: %w", ordinaryErr)
			}
			return ErrDenied
		}
		return ErrUnavailable
	}
	return nil
}
