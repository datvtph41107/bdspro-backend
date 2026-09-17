package domain

import "testing"

func TestOrderReferenceCanonicalOwner(t *testing.T) {
	entropy := [8]byte{0x0b, 0xdf, 0x8b, 0x7a, 0xff, 0x3a, 0xa0, 0xa2}
	canonical := FormatOrderReference(entropy)
	if canonical != "QHP-0BDF8B7AFF3AA0A2" || !IsCanonicalOrderReference(canonical) {
		t.Fatalf("canonical reference = %q", canonical)
	}

	for _, input := range []string{
		"bank transfer qhp-0bdf8b7aff3aa0a2",
		"BankAPINotify QHP0BDF8B7AFF3AA0A2",
	} {
		got, ok := ExtractOrderReference(input)
		if !ok || got != canonical {
			t.Fatalf("ExtractOrderReference(%q) = %q, %v", input, got, ok)
		}
	}

	if _, ok := ExtractOrderReference("QHP-TOO-SHORT"); ok {
		t.Fatal("invalid reference accepted")
	}
}
