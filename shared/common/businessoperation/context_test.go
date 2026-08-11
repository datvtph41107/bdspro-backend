package businessoperation

import (
	"context"
	"errors"
	"testing"
)

func TestBind(t *testing.T) {
	ctx := context.Background()

	bound, err := Bind(ctx, WorkspaceReportGenerate)
	if err != nil {
		t.Fatalf("Bind() error = %v", err)
	}

	operation, ok := FromContext(bound)
	if !ok || operation != WorkspaceReportGenerate {
		t.Fatalf("FromContext() = %q, %v", operation, ok)
	}

	same, err := Bind(bound, WorkspaceReportGenerate)
	if err != nil {
		t.Fatalf("Bind() same value error = %v", err)
	}
	if same != bound {
		t.Fatal("same operation should keep the existing context")
	}

	_, err = Bind(bound, Code("parcel.view"))
	if !errors.Is(err, ErrContextConflict) {
		t.Fatalf("Bind() conflict error = %v", err)
	}
}
