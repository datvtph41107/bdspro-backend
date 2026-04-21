package validate

import (
	"bdspro/core/pkg/validate"
	"testing"
)

func ptr[T any](v T) *T { return &v }

func TestViolations_AddMergeListClear(t *testing.T) {
	var v validate.Violations
	v.Add("title", validate.CodeRequired, "required")
	v.Addf("title", validate.CodeTooLong, "max %d", 10)

	if got, want := v.Len(), 2; got != want {
		t.Fatalf("Len() = %d, want %d", got, want)
	}

	list := v.List()
	if len(list) != 2 {
		t.Fatalf("List() len = %d, want 2", len(list))
	}

	list[0].Field = "mutated"
	if v.List()[0].Field == "mutated" {
		t.Fatal("List() must return a copy")
	}

	var other validate.Violations
	other.Add("lat", validate.CodeOutOfRange, "out of range")
	v.Merge(&other)

	if got, want := v.Len(), 3; got != want {
		t.Fatalf("after Merge, Len() = %d, want %d", got, want)
	}

	v.Clear()
	if v.HasErrors() {
		t.Fatal("Clear() must remove all violations")
	}
}

func TestViolations_ErrorString(t *testing.T) {
	var v validate.Violations
	v.Add("title", validate.CodeRequired, "required")

	got := v.Error()
	if got == "" {
		t.Fatal("Error() must not be empty")
	}
}

func TestViolations_AsError(t *testing.T) {
	var v validate.Violations
	if err := v.AsError(); err != nil {
		t.Fatal("empty violations must return nil")
	}

	v.Add("title", validate.CodeRequired, "required")
	if err := v.AsError(); err == nil {
		t.Fatal("non-empty violations must return error")
	}
}
