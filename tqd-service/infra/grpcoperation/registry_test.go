package grpcoperation

import (
	"common/businessoperation"
	"errors"
	"testing"
)

func TestRegistry(t *testing.T) {
	registry := NewRegistry()
	method := "/types.tqd.MapWorkspaceService/CreateGeneratedReport"

	if err := registry.Add(method, businessoperation.WorkspaceReportGenerate); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	operation, ok := registry.Find(method)
	if !ok || operation != businessoperation.WorkspaceReportGenerate {
		t.Fatalf("Find() = %q, %v", operation, ok)
	}

	if err := registry.Add(method, businessoperation.Code("parcel.view")); !errors.Is(err, ErrDuplicateMethod) {
		t.Fatalf("duplicate Add() error = %v", err)
	}
}
