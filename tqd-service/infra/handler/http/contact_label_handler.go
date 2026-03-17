package handler_http

import (
	"context"
	"net/http"
	"strconv"

	_dto "common/domain/dto"
	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/infra/validator"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ContactLabelHandler handles ContactLabel gRPC requests
type ContactLabelHandler struct {
	tqdpb.UnimplementedContactLabelServiceServer
	contactLabelUsecase usecase.ContactLabelUsecase
	mapper              *mapper.ContactLabelMapper
	validator           *validator.ContactLabelValidator
}

// NewContactLabelHandler creates a new ContactLabelHandler
func NewContactLabelHandler(
	contactLabelUsecase usecase.ContactLabelUsecase,
	mapper *mapper.ContactLabelMapper,
	validator *validator.ContactLabelValidator,
) *ContactLabelHandler {
	return &ContactLabelHandler{
		contactLabelUsecase: contactLabelUsecase,
		mapper:              mapper,
		validator:           validator,
	}
}

// CreateContactLabel creates a new contact label
// Swagger: route /v2/tqd/contact/labels [post]
// @Summary Create a new contact label
// @Description Create a new contact label
// @Tags Contact Label
// @Accept json
// @Produce json
// @Param request body tqdpb.ContactLabel true "Contact label"
// @Success 200 {object} tqdpb.ContactLabel "Contact label"
// @Router /contact/labels [post]
func (h *ContactLabelHandler) CreateContactLabel(ctx context.Context, req *tqdpb.ContactLabel) (*tqdpb.ContactLabel, error) {
	// Convert proto to DTO
	createReq := &dto.CreateContactLabelRequestDTO{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		IsActive:    req.IsActive,
		// SortOrder:   int(req.SortOrder),
		IsSystem: req.IsSystem,
		// Type:        int(req.Type),
	}

	// Validate request
	if err := h.validator.ValidateCreateRequest(createReq); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	// Create contact label
	result, err := h.contactLabelUsecase.Create(ctx, createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create contact label: %v", err)
	}

	// Convert DTO to proto
	response := &tqdpb.ContactLabel{
		Id:           result.ID,
		Name:         result.Name,
		Code:         result.Code,
		Description:  result.Description,
		Color:        result.Color,
		Icon:         result.Icon,
		IsActive:     result.IsActive,
		SortOrder:    int32(result.SortOrder),
		ContactCount: int32(result.ContactCount),
		IsSystem:     result.IsSystem,
		Status:       int32(result.Status),
		Type:         int32(result.Type),
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}

	return response, nil
}

// @Summary Get a contact label by ID
// @Description Get a contact label by ID
// @Tags Contact Label
// @Accept json
// @Produce json
// @Param id path string true "Contact label ID"
// @Success 200 {object} tqdpb.ContactLabel "Contact label"
// @Router /contact/labels/{id} [get]
func (h *ContactLabelHandler) GetContactLabel(ctx context.Context, req *tqdpb.GetContactLabelRequest) (*tqdpb.ContactLabel, error) {
	// Get contact label
	result, err := h.contactLabelUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get contact label: %v", err)
	}

	// Convert DTO to proto
	response := &tqdpb.ContactLabel{
		Id:           uint64(result.ID),
		Name:         result.Name,
		Code:         result.Code,
		Description:  result.Description,
		Color:        result.Color,
		Icon:         result.Icon,
		IsActive:     result.IsActive,
		SortOrder:    int32(result.SortOrder),
		ContactCount: int32(result.ContactCount),
		IsSystem:     result.IsSystem,
		Status:       int32(result.Status),
		Type:         int32(result.Type),
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}

	return response, nil
}

// @Summary Update a contact label
// @Description Update a contact label
// @Tags Contact Label
// @Accept json
// @Produce json
// @Param id path string true "Contact label ID"
// @Param request body tqdpb.ContactLabel true "Contact label"
// @Success 200 {object} tqdpb.ContactLabel "Contact label"
// @Router /contact/labels/{id} [put]
func (h *ContactLabelHandler) UpdateContactLabel(ctx context.Context, req *tqdpb.ContactLabel) (*tqdpb.ContactLabel, error) {
	// Convert proto to DTO
	updateReq := &dto.UpdateContactLabelRequestDTO{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		IsActive:    req.IsActive,
		// SortOrder:   int(req.SortOrder),
		IsSystem: req.IsSystem,
		// Type:     req.Type,
	}

	// Update contact label
	result, err := h.contactLabelUsecase.Update(ctx, req.Id, updateReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update contact label: %v", err)
	}

	// Convert DTO to proto
	response := &tqdpb.ContactLabel{
		Id:           uint64(result.ID),
		Name:         result.Name,
		Code:         result.Code,
		Description:  result.Description,
		Color:        result.Color,
		Icon:         result.Icon,
		IsActive:     result.IsActive,
		SortOrder:    int32(result.SortOrder),
		ContactCount: int32(result.ContactCount),
		IsSystem:     result.IsSystem,
		Status:       int32(result.Status),
		Type:         int32(result.Type),
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}

	return response, nil
}

// Swagger: route /v2/tqd/contact/labels/{id} [delete]
// @Summary Delete a contact label
// @Description Delete a contact label
// @Tags Contact Label
// @Accept json
// @Produce json
// @Param id path string true "Contact label ID"
// @Success 200 {object} sharepb.SubmitResponse "Submit response"
// @Router /contact/labels/{id} [delete]
func (h *ContactLabelHandler) DeleteContactLabel(ctx context.Context, req *tqdpb.DeleteContactLabelRequest) (*sharepb.SubmitResponse, error) {
	// Delete contact label
	err := h.contactLabelUsecase.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete contact label: %v", err)
	}

	// Return success response
	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Contact label deleted successfully",
	}, nil
}

// @Summary List contact labels with pagination
// @Description List contact labels with pagination
// @Tags Contact Label
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Param search query string false "Search"
// @Param isActive query bool false "Is active"
// @Param isSystem query bool false "Is system"
// @Success 200 {object} tqdpb.ListContactLabelsResponse "List contact labels response"
// @Router /contact/labels [get]
func (h *ContactLabelHandler) ListContactLabels(ctx context.Context, req *tqdpb.ListContactLabelsRequest) (*tqdpb.ListContactLabelsResponse, error) {
	// Convert proto to DTO
	listReq := &dto.ListContactLabelsRequestDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		Search:   req.Search,
		IsActive: &req.IsActive,
		IsSystem: &req.IsSystem,
	}

	// List contact labels
	result, err := h.contactLabelUsecase.List(ctx, listReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list contact labels: %v", err)
	}

	// Convert DTOs to protos
	contactLabels := make([]*tqdpb.ContactLabel, len(result.Data))
	for i, item := range result.Data {
		contactLabels[i] = &tqdpb.ContactLabel{
			Id:           uint64(item.ID),
			Name:         item.Name,
			Code:         item.Code,
			Description:  item.Description,
			Color:        item.Color,
			Icon:         item.Icon,
			IsActive:     item.IsActive,
			SortOrder:    int32(item.SortOrder),
			ContactCount: int32(item.ContactCount),
			IsSystem:     item.IsSystem,
			Status:       int32(item.Status),
			Type:         int32(item.Type),
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		}
	}

	// Return response
	return &tqdpb.ListContactLabelsResponse{
		Data:  contactLabels,
		Total: result.Total,
	}, nil
}

// HTTP Methods

// CreateContactLabelHTTP handles HTTP POST request for creating contact label
// @Summary Create a new contact label
// @Description Create a new contact label
// @Tags Contact Label
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tqdpb.ContactLabel true "Contact label"
// @Success 200 {object} tqdpb.ContactLabel "Contact label"
// @Router /v2/tqd/contact/labels [post]
func (h *ContactLabelHandler) CreateContactLabelHTTP(c *gin.Context) {
	var req tqdpb.ContactLabel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert proto to DTO
	createReq := &dto.CreateContactLabelRequestDTO{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		IsActive:    req.IsActive,
		// SortOrder:   int(req.SortOrder),
		IsSystem: req.IsSystem,
		// Type:        int(req.Type),
	}

	// Validate request
	if err := h.validator.ValidateCreateRequest(createReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create contact label
	result, err := h.contactLabelUsecase.Create(c.Request.Context(), createReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert DTO to proto
	response := &tqdpb.ContactLabel{
		Id:           result.ID,
		Name:         result.Name,
		Code:         result.Code,
		Description:  result.Description,
		Color:        result.Color,
		Icon:         result.Icon,
		IsActive:     result.IsActive,
		SortOrder:    int32(result.SortOrder),
		ContactCount: int32(result.ContactCount),
		IsSystem:     result.IsSystem,
		Status:       int32(result.Status),
		Type:         int32(result.Type),
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// GetContactLabelHTTP handles HTTP GET request for getting contact label by ID
// @Summary Get a contact label by ID
// @Description Get a contact label by ID
// @Tags Contact Label
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contact label ID"
// @Success 200 {object} tqdpb.ContactLabel "Contact label"
// @Router /v2/tqd/contact/labels/{id} [get]
func (h *ContactLabelHandler) GetContactLabelHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get contact label
	result, err := h.contactLabelUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contact label not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert DTO to proto
	response := &tqdpb.ContactLabel{
		Id:           result.ID,
		Name:         result.Name,
		Code:         result.Code,
		Description:  result.Description,
		Color:        result.Color,
		Icon:         result.Icon,
		IsActive:     result.IsActive,
		SortOrder:    int32(result.SortOrder),
		ContactCount: int32(result.ContactCount),
		IsSystem:     result.IsSystem,
		Status:       int32(result.Status),
		Type:         int32(result.Type),
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateContactLabelHTTP handles HTTP PUT request for updating contact label
// @Summary Update a contact label
// @Description Update a contact label
// @Tags Contact Label
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contact label ID"
// @Param request body tqdpb.ContactLabel true "Contact label"
// @Success 200 {object} tqdpb.ContactLabel "Contact label"
// @Router /v2/tqd/contact/labels/{id} [put]
func (h *ContactLabelHandler) UpdateContactLabelHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req tqdpb.ContactLabel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert proto to DTO
	updateReq := &dto.UpdateContactLabelRequestDTO{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		IsActive:    req.IsActive,
		// SortOrder:   int(req.SortOrder),
		IsSystem: req.IsSystem,
		// Type:     req.Type,
	}

	// Validate request
	if err := h.validator.ValidateUpdateRequest(updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update contact label
	result, err := h.contactLabelUsecase.Update(c.Request.Context(), id, updateReq)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contact label not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert DTO to proto
	response := &tqdpb.ContactLabel{
		Id:           result.ID,
		Name:         result.Name,
		Code:         result.Code,
		Description:  result.Description,
		Color:        result.Color,
		Icon:         result.Icon,
		IsActive:     result.IsActive,
		SortOrder:    int32(result.SortOrder),
		ContactCount: int32(result.ContactCount),
		IsSystem:     result.IsSystem,
		Status:       int32(result.Status),
		Type:         int32(result.Type),
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteContactLabelHTTP handles HTTP DELETE request for deleting contact label
// @Summary Delete a contact label
// @Description Delete a contact label
// @Tags Contact Label
// @Accept json
// @Produce json
// @Param id path string true "Contact label ID"
// @Success 200 {object} sharepb.SubmitResponse "Submit response"
// @Router /v2/tqd/contact/labels/{id} [delete]
func (h *ContactLabelHandler) DeleteContactLabelHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Delete contact label
	err = h.contactLabelUsecase.Delete(c.Request.Context(), id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contact label not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, &sharepb.SubmitResponse{
		Id:      id,
		Message: "Contact label deleted successfully",
	})
}

// ListContactLabelsHTTP handles HTTP GET request for listing contact labels
// @Summary List contact labels with pagination
// @Description List contact labels with pagination
// @Tags Contact Label
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Param search query string false "Search"
// @Param isActive query bool false "Is active"
// @Param isSystem query bool false "Is system"
// @Success 200 {object} tqdpb.ListContactLabelsResponse "List contact labels response"
// @Router /v2/tqd/contact/labels [get]
func (h *ContactLabelHandler) ListContactLabelsHTTP(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 32)
	size, _ := strconv.ParseUint(c.DefaultQuery("size", "10"), 10, 32)
	search := c.Query("search")

	var isActive *bool
	var isSystem *bool

	if isActiveStr := c.Query("isActive"); isActiveStr != "" {
		if val, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &val
		}
	}

	if isSystemStr := c.Query("isSystem"); isSystemStr != "" {
		if val, err := strconv.ParseBool(isSystemStr); err == nil {
			isSystem = &val
		}
	}

	// Convert to DTO
	listReq := &dto.ListContactLabelsRequestDTO{
		Pagable: _dto.Pagable{
			Page: uint32(page),
			Size: uint32(size),
		},
		Search:   search,
		IsActive: isActive,
		IsSystem: isSystem,
	}

	// List contact labels
	result, err := h.contactLabelUsecase.List(c.Request.Context(), listReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert DTOs to protos
	contactLabels := make([]*tqdpb.ContactLabel, len(result.Data))
	for i, item := range result.Data {
		contactLabels[i] = &tqdpb.ContactLabel{
			Id:           item.ID,
			Name:         item.Name,
			Code:         item.Code,
			Description:  item.Description,
			Color:        item.Color,
			Icon:         item.Icon,
			IsActive:     item.IsActive,
			SortOrder:    int32(item.SortOrder),
			ContactCount: int32(item.ContactCount),
			IsSystem:     item.IsSystem,
			Status:       int32(item.Status),
			Type:         int32(item.Type),
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		}
	}

	// Return response
	c.JSON(http.StatusOK, &tqdpb.ListContactLabelsResponse{
		Data:  contactLabels,
		Total: result.Total,
	})
}
