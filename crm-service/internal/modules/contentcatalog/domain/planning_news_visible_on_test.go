package domain

import "testing"

func TestPlanningNewsVisibleOn(t *testing.T) {
	t.Parallel()
	for _, value := range []PlanningNewsVisibleOn{10, 20, 30, 40, 50} {
		if !value.IsValid() || value.String() == "" || value.Label() == "" {
			t.Fatalf("expected %d to be valid with code and label", value)
		}
	}
	for _, value := range []PlanningNewsVisibleOn{0, 1, 9, 11, 49, 51, 99} {
		if value.IsValid() {
			t.Fatalf("expected %d to be invalid", value)
		}
	}
}

func TestParsePlanningNewsVisibleOn(t *testing.T) {
	t.Parallel()
	value, ok := ParsePlanningNewsVisibleOn(" MOST_READ ")
	if !ok || value != PlanningNewsVisibleOnMostRead {
		t.Fatalf("unexpected parse result: value=%d ok=%v", value, ok)
	}
}
