package publiccontent

import (
	"net/http"
	"strconv"
	"strings"

	publiccontent_postgres "tqd/infra/postgres/publiccontent"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repository *publiccontent_postgres.Repository
}

func NewHandler(repository *publiccontent_postgres.Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) GetParcelQuickView(c *gin.Context) {
	parcelID, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil || parcelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_PARCEL_ID",
			"message": "Parcel id must be a positive integer.",
		})
		return
	}

	result, err := h.repository.GetParcelQuickView(c.Request.Context(), parcelID)
	if err != nil {
		if publiccontent_postgres.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "PARCEL_NOT_FOUND",
				"message": "Không tìm thấy thửa đất.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "PARCEL_QUICK_VIEW_FAILED",
			"message": err.Error(),
		})
		return
	}

	c.Header("Cache-Control", "public, max-age=30, stale-while-revalidate=120")
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetAdministrativeUnitProjection(c *gin.Context) {
	identity := strings.TrimSpace(c.Param("identity"))
	if identity == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_ADMINISTRATIVE_UNIT_IDENTITY",
			"message": "Administrative unit identity is required.",
		})
		return
	}

	result, err := h.repository.GetAdministrativeUnitProjection(c.Request.Context(), identity)
	if err != nil {
		if publiccontent_postgres.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "ADMINISTRATIVE_UNIT_NOT_FOUND",
				"message": "Không tìm thấy đơn vị hành chính.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "ADMINISTRATIVE_UNIT_PROJECTION_FAILED",
			"message": err.Error(),
		})
		return
	}

	c.Header("Cache-Control", "public, max-age=60, stale-while-revalidate=300")
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetPlanningProjectProjection(c *gin.Context) {
	identity := strings.TrimSpace(c.Param("identity"))
	if identity == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_PLANNING_PROJECT_IDENTITY",
			"message": "Planning project identity is required.",
		})
		return
	}

	result, err := h.repository.GetPlanningProjectProjection(c.Request.Context(), identity)
	if err != nil {
		if publiccontent_postgres.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "PLANNING_PROJECT_NOT_FOUND",
				"message": "Không tìm thấy đồ án quy hoạch.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "PLANNING_PROJECT_PROJECTION_FAILED",
			"message": err.Error(),
		})
		return
	}

	c.Header("Cache-Control", "public, max-age=60, stale-while-revalidate=300")
	c.JSON(http.StatusOK, result)
}
