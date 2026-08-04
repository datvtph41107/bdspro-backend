package quotasubject

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestSubjectValidationAndKey(t *testing.T) {
	t.Parallel()

	subject := Subject{Kind: KindUser, ID: "42", Source: SourceActorProfile}
	if !subject.IsValid() {
		t.Fatal("user subject should be valid")
	}
	if got := subject.Key(); got != "user:42" {
		t.Fatalf("Key() = %q, want user:42", got)
	}

	invalid := Subject{Kind: KindUser, ID: "42", Source: SourceAPIKeyAuthority}
	if invalid.IsValid() {
		t.Fatal("mismatched source should be invalid")
	}
}

func TestSubjectContextAndMetadataRoundTrip(t *testing.T) {
	t.Parallel()

	want := Subject{Kind: KindProject, ID: "project_9", Source: SourceProjectAccess}
	ctx := WithSubject(context.Background(), want)
	got, ok := FromContext(ctx)
	if !ok || got != want {
		t.Fatalf("context subject = %+v, %v, want %+v", got, ok, want)
	}

	md := metadata.MD{}
	AppendMetadata(md, want)
	parsed, present, err := FromMetadata(md)
	if err != nil || !present || parsed != want {
		t.Fatalf("metadata subject = %+v, %v, %v, want %+v", parsed, present, err, want)
	}
}

func TestSubjectMetadataRejectsPartialOrDuplicateValues(t *testing.T) {
	t.Parallel()

	partial := metadata.Pairs(KindMetadataKey, string(KindUser))
	if _, _, err := FromMetadata(partial); err == nil {
		t.Fatal("partial subject metadata should fail")
	}

	duplicate := metadata.MD{
		KindMetadataKey:   {string(KindUser), string(KindUser)},
		IDMetadataKey:     {"42"},
		SourceMetadataKey: {string(SourceActorProfile)},
	}
	if _, _, err := FromMetadata(duplicate); err == nil {
		t.Fatal("duplicate subject metadata should fail")
	}
}

func TestWithSubjectDoesNotOverwriteExistingSubject(t *testing.T) {
	t.Parallel()

	first := Subject{Kind: KindUser, ID: "42", Source: SourceActorProfile}
	second := Subject{Kind: KindAPIKey, ID: "88", Source: SourceAPIKeyAuthority}
	ctx := WithSubject(context.Background(), first)
	ctx = WithSubject(ctx, second)

	got, ok := FromContext(ctx)
	if !ok || got != first {
		t.Fatalf("subject = %+v, %v, want %+v", got, ok, first)
	}
}
