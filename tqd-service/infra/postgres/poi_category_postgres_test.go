package postgres

import "testing"

func TestParentIDChanged(t *testing.T) {
	one := uint64(1)
	oneAgain := uint64(1)
	two := uint64(2)

	tests := []struct {
		name          string
		before, after *uint64
		want          bool
	}{
		{name: "both absent", want: false},
		{name: "parent added", after: &one, want: true},
		{name: "parent removed", before: &one, want: true},
		{name: "same parent", before: &one, after: &oneAgain, want: false},
		{name: "parent changed", before: &one, after: &two, want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := parentIDChanged(test.before, test.after); got != test.want {
				t.Fatalf("parentIDChanged() = %v, want %v", got, test.want)
			}
		})
	}
}
