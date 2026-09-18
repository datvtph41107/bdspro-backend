package _utils

import "testing"

func TestEncodeDecodePhoneRoundTrip(t *testing.T) {
	phone := "0961745160"
	encoded, err := EncodePhone(phone)
	if err != nil {
		t.Fatalf("EncodePhone: %v", err)
	}
	if encoded == phone {
		t.Fatalf("encoded phone must differ from raw phone, got %q", encoded)
	}
	for _, ch := range encoded {
		if ch >= '0' && ch <= '9' {
			t.Fatalf("encoded phone must not contain raw digits, got %q", encoded)
		}
	}

	decoded, err := DecodePhone(encoded)
	if err != nil {
		t.Fatalf("DecodePhone: %v", err)
	}
	if decoded != phone {
		t.Fatalf("DecodePhone = %q, want %q", decoded, phone)
	}
}

func TestDecodePhoneRejectsRawPhone(t *testing.T) {
	if _, err := DecodePhone("0961745160"); err == nil {
		t.Fatal("expected error for raw phone in URL")
	}
}

func TestEncodePhoneRejectsNonDigits(t *testing.T) {
	if _, err := EncodePhone("+84961745160"); err == nil {
		t.Fatal("expected error for non-digit phone")
	}
}
