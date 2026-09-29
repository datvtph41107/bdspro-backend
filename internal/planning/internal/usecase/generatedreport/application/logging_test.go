package application

import (
	"strings"
	"testing"
)

func TestCommandRefIsStableAndDoesNotExposeRawKey(t *testing.T) {
	raw := "idem-secret-command"

	first := commandRef(raw)
	second := commandRef(raw)

	if first == "" || first != second {
		t.Fatalf("commandRef unstable: %q %q", first, second)
	}
	if strings.Contains(first, raw) {
		t.Fatalf("commandRef leaked raw key: %q", first)
	}
	if commandRef("other-command") == first {
		t.Fatal("different command keys produced same test fingerprint")
	}
}
