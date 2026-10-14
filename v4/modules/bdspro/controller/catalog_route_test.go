package controller

import (
	"testing"

	"userv4/modules/bdspro/dto"
)

func TestCatalogSearchRequest_Pagable(t *testing.T) {
	req := dto.CatalogSearchRequest{Page: 0, Size: 0, Text: "abc"}
	p := req.Pagable()
	if p.GetPage() != 0 {
		t.Fatalf("page: %d", p.GetPage())
	}
	if p.GetSize() != 20 {
		t.Fatalf("default size: %d", p.GetSize())
	}
	if p.GetOffset() != 0 {
		t.Fatalf("offset: %d", p.GetOffset())
	}
}
