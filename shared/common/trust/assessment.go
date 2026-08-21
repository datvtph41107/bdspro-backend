package trust

import "google.golang.org/grpc/metadata"

// Assessment separates cryptographic verification from transport middleware.
type Assessment struct {
	Caller                Caller
	Err                   error
	Mode                  Mode
	HasPrivilegedMetadata bool
	HasTrustAssertion     bool
	UnsignedAllowed       bool
}

// Assess evaluates one incoming metadata envelope against an explicit contract.
func Assess(method string, md metadata.MD, cfg Config, contract MetadataContract) Assessment {
	caller, err := Verify(method, md, cfg, contract)
	hasPrivilegedMetadata := HasPrivilegedMetadata(md, contract)
	hasTrustAssertion := HasAssertion(md)
	return Assessment{
		Caller:                caller,
		Err:                   err,
		Mode:                  cfg.Mode,
		HasPrivilegedMetadata: hasPrivilegedMetadata,
		HasTrustAssertion:     hasTrustAssertion,
		UnsignedAllowed: contract.allowsUnsignedMethod(method) &&
			!hasPrivilegedMetadata && !hasTrustAssertion,
	}
}

func (a Assessment) Verified() bool {
	return a.Err == nil && a.Caller.Level == LevelVerified
}

func (a Assessment) Reject() bool {
	return !a.Verified() && !a.UnsignedAllowed && a.Mode == ModeEnforce
}

func (a Assessment) LegacyUnverified() bool {
	return !a.Verified() && a.Mode == ModeAudit &&
		(a.HasPrivilegedMetadata || a.HasTrustAssertion)
}

func (a Assessment) CanPromotePrivilegedContext() bool {
	return a.Verified()
}
