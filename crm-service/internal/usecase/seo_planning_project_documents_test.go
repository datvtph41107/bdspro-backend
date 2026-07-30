package usecase

import (
	"reflect"
	"testing"

	"crm/internal/dto"
)

func TestPlanningDocumentItemsPreservesDuplicateLabelsByID(
	t *testing.T,
) {
	got := planningDocumentItems([]dto.TqdPlanningDocument{
		{
			ID:               9,
			DocumentTypeName: "Tệp khác",
			Title:            "Không xác định",
		},
		{
			ID:               8,
			DocumentTypeName: "Tệp khác",
			Title:            "Không xác định",
		},
	})

	want := []string{
		"Tệp khác: Không xác định · Hồ sơ #9",
		"Tệp khác: Không xác định · Hồ sơ #8",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected documents:\ngot:  %#v\nwant: %#v", got, want)
	}
}

func TestPlanningDocumentItemsKeepsDistinctLabelsClean(
	t *testing.T,
) {
	got := planningDocumentItems([]dto.TqdPlanningDocument{
		{
			ID:               1,
			DocumentTypeName: "Quyết định",
			Title:            "Quyết định phê duyệt",
		},
		{
			ID:               2,
			DocumentTypeName: "Bản vẽ",
			Title:            "Sơ đồ sử dụng đất",
		},
	})

	want := []string{
		"Quyết định: Quyết định phê duyệt",
		"Bản vẽ: Sơ đồ sử dụng đất",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected documents:\ngot:  %#v\nwant: %#v", got, want)
	}
}
