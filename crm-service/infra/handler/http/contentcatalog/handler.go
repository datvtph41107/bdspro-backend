package contentcatalog

import (
	"net/http"
	"strconv"
	"strings"

	catalog_postgres "crm/infra/postgre/contentcatalog"

	"github.com/gin-gonic/gin"
)

type Handler struct{ repository *catalog_postgres.Repository }

func NewHandler(repository *catalog_postgres.Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) ListDomains(c *gin.Context) {
	items, err := h.repository.ListDomains(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "CONTENT_DOMAINS_FAILED", "message": err.Error()})
		return
	}
	c.Header("Cache-Control", "public, max-age=300, stale-while-revalidate=3600")
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *Handler) ListCategories(c *gin.Context) {
	items, err := h.repository.ListCategories(c.Request.Context(), c.Query("domain"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "CONTENT_CATEGORIES_FAILED", "message": err.Error()})
		return
	}
	c.Header("Cache-Control", "public, max-age=300, stale-while-revalidate=3600")
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *Handler) ListFeed(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	result, err := h.repository.ListFeed(
		c.Request.Context(),
		strings.TrimSpace(c.Param("feed")),
		c.Query("domain"),
		c.Query("category"),
		limit,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "CONTENT_FEED_INVALID", "message": err.Error()})
		return
	}
	c.Header("Cache-Control", "public, max-age=60, stale-while-revalidate=300")
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Search(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	kinds := []string{}
	for _, kind := range strings.Split(c.Query("kinds"), ",") {
		if v := strings.TrimSpace(kind); v != "" {
			kinds = append(kinds, v)
		}
	}
	result, err := h.repository.SearchPublishedContent(c.Request.Context(), c.Query("q"), kinds, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "CONTENT_SEARCH_FAILED", "message": err.Error()})
		return
	}
	c.Header("Cache-Control", "public, max-age=30, stale-while-revalidate=120")
	c.JSON(http.StatusOK, result)
}
