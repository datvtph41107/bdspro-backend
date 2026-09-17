package filehttp

import (
	"file/internal/filecontent"
	"file/internal/serviceauth"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOwnedLookupRejectsUserCredentialAndNeverExposesCompletedEffect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	verifier, _ := serviceauth.NewVerifier("local-unit-test-service-key")
	handler := NewHandler(&fakeFileOperations{saved: filecontent.SavedFile{Path: "private-owner-path"}}, verifier)
	router := gin.New()
	router.GET("/internal/v1/file/owned", handler.FindOwnedFile)
	for _, credential := range []string{"", "Bearer user-token", "local-unit-test-service-key"} {
		request := httptest.NewRequest(http.MethodGet, "/internal/v1/file/owned?owner_namespace=generated_report&owner_key=job-1", nil)
		request.Header.Set("Authorization", credential)
		request.Header.Set("X-Service-Auth", credential)
		request.Header.Set("X-Service-Name", "tqd")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != http.StatusUnauthorized || response.Body.String() == "private-owner-path" {
			t.Fatalf("user/raw credential exposed effect: %d %s", response.Code, response.Body.String())
		}
	}
}
