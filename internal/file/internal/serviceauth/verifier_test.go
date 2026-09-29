package serviceauth

import (
	"strconv"
	"testing"
	"time"
)

func TestVerifierAcceptsCanonicalTimestampHMAC(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	v, ok := NewVerifier("test-file-service-secret")
	if !ok {
		t.Fatal("NewVerifier rejected non-empty secret")
	}
	timestamp := "1700000000"
	header := timestamp + ":" + sign(timestamp+":assistant-service", []byte("test-file-service-secret"))
	if !v.Verify(header, "assistant-service", now) {
		t.Fatal("canonical HMAC request rejected")
	}
}

func TestVerifierRejectsRawSecretAndMissingServiceName(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	v, _ := NewVerifier("test-file-service-secret")
	if v.Verify("test-file-service-secret", "assistant-service", now) {
		t.Fatal("raw secret must not be accepted as a wire credential")
	}
	timestamp := "1700000000"
	header := timestamp + ":" + sign(timestamp+":assistant-service", []byte("test-file-service-secret"))
	if v.Verify(header, "", now) {
		t.Fatal("missing service name must fail closed")
	}
}

func TestVerifierRejectsReplayFutureAndWrongService(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	v, _ := NewVerifier("test-file-service-secret")

	cases := []struct {
		name        string
		timestamp   int64
		signedName  string
		requestName string
	}{
		{name: "replay", timestamp: now.Add(-6 * time.Minute).Unix(), signedName: "assistant-service", requestName: "assistant-service"},
		{name: "future", timestamp: now.Add(2 * time.Minute).Unix(), signedName: "assistant-service", requestName: "assistant-service"},
		{name: "wrong service", timestamp: now.Unix(), signedName: "assistant-service", requestName: "tqd-service"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ts := time.Unix(tc.timestamp, 0).Format("150405")
			_ = ts // keep the test table readable; wire timestamp is decimal below.
			timestamp := formatInt(tc.timestamp)
			header := timestamp + ":" + sign(timestamp+":"+tc.signedName, []byte("test-file-service-secret"))
			if v.Verify(header, tc.requestName, now) {
				t.Fatalf("request unexpectedly accepted: %+v", tc)
			}
		})
	}
}

func TestNewVerifierRejectsEmptySecret(t *testing.T) {
	if _, ok := NewVerifier("   "); ok {
		t.Fatal("empty secret must be rejected")
	}
}

func formatInt(value int64) string {
	return strconv.FormatInt(value, 10)
}
