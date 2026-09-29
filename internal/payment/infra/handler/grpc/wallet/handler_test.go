package wallethandler

import "testing"

func TestExtractCommercialOrderReferenceCanonicalizesCase(t *testing.T) {
	got, ok := ExtractCommercialOrderReference(
		"bank transfer qhp-d9ee180c81828aa0",
	)
	if !ok || got != "QHP-D9EE180C81828AA0" {
		t.Fatalf("reference = %q, found = %v", got, ok)
	}
}
