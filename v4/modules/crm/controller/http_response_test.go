package controller

import (
	"testing"

	domain "mdn4/modules/crm/entity"
)

func TestContactDetailResult_IncludesProductIds(t *testing.T) {
	entity := &domain.Contact{
		ProductIDs: []uint64{10, 20},
	}
	entity.ID = 1
	entity.FullName = "Test"
	entity.Phone = "0900000000"

	res, err := contactDetailResult(entity)
	if err != nil {
		t.Fatalf("contactDetailResult: %v", err)
	}
	m, ok := res.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected map data, got %T", res.Data)
	}
	raw, ok := m["productIds"]
	if !ok {
		t.Fatalf("productIds missing: %+v", m)
	}
	ids, ok := raw.([]uint64)
	if !ok {
		// json round-trip via proto may use []any
		if anySlice, ok := raw.([]any); ok && len(anySlice) == 2 {
			return
		}
		t.Fatalf("unexpected productIds type %T", raw)
	}
	if len(ids) != 2 || ids[0] != 10 || ids[1] != 20 {
		t.Fatalf("unexpected productIds: %v", ids)
	}
}

func TestContactDetailResult_IncludesProfileId(t *testing.T) {
	pid := uint64(267)
	entity := &domain.Contact{ProfileID: &pid, ProductIDs: []uint64{}}
	entity.ID = 1
	entity.FullName = "Test"
	entity.Phone = "0961745160"

	res, err := contactDetailResult(entity)
	if err != nil {
		t.Fatalf("contactDetailResult: %v", err)
	}
	m := res.Data.(map[string]any)
	if m["profileId"] != pid {
		t.Fatalf("expected profileId %d, got %+v", pid, m["profileId"])
	}
}

func TestContactDetailResult_EmptyProductIds(t *testing.T) {
	entity := &domain.Contact{ProductIDs: nil}
	entity.ID = 1
	entity.FullName = "Test"
	entity.Phone = "0900000000"

	res, err := contactDetailResult(entity)
	if err != nil {
		t.Fatalf("contactDetailResult: %v", err)
	}
	m := res.Data.(map[string]any)
	raw, ok := m["productIds"]
	if !ok {
		t.Fatalf("productIds missing: %+v", m)
	}
	ids, ok := raw.([]uint64)
	if !ok || len(ids) != 0 {
		t.Fatalf("expected empty productIds, got %v (%T)", raw, raw)
	}
}
