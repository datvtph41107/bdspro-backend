package controller

import (
	"context"
	"crypto/subtle"
	"net/http"
	"strings"
	"time"

	"userv4/internal/httpapi"
	"userv4/modules/seo/entity"
	"userv4/modules/seo/services"
)

const seoSecretHeader = "X-Seo-Secret-Key"

// InternalController serves /v2/internal/seo/* to the SEO renderer. Callers must send
// the shared key in X-Seo-Secret-Key; an empty configured key rejects every request.
type InternalController struct {
	svc    services.SeoDomainService
	secret string
}

func NewInternalController(svc services.SeoDomainService, secret string) *InternalController {
	return &InternalController{svc: svc, secret: strings.TrimSpace(secret)}
}

func (c *InternalController) RegisterHTTP(mux *http.ServeMux) {
	mux.Handle("GET /v2/internal/seo/render-data/by-slug/{slug}", c.guard(httpapi.Bind(c.renderDataBySlug)))
	mux.Handle("GET /v2/internal/seo/render-data/{refType}/{refId}", c.guard(httpapi.Bind(c.renderDataByRef)))
	mux.Handle("GET /v2/internal/seo/static-paths", c.guard(httpapi.Bind(c.staticPaths)))
}

func (c *InternalController) guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get(seoSecretHeader)
		if c.secret == "" || subtle.ConstantTimeCompare([]byte(got), []byte(c.secret)) != 1 {
			httpapi.WriteError(w, http.StatusUnauthorized, "secretKey không hợp lệ")
			return
		}
		next.ServeHTTP(w, r)
	})
}

type renderDataBySlugPath struct {
	Slug string `path:"slug"`
}

type renderDataByRefPath struct {
	RefType uint32 `path:"refType"`
	RefID   uint64 `path:"refId"`
}

type staticPathsQuery struct {
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
}

type seoStaticPathItem struct {
	ID           uint64  `json:"id"`
	Slug         string  `json:"slug"`
	CanonicalURL string  `json:"canonicalUrl"`
	RefType      uint32  `json:"refType"`
	RefID        *uint64 `json:"refId,omitempty"`
	UpdatedAt    string  `json:"updatedAt"`
	PageStatus   string  `json:"pageStatus"`
	RenderStatus string  `json:"renderStatus"`
	NeedGenerate bool    `json:"needGenerate"`
}

func (c *InternalController) renderDataBySlug(ctx context.Context, req renderDataBySlugPath) (*httpapi.Result, error) {
	seo, err := c.svc.GetRenderDataBySlug(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	return dataResult(renderData(seo))
}

func (c *InternalController) renderDataByRef(ctx context.Context, req renderDataByRefPath) (*httpapi.Result, error) {
	seo, err := c.svc.GetRenderDataByRef(ctx, req.RefType, req.RefID)
	if err != nil {
		return nil, err
	}
	return dataResult(renderData(seo))
}

func (c *InternalController) staticPaths(ctx context.Context, req staticPathsQuery) (*httpapi.Result, error) {
	items, total, err := c.svc.ListStaticPaths(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	out := make([]seoStaticPathItem, 0, len(items))
	for i := range items {
		seo := &items[i]
		out = append(out, seoStaticPathItem{
			ID:           seo.ID,
			Slug:         seo.Slug,
			CanonicalURL: seo.CanonicalURL,
			RefType:      uint32(seo.RefType),
			RefID:        seo.RefID,
			UpdatedAt:    seo.UpdatedAt.Format(time.RFC3339),
			PageStatus:   seo.PageStatus,
			RenderStatus: seo.RenderStatus,
			NeedGenerate: seo.NeedGenerate,
		})
	}
	return listResult(out, total)
}

func renderData(seo *entity.SeoDomain) map[string]any {
	return map[string]any{
		"seoDomain":     seo,
		"internalLinks": seo.InternalLinks,
		"relatives":     seo.Relatives,
	}
}
