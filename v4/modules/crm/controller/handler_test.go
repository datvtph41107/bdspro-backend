package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newRouteOnlyHandler() *Handler {
	return NewHandler(
		&ContactController{}, &AdminSalesController{}, &FriendController{}, &FriendGroupController{},
		&FollowController{}, &BlockController{}, &TagController{}, &PipelineController{},
		&StageController{}, &LeadController{}, &RuleController{}, &SharingController{},
		&InvitationController{}, &SupportTicketController{},
		&FeedbackController{}, &AppointmentController{}, &EnumController{})
}

func TestRegisterHTTP_ProtoPathsResolve(t *testing.T) {
	mux := http.NewServeMux()
	newRouteOnlyHandler().RegisterHTTP(mux)

	cases := []struct{ method, path, pattern string }{
		{http.MethodGet, "/v2/crm/contact/list/me", "GET /v2/crm/contact/list/me"},
		{http.MethodGet, "/v2/crm/contact/638/by-product?size=999", "GET /v2/crm/contact/{id}/by-product"},
		{http.MethodGet, "/v2/crm/contact/5/detail", "GET /v2/crm/contact/{id}/detail"},
		{http.MethodGet, "/v2/crm/contact/5/products", "GET /v2/crm/contact/{id}/products"},
		{http.MethodGet, "/v2/crm/contact/5/history", "GET /v2/crm/contact/{id}/history"},
		{http.MethodGet, "/v2/crm/contact/5/relation-ship", "GET /v2/crm/contact/{id}/relation-ship"},
		{http.MethodPut, "/v2/crm/contact/5/edit", "PUT /v2/crm/contact/{id}/edit"},
		{http.MethodPut, "/v2/crm/contact/5/note", "PUT /v2/crm/contact/{id}/note"},
		{http.MethodPut, "/v2/crm/contact/5/products", "PUT /v2/crm/contact/{id}/products"},
		{http.MethodDelete, "/v2/crm/contact/5", "DELETE /v2/crm/contact/{id}"},
		{http.MethodGet, "/v2/feedback/report/reasons", "GET /v2/feedback/report/reasons"},
		{http.MethodGet, "/v2/feedback/rate/detail/3", "GET /v2/feedback/rate/detail/{id}"},
		{http.MethodGet, "/v2/feedback/rate/product/stats/42", "GET /v2/feedback/rate/{registry}/stats/{ownerId}"},
		{http.MethodPost, "/v2/feedback/admin/report/3/approve", "POST /v2/feedback/admin/report/{id}/approve"},
		{http.MethodGet, "/v2/feedback/report-reason/3", "GET /v2/feedback/report-reason/{id}"},
		{http.MethodGet, "/v2/feedback/registry", "GET /v2/feedback/registry"},
		{http.MethodGet, "/v2.1/appointment/stats", "GET /v2.1/appointment/stats"},
		{http.MethodGet, "/v2.1/appointment/admin/list", "GET /v2.1/appointment/admin/list"},
		{http.MethodGet, "/v2.1/appointment/12", "GET /v2.1/appointment/{id}"},
		{http.MethodPut, "/v2.1/appointment/12/status/approve", "PUT /v2.1/appointment/{id}/status/{action}"},
		{http.MethodGet, "/v2/enums/crm/priority", "GET /v2/enums/crm/{name}"},
	}
	for _, tc := range cases {
		_, pattern := mux.Handler(httptest.NewRequest(tc.method, tc.path, nil))
		if pattern != tc.pattern {
			t.Errorf("%s %s matched %q, want %q", tc.method, tc.path, pattern, tc.pattern)
		}
	}
}