package grpcmetadata

import (
	"common/callercontext"
	"common/identitycontext"
	"common/requestcontext"
	"context"
	"net/http"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestFromHTTPRequestUsesCanonicalContextOnly(t *testing.T) {
	ctx := requestcontext.WithRequestID(context.Background(), "req_123")
	var err error
	ctx, err = requestcontext.BindOperationID(ctx, "op_123")
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = requestcontext.BindIdempotencyKey(ctx, "idem_123")
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = identitycontext.BindActor(ctx, identitycontext.Actor{
		AuthID:    10,
		ProfileID: 42,
		OriginID:  20,
		SessionID: 30,
		Role:      "USER",
		TokenType: "access",
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = callercontext.BindCaller(ctx, callercontext.Caller{Kind: callercontext.KindUser})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://gateway.local/v2/tqd/map-workspace/reports", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer raw-token-must-not-be-forwarded")
	req.Header.Set("X-API-Key", "raw-api-key-must-not-be-forwarded")
	req.Header.Set("device-id", "device-1")

	md := FromHTTPRequest(ctx, req)

	wantSingle := map[string]string{
		"x-request-id":        "req_123",
		"x-operation-id":      "op_123",
		"idempotency-key":     "idem_123",
		"authid":              "10",
		"profileid":           "42",
		"originid":            "20",
		"session":             "30",
		"role":                "USER",
		"type":                "access",
		"x-qhpro-caller-kind": "user",
		"device-id":           "device-1",
	}
	for key, want := range wantSingle {
		values := md.Get(key)
		if len(values) != 1 || values[0] != want {
			t.Fatalf("%s = %v, want %q", key, values, want)
		}
	}

	for _, forbidden := range []string{"authorization", "x-api-key", "planid", "planfrom"} {
		if values := md.Get(forbidden); len(values) != 0 {
			t.Fatalf("%s must not be forwarded: %v", forbidden, values)
		}
	}
}

func TestDuringReportMigrationUsesTargetAdapterOnlyForEnabledTargetRoute(t *testing.T) {
	legacyCalls := 0
	legacy := func(_ context.Context, _ *http.Request) metadata.MD {
		legacyCalls++
		return metadata.Pairs("bridge", "legacy")
	}

	tests := []struct {
		name          string
		enabled       bool
		method        string
		path          string
		wantLegacy    bool
		wantOperation string
	}{
		{
			name:       "flag disabled keeps target route legacy",
			enabled:    false,
			method:     http.MethodPost,
			path:       createGeneratedReportPath,
			wantLegacy: true,
		},
		{
			name:          "enabled target route uses canonical adapter",
			enabled:       true,
			method:        http.MethodPost,
			path:          createGeneratedReportPath + "/",
			wantOperation: "op_target",
		},
		{
			name:       "enabled non-target route stays legacy",
			enabled:    true,
			method:     http.MethodGet,
			path:       "/v2/tqd/regions",
			wantLegacy: true,
		},
		{
			name:       "enabled wrong method stays legacy",
			enabled:    true,
			method:     http.MethodGet,
			path:       createGeneratedReportPath,
			wantLegacy: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := legacyCalls
			ctx, err := requestcontext.BindOperationID(context.Background(), "op_target")
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequestWithContext(ctx, tt.method, "http://gateway.local"+tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			md := DuringReportMigration(tt.enabled, legacy)(ctx, req)
			if tt.wantLegacy {
				if got := md.Get("bridge"); len(got) != 1 || got[0] != "legacy" {
					t.Fatalf("legacy bridge metadata = %v", got)
				}
				if legacyCalls != before+1 {
					t.Fatalf("legacy bridge calls = %d, want %d", legacyCalls, before+1)
				}
				return
			}

			if legacyCalls != before {
				t.Fatalf("legacy bridge must not be called for target route")
			}
			if got := md.Get("x-operation-id"); len(got) != 1 || got[0] != tt.wantOperation {
				t.Fatalf("operation metadata = %v, want %q", got, tt.wantOperation)
			}
		})
	}
}

func TestReportTargetEnabled(t *testing.T) {
	for _, value := range []string{"1", "true", "TRUE", " yes "} {
		t.Run(value, func(t *testing.T) {
			t.Setenv("QHPRO_REPORT_TARGET_ENABLED", value)
			if !ReportTargetEnabled() {
				t.Fatalf("value %q must enable target", value)
			}
		})
	}

	for _, value := range []string{"", "0", "false", "enabled"} {
		t.Run("disabled_"+value, func(t *testing.T) {
			t.Setenv("QHPRO_REPORT_TARGET_ENABLED", value)
			if ReportTargetEnabled() {
				t.Fatalf("value %q must not enable target", value)
			}
		})
	}
}
