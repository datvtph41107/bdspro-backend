package services

import "testing"

func TestCompareVersionStrings(t *testing.T) {
	cases := []struct {
		current, latest string
		want            int
	}{
		{"1.0.0", "1.0.1", -1},
		{"2.0.0", "1.9.9", 1},
		{"1.2.3", "1.2.3", 0},
		{"1.2", "1.2.0", 0},
		{"", "1.0.0", -1},
	}
	for _, tc := range cases {
		if got := compareVersionStrings(tc.current, tc.latest); got != tc.want {
			t.Fatalf("compareVersionStrings(%q, %q)=%d want %d", tc.current, tc.latest, got, tc.want)
		}
	}
}
