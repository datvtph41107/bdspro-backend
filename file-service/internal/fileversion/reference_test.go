package fileversion

import (
	"strings"
	"testing"

	_utils "common/utils"
)

func TestReferenceCodecRoundTripAndTamperDetection(t *testing.T) {
	codec, err := NewReferenceCodec("xor-key", "signature-key")
	if err != nil {
		t.Fatal(err)
	}
	const path = "files/public/version/2026-08-26/abc/42.apk"

	reference, err := codec.Encode(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(reference, "v1.") {
		t.Fatalf("reference = %q, want v1 envelope", reference)
	}
	decoded, err := codec.Decode(reference)
	if err != nil {
		t.Fatal(err)
	}
	if decoded != path {
		t.Fatalf("decoded = %q, want %q", decoded, path)
	}

	tampered := reference[:len(reference)-1] + "A"
	if tampered == reference {
		tampered = reference[:len(reference)-1] + "B"
	}
	if _, err := codec.Decode(tampered); err == nil {
		t.Fatal("tampered version reference was accepted")
	}
}

func TestReferenceCodecReadsLegacyVersionReference(t *testing.T) {
	const xorKey = "xor-key"
	codec, err := NewReferenceCodec(xorKey, "signature-key")
	if err != nil {
		t.Fatal(err)
	}
	const path = "files/public/version/legacy/42.apk"
	legacy := "p" + _utils.XorEncode(path, xorKey)

	decoded, err := codec.Decode(legacy)
	if err != nil {
		t.Fatal(err)
	}
	if decoded != path {
		t.Fatalf("decoded = %q, want %q", decoded, path)
	}
}
