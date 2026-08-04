package operationcontext

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestBindIsIdempotentAndImmutable(t *testing.T) {
	t.Parallel()

	selection, err := New("parcel.lookup", "parcel.lookup_transport", SourceDirectBinding)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	ctx, err := Bind(context.Background(), selection)
	if err != nil {
		t.Fatalf("Bind() error = %v", err)
	}
	if _, err := Bind(ctx, selection); err != nil {
		t.Fatalf("Bind() same selection error = %v", err)
	}

	other, _ := New("parcel.planning_view", "parcel.planning_transport", SourceDomainResolver)
	if _, err := Bind(ctx, other); !errors.Is(err, ErrOperationAlreadySelected) {
		t.Fatalf("Bind() error = %v, want ErrOperationAlreadySelected", err)
	}
}

func TestMetadataRoundTrip(t *testing.T) {
	t.Parallel()

	selection, _ := New("parcel.lookup", "parcel.lookup_transport", SourceDirectBinding)
	md := metadata.MD{}
	AppendMetadata(md, selection)

	got, ok, err := FromMetadata(md)
	if err != nil {
		t.Fatalf("FromMetadata() error = %v", err)
	}
	if !ok || got != selection {
		t.Fatalf("FromMetadata() = %#v, %v", got, ok)
	}
}

func TestMetadataRejectsPartialSelection(t *testing.T) {
	t.Parallel()

	md := metadata.Pairs(CodeMetadataKey, "parcel.lookup")
	if _, _, err := FromMetadata(md); err == nil {
		t.Fatal("FromMetadata() error = nil, want incomplete metadata error")
	}
}
