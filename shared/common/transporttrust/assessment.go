package transporttrust

import (
	"context"

	"google.golang.org/grpc/metadata"
)

/**
 * Assessment separates cryptographic verification from transport middleware.
 */
type Assessment struct {
	Caller                Caller
	Err                   error
	Mode                  Mode
	HasPrivilegedMetadata bool
	HasTrustAssertion     bool
	UnsignedAllowed       bool
}

/**
 * Assess evaluates one incoming metadata envelope.
 */
func Assess(method string, md metadata.MD, cfg Config) Assessment {
	caller, err := Verify(method, md, cfg)
	hasPrivilegedMetadata := HasPrivilegedMetadata(md)
	hasTrustAssertion := HasAssertion(md)
	return Assessment{
		Caller:                caller,
		Err:                   err,
		Mode:                  cfg.Mode,
		HasPrivilegedMetadata: hasPrivilegedMetadata,
		HasTrustAssertion:     hasTrustAssertion,
		UnsignedAllowed: isUnsignedInfrastructureMethod(method) &&
			!hasPrivilegedMetadata && !hasTrustAssertion,
	}
}

/** Verified reports whether the transport assertion is valid. */
func (a Assessment) Verified() bool {
	return a.Err == nil && a.Caller.Level == LevelVerified
}

/** Reject reports whether the request must fail at the transport boundary. */
func (a Assessment) Reject() bool {
	return !a.Verified() && !a.UnsignedAllowed && a.Mode == ModeEnforce
}

/** LegacyUnverified reports an auditable compatibility request. */
func (a Assessment) LegacyUnverified() bool {
	return !a.Verified() && a.Mode == ModeAudit &&
		(a.HasPrivilegedMetadata || a.HasTrustAssertion)
}

/** CanPromotePrivilegedContext gates canonical Actor/Caller/operation context. */
func (a Assessment) CanPromotePrivilegedContext() bool {
	return a.Verified()
}

/**
 * CanPropagateLegacyPrivilegedContext prevents an unsigned inbound identity
 * from being signed again by a downstream internal client.
 *
 * A context without transport provenance is considered locally constructed
 * compatibility context. A verified inbound context is also safe. Only an
 * explicitly legacy-unverified inbound context is blocked.
 */
func CanPropagateLegacyPrivilegedContext(ctx context.Context) bool {
	caller, ok := FromContext(ctx)
	return !ok || caller.Level == LevelVerified
}

func isUnsignedInfrastructureMethod(method string) bool {
	switch method {
	case "/grpc.health.v1.Health/Check", "/grpc.health.v1.Health/Watch":
		return true
	default:
		return false
	}
}
