package grpcmetadata

import (
	"common/identity"
	"common/request"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"testing"
)

func TestFromHTTPRequestUsesCanonicalContextOnly(t *testing.T) {
	ctx := request.WithRequestID(context.Background(), "req_123")
	var err error
	ctx, err = request.BindOperationID(ctx, "op_123")
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = request.BindIdempotencyKey(ctx, "idem_123")
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = identity.BindActor(ctx, identity.Actor{
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
	ctx, err = identity.BindCaller(ctx, identity.Caller{Kind: identity.CallerUser})
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

func TestFromHTTPRequestProtectsSepayIPNSecret(t *testing.T) {
	request, err := http.NewRequest(http.MethodPost, "http://gateway.local/v2/payment/sepay/ipn", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("X-Secret-Key", "provider-secret")
	md := FromHTTPRequest(context.Background(), request)

	wantRaw := sha256.Sum256([]byte("provider-secret"))
	want := hex.EncodeToString(wantRaw[:])
	values := md.Get("x-qhpro-provider-credential-sha256")
	if len(values) != 1 || values[0] != want {
		t.Fatalf("provider digest = %v, want %q", values, want)
	}
	if values := md.Get("x-secret-key"); len(values) != 0 {
		t.Fatalf("raw provider secret must not be forwarded: %v", values)
	}
}
