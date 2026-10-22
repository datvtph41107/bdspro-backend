package memoryusage

import (
	"common/operation"
	"context"
	"mdn4/modules/tqd/access"
	"mdn4/modules/tqd/services/usage"
	"strings"
	"testing"
)

func TestReadCommandRejectsPayloadReuseAndAcceptsOnlyOnce(t *testing.T) {
	store := NewStore()
	base := usage.ReadCommand{Subject: access.Subject{Type: access.SubjectProfile, ID: "42"}, Operation: operation.Code("workspace.parcel.detail.view"), CommandKey: "idem_read_1", RequestHash: strings.Repeat("a", 64)}
	if accepted, err := store.BindReadCommand(context.Background(), base); err != nil || accepted {
		t.Fatalf("bind=%v err=%v", accepted, err)
	}
	changed := base
	changed.RequestHash = strings.Repeat("b", 64)
	if _, err := store.BindReadCommand(context.Background(), changed); err != usage.ErrReadCommandConflict {
		t.Fatalf("changed payload err=%v", err)
	}
	if _, err := store.AcceptReadCommand(context.Background(), base, nil); err != nil {
		t.Fatal(err)
	}
	accepted, err := store.BindReadCommand(context.Background(), base)
	if err != nil || !accepted {
		t.Fatalf("replay bind=%v err=%v", accepted, err)
	}
}
