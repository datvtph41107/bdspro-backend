package usecase

import (
	"testing"

	"user/internal/domain/access"
)

func TestIsSystemRootRole(t *testing.T) {
	for _, test := range []struct {
		name string
		role *access.Role
		want bool
	}{
		{name: "nil", role: nil, want: false},
		{name: "ordinary", role: &access.Role{Key: "SUPPORT"}, want: false},
		{name: "canonical", role: &access.Role{Key: "QHPRO_SYSTEM_ROOT"}, want: true},
		{name: "normalized", role: &access.Role{Key: " qhpro_system_root "}, want: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := isSystemRootRole(test.role); got != test.want {
				t.Fatalf("isSystemRootRole() = %t, want %t", got, test.want)
			}
		})
	}
}
