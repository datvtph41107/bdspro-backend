package services

import (
	"context"
	"testing"
)

func TestLocalPublicService_GetEnums(t *testing.T) {
	svc := NewLocalPublicService(nil)
	resp, err := svc.GetEnums(context.Background(), 0)
	if err != nil {
		t.Fatalf("GetEnums: %v", err)
	}
	if resp == nil || !resp.GetLastest() || len(resp.GetEnums()) == 0 {
		t.Fatalf("expected latest enums catalog, got %+v", resp)
	}
	keys := map[string]bool{}
	for _, item := range resp.GetEnums() {
		keys[item.GetKey()] = true
		if len(item.GetData()) == 0 {
			t.Fatalf("enum %q has empty data", item.GetKey())
		}
	}
	for _, want := range []string{
		"problem_report",
		"legal_status",
		"planning_type",
		"location_unit_type",
	} {
		if !keys[want] {
			t.Fatalf("missing enum key %q in %#v", want, keys)
		}
	}
}

func TestLocalPublicService_GetEnums_Unchanged(t *testing.T) {
	svc := NewLocalPublicService(nil)
	first, err := svc.GetEnums(context.Background(), 0)
	if err != nil {
		t.Fatalf("first: %v", err)
	}
	second, err := svc.GetEnums(context.Background(), first.GetTimestamp())
	if err != nil {
		t.Fatalf("second: %v", err)
	}
	if second.GetLastest() || len(second.GetEnums()) != 0 {
		t.Fatalf("expected empty not-latest response, got %+v", second)
	}
}
