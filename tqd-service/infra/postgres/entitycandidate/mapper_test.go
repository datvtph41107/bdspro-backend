package entitycandidate

import (
	"testing"

	discoverydomain "tqd/internal/domain/discovery/model"
)

func TestBuildKeepsCanonicalLinksAndCapabilitiesConsistent(t *testing.T) {
	cases := []struct {
		name          string
		kind          discoverydomain.EntityKind
		row           Row
		wantDetail    string
		wantCanonical string
	}{
		{
			name:          "planning project",
			kind:          discoverydomain.EntityKindPlanningProject,
			row:           Row{ID: "12", Title: "Đồ án mẫu"},
			wantDetail:    "/tin-tuc/do-an-quy-hoach/12/do-an-mau",
			wantCanonical: "/tin-tuc/do-an-quy-hoach/12/do-an-mau",
		},
		{
			name:          "planning region",
			kind:          discoverydomain.EntityKindPlanningRegion,
			row:           Row{ID: "34", Title: "Vùng mẫu"},
			wantDetail:    "/tin-tuc/vung-quy-hoach/34/vung-mau",
			wantCanonical: "/tin-tuc/vung-quy-hoach/34/vung-mau",
		},
		{
			name:          "administrative ward",
			kind:          discoverydomain.EntityKindAdministrativeUnit,
			row:           Row{ID: "ward:9", Title: "Phường Mẫu", WardCode: "00009"},
			wantDetail:    "/tin-tuc/dia-ban/000092/phuong-mau",
			wantCanonical: "/tin-tuc/dia-ban/000092/phuong-mau",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			candidate := Build(test.kind, test.row, BuildOptions{MatchType: "test"})
			if candidate.Links.DetailURL != test.wantDetail {
				t.Fatalf("DetailURL = %q, want %q", candidate.Links.DetailURL, test.wantDetail)
			}
			if candidate.Links.CanonicalURL != test.wantCanonical {
				t.Fatalf("CanonicalURL = %q, want %q", candidate.Links.CanonicalURL, test.wantCanonical)
			}
			if !candidate.Capabilities.CanOpenDetail {
				t.Fatal("CanOpenDetail = false")
			}
		})
	}
}

func TestBuildPlanningProjectCapabilities(t *testing.T) {
	candidate := Build(discoverydomain.EntityKindPlanningProject, Row{ID: "12"}, BuildOptions{})
	if !candidate.Capabilities.CanCreateReport || !candidate.Capabilities.CanCompare {
		t.Fatalf("planning project capabilities = %#v", candidate.Capabilities)
	}
}
