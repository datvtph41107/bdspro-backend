package trust

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

func assessmentConfig(now time.Time, mode Mode) Config {
	return Config{
		ServiceID: "gateway-service",
		Secret:    "test-secret",
		Mode:      mode,
		MaxAge:    30 * time.Second,
		ClockSkew: 5 * time.Second,
		Now:       func() time.Time { return now },
	}
}

func TestAssessmentPromotesOnlyVerifiedMetadata(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := assessmentConfig(now, ModeAudit)
	md := metadata.Pairs("profileid", "42")
	if err := Sign("/test.Service/Call", md, cfg, testMetadataContract()); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	assessment := Assess("/test.Service/Call", md, cfg, testMetadataContract())
	if !assessment.Verified() || !assessment.CanPromotePrivilegedContext() {
		t.Fatalf("assessment = %+v", assessment)
	}
	if assessment.Reject() || assessment.LegacyUnverified() {
		t.Fatalf("verified metadata classified incorrectly: %+v", assessment)
	}
}

func TestAssessmentAuditMarksUnsignedPrivilegedMetadataLegacyOnly(t *testing.T) {
	assessment := Assess(
		"/test.Service/Call",
		metadata.Pairs("profileid", "42"),
		Config{Mode: ModeAudit},
		testMetadataContract(),
	)
	if assessment.Verified() || assessment.CanPromotePrivilegedContext() {
		t.Fatal("unsigned privileged metadata was promoted")
	}
	if assessment.Reject() || !assessment.LegacyUnverified() {
		t.Fatalf("assessment = %+v", assessment)
	}
}

func TestAssessmentEnforceRejectsUnsignedMetadata(t *testing.T) {
	assessment := Assess(
		"/test.Service/Call",
		metadata.Pairs("profileid", "42"),
		Config{Mode: ModeEnforce, Secret: "test-secret"},
		testMetadataContract(),
	)
	if !assessment.Reject() || assessment.LegacyUnverified() {
		t.Fatalf("assessment = %+v", assessment)
	}
}

func TestAssessmentAllowsOnlyDeclaredUnsignedInfrastructureMethods(t *testing.T) {
	for _, method := range []string{
		"/grpc.health.v1.Health/Check",
		"/grpc.health.v1.Health/Watch",
	} {
		assessment := Assess(method, metadata.MD{}, Config{
			Mode:   ModeEnforce,
			Secret: "test-secret",
		}, testMetadataContract())
		if assessment.Reject() || !assessment.UnsignedAllowed {
			t.Fatalf("method %s assessment = %+v", method, assessment)
		}
	}
}

func TestAssessmentRejectsIdentitySmugglingOnUnsignedMethod(t *testing.T) {
	assessment := Assess(
		"/grpc.health.v1.Health/Check",
		metadata.Pairs("profileid", "42"),
		Config{Mode: ModeEnforce, Secret: "test-secret"},
		testMetadataContract(),
	)
	if !assessment.Reject() || assessment.UnsignedAllowed {
		t.Fatalf("assessment = %+v", assessment)
	}
}

func TestLegacyUnverifiedContextCannotBeResigned(t *testing.T) {
	ctx, err := BindCaller(context.Background(), Caller{
		ServiceID: "legacy-unverified",
		Level:     LevelLegacyUnverified,
	})
	if err != nil {
		t.Fatal(err)
	}
	if CanPropagateLegacyPrivilegedContext(ctx) {
		t.Fatal("legacy-unverified context may launder privileged metadata")
	}
}

func TestVerifiedAndLocalContextsMayPropagateCompatibilityMetadata(t *testing.T) {
	if !CanPropagateLegacyPrivilegedContext(context.Background()) {
		t.Fatal("locally constructed compatibility context was blocked")
	}
	ctx, err := BindCaller(context.Background(), Caller{
		ServiceID: "gateway-service",
		Level:     LevelVerified,
		IssuedAt:  time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !CanPropagateLegacyPrivilegedContext(ctx) {
		t.Fatal("verified context was blocked")
	}
}
