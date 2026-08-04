package transporttrust

import (
	_quotasubject "common/quotasubject"

	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

func TestVerifyRejectsTamperedQuotaSubject(t *testing.T) {
	t.Parallel()

	now := time.Unix(1_700_000_000, 0)
	cfg := Config{ServiceID: "gateway-service", Secret: "test-secret", Now: func() time.Time { return now }}
	md := metadata.MD{}
	_quotasubject.AppendMetadata(md, _quotasubject.Subject{
		Kind: _quotasubject.KindUser, ID: "42", Source: _quotasubject.SourceActorProfile,
	})
	if err := Sign("/test.Service/Call", md, cfg); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	md.Set(_quotasubject.IDMetadataKey, "99")

	_, err := Verify("/test.Service/Call", md, cfg)
	if !errors.Is(err, ErrInvalidAssertion) {
		t.Fatalf("Verify() error = %v, want ErrInvalidAssertion", err)
	}
}

func TestQuotaSubjectMetadataIsPrivileged(t *testing.T) {
	t.Parallel()

	md := metadata.Pairs(_quotasubject.KindMetadataKey, string(_quotasubject.KindUser))
	if !HasPrivilegedMetadata(md) {
		t.Fatal("quota subject metadata should be privileged")
	}
}
