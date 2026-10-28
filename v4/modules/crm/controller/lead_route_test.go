package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLeadBodiesUseProtoFieldNames(t *testing.T) {
	var stage leadSwitchStageReq
	if err := json.Unmarshal([]byte(`{"leadId":2,"stageId":9,"stageNote":"n"}`), &stage); err != nil {
		t.Fatal(err)
	}
	if stage.LeadID != 2 || stage.StageID != 9 || stage.StageNote != "n" {
		t.Fatalf("stage body: %+v", stage)
	}

	var assign leadAssignPersonReq
	body := `{"leadId":3,"ownerId":40,"ownerType":10,"chargePersonId":7,"chargePersonType":10,"noteAssign":"x"}`
	if err := json.Unmarshal([]byte(body), &assign); err != nil {
		t.Fatal(err)
	}
	if assign.LeadID != 3 || assign.OwnerID != 40 || assign.ChargePersonID != 7 || assign.NoteAssign != "x" {
		t.Fatalf("assign-person body: %+v", assign)
	}
}

func TestLeadManageRouteAcceptsOwnerOfName(t *testing.T) {
	mux := http.NewServeMux()
	NewLeadController(nil).RegisterHTTP(mux)

	for _, path := range []string{"/v2/crm/lead/manage/user/40/4", "/v2/crm/lead/manage/10/40/4"} {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if strings.Contains(rec.Body.String(), "invalid ownerOf") {
			t.Fatalf("%s: %d %s", path, rec.Code, rec.Body.String())
		}
	}
}
