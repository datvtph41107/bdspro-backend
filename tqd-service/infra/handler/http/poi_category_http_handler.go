package handler_http

import (
	"net/http"
	"strconv"

	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"github.com/gin-gonic/gin"
)

type PoiCategoryHandler struct {
	poiCategoryUsecase usecase.PoiCategoryUsecase
}

func NewPoiCategoryHandler(
	poiCategoryUsecase usecase.PoiCategoryUsecase,
) *PoiCategoryHandler {
	return &PoiCategoryHandler{
		poiCategoryUsecase: poiCategoryUsecase,
	}
}

// CreatePoiCategoryHTTP handles HTTP POST request for creating POI category
// @Summary Create a new POI category
// @Description Create a new POI category
// @Tags POI Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.PoiCategory true "POI category"
// @Success 200 {object} map[string]interface{} "POI category created successfully"
// @Router /v2/tqd/poi-categories [post]
func (h *PoiCategoryHandler) CreatePoiCategoryHTTP(c *gin.Context) {
	var category domain.PoiCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	createdCategory, err := h.poiCategoryUsecase.Create(c.Request.Context(), &category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create POI category",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POI category created successfully",
		"data":    createdCategory,
	})
}

// GetPoiCategoryHTTP handles HTTP GET request for getting POI category by ID
// @Summary Get POI category by ID
// @Description Get POI category by ID
// @Tags POI Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "POI category ID"
// @Success 200 {object} map[string]interface{} "POI category"
// @Router /v2/tqd/poi-categories/{id} [get]
func (h *PoiCategoryHandler) GetPoiCategoryHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid ID",
			"details": err.Error(),
		})
		return
	}

	category, err := h.poiCategoryUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "POI category not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POI category retrieved successfully",
		"data":    category,
	})
}

// UpdatePoiCategoryHTTP handles HTTP PUT request for updating POI category
// @Summary Update POI category
// @Description Update POI category
// @Tags POI Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "POI category ID"
// @Param request body domain.PoiCategory true "POI category"
// @Success 200 {object} map[string]interface{} "POI category updated successfully"
// @Router /v2/tqd/poi-categories/{id} [put]
func (h *PoiCategoryHandler) UpdatePoiCategoryHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid ID",
			"details": err.Error(),
		})
		return
	}

	var category domain.PoiCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	updatedCategory, err := h.poiCategoryUsecase.Update(c.Request.Context(), id, &category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update POI category",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POI category updated successfully",
		"data":    updatedCategory,
	})
}

// DeletePoiCategoryHTTP handles HTTP DELETE request for deleting POI category
// @Summary Delete POI category
// @Description Delete POI category
// @Tags POI Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "POI category ID"
// @Success 200 {object} map[string]interface{} "POI category deleted successfully"
// @Router /v2/tqd/poi-categories/{id} [delete]
func (h *PoiCategoryHandler) DeletePoiCategoryHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid ID",
			"details": err.Error(),
		})
		return
	}

	ok, err := h.poiCategoryUsecase.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to delete POI category",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "POI category deleted successfully",
		"data":    ok,
	})
}

// ListPoiCategoriesHTTP handles HTTP GET request for listing POI categories
// @Summary List POI categories
// @Description List POI categories with pagination and filters
// @Tags POI Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param code query string false "Filter by code"
// @Param name query string false "Filter by name"
// @Param parentId query int false "Filter by parent ID"
// @Param isActive query bool false "Filter by active status"
// @Success 200 {object} map[string]interface{} "POI categories list"
// @Router /v2/tqd/poi-categories [get]
func (h *PoiCategoryHandler) ListPoiCategoriesHTTP(c *gin.Context) {
	var req dto.PoiCategoryListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, total, err := h.poiCategoryUsecase.GetList(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to list POI categories",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "POI categories retrieved successfully",
		"data":    response,
		"total":   total,
	})
}

// GetPoiCategoryTreeHTTP handles HTTP GET request for getting POI categories in tree structure
// @Summary Get POI categories tree
// @Description Get POI categories in hierarchical tree structure (parent-children)
// @Tags POI Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "POI categories tree"
// @Router /v2/tqd/poi-categories/tree [get]
func (h *PoiCategoryHandler) GetPoiCategoryTreeHTTP(c *gin.Context) {
	treeData, errDTO := h.poiCategoryUsecase.GetTree(c.Request.Context())
	if errDTO != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    errDTO.Code,
			"message": errDTO.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POI categories tree retrieved successfully",
		"data":    treeData,
	})
}
