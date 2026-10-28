package dto

import (
	"testing"

	base_enum "base/enum"
)

func TestParseOwnerOf(t *testing.T) {
	cases := map[string]base_enum.EOwnerOf{
		"user":         base_enum.EOwnerOfMember,
		"User":         base_enum.EOwnerOfMember,
		"group":        base_enum.EOwnerOfGroup,
		"organization": base_enum.EOwnerOfOrgnization,
		"admin":        base_enum.EOwnerOfAdmin,
		"10":           base_enum.EOwnerOfMember,
		"30":           base_enum.EOwnerOfOrgnization,
	}
	for in, want := range cases {
		got, err := ParseOwnerOf(in)
		if err != nil || got != want {
			t.Fatalf("ParseOwnerOf(%q) = %d, %v; want %d", in, got, err, want)
		}
	}
	if _, err := ParseOwnerOf("team"); err == nil {
		t.Fatal("expected error for unknown ownerOf")
	}
}
