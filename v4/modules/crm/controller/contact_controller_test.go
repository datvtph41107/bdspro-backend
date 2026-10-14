package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"userv4/internal/httpapi"
	"userv4/modules/crm/dto"
)

func TestContactCreate_DecodeWithoutPathID(t *testing.T) {
	body := `{"phone":"0912322222","fullName":"Ko có","ownerOf":10,"ownerId":40,"productIds":[]}`
	req := httptest.NewRequest(http.MethodPost, "/v2/crm/contact", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	var got dto.ContactCreateHTTPRequest
	if err := httpapi.DecodePathJSON(req, &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Phone != "0912322222" || got.FullName != "Ko có" || got.OwnerOf != 10 || got.OwnerID != 40 {
		t.Fatalf("unexpected: %+v", got)
	}
}
