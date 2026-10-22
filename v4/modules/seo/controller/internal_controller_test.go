package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mdn4/modules/seo/entity"
	"mdn4/modules/seo/services"
)

type fakeRenderService struct {
	services.SeoDomainService
	slug string
}

func (f *fakeRenderService) GetRenderDataBySlug(_ context.Context, slug string) (*entity.SeoDomain, error) {
	f.slug = slug
	return &entity.SeoDomain{ID: 1, Slug: slug, InternalLinks: []entity.SeoLink{{ID: 2}}}, nil
}

func TestInternalController_RequiresSecret(t *testing.T) {
	cases := []struct {
		name, configured, header string
		want                     int
	}{
		{"missing header", "s3cret", "", http.StatusUnauthorized},
		{"wrong header", "s3cret", "nope", http.StatusUnauthorized},
		{"unconfigured rejects all", "", "", http.StatusUnauthorized},
		{"valid", "s3cret", "s3cret", http.StatusOK},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := &fakeRenderService{}
			mux := http.NewServeMux()
			NewInternalController(svc, tc.configured).RegisterHTTP(mux)

			r := httptest.NewRequest(http.MethodGet, "/v2/internal/seo/render-data/by-slug/quy-hoach-ha-noi", nil)
			if tc.header != "" {
				r.Header.Set("x-seo-secret-key", tc.header)
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)

			if w.Code != tc.want {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if tc.want == http.StatusOK {
				if svc.slug != "quy-hoach-ha-noi" || !strings.Contains(w.Body.String(), `"internalLinks"`) {
					t.Fatalf("slug=%q body=%s", svc.slug, w.Body.String())
				}
			} else if svc.slug != "" {
				t.Fatal("service must not be called without a valid secret")
			}
		})
	}
}

func TestInternalController_BySlugWinsOverRefPattern(t *testing.T) {
	mux := http.NewServeMux()
	NewInternalController(&fakeRenderService{}, "k").RegisterHTTP(mux)
	_, pattern := mux.Handler(httptest.NewRequest(http.MethodGet, "/v2/internal/seo/render-data/by-slug/x", nil))
	if pattern != "GET /v2/internal/seo/render-data/by-slug/{slug}" {
		t.Fatalf("pattern=%q", pattern)
	}
}
