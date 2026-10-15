package mapper

import (
	_entity "common/domain/entity"
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"time"
	"userv4/modules/tqd/dto"
	domain "userv4/modules/tqd/entity"
)

type PoiCategoryMapper struct{}

func NewPoiCategoryMapper() *PoiCategoryMapper {
	return &PoiCategoryMapper{}
}

func (m *PoiCategoryMapper) ToResponse(pc *domain.PoiCategory) *dto.PoiCategoryResponse {
	if pc == nil {
		return nil
	}
	resp := &dto.PoiCategoryResponse{}
	resp.FromDomain(pc)
	resp.IconPath = m.resolveIconPath(pc.Icon)
	return resp
}

func (m *PoiCategoryMapper) ToResponseList(categories []domain.PoiCategory) []dto.PoiCategoryResponse {
	if categories == nil {
		return nil
	}
	result := make([]dto.PoiCategoryResponse, len(categories))
	for i, pc := range categories {
		result[i].FromDomain(&pc)
		result[i].IconPath = m.resolveIconPath(pc.Icon)
	}
	return result
}

func (m *PoiCategoryMapper) ToTreeResponse(pc *domain.PoiCategory) *dto.PoiCategoryTreeResponse {
	if pc == nil {
		return nil
	}
	resp := &dto.PoiCategoryTreeResponse{}
	resp.FromDomainTree(pc)
	resp.IconPath = m.resolveIconPath(pc.Icon)
	return resp
}

func (m *PoiCategoryMapper) ToTreeResponseList(categories []domain.PoiCategory) []dto.PoiCategoryTreeResponse {
	if categories == nil {
		return nil
	}
	result := make([]dto.PoiCategoryTreeResponse, len(categories))
	for i, pc := range categories {
		result[i].FromDomainTree(&pc)
		result[i].IconPath = m.resolveIconPath(pc.Icon)
	}
	return result
}

func (m *PoiCategoryMapper) BuildTree(categories []domain.PoiCategory) []dto.PoiCategoryTreeResponse {
	categoryMap := make(map[uint64]*domain.PoiCategory)
	for i := range categories {
		categoryMap[categories[i].ID] = &categories[i]
	}

	var roots []dto.PoiCategoryTreeResponse
	processed := make(map[uint64]bool)

	for i := range categories {
		if processed[categories[i].ID] {
			continue
		}
		if categories[i].ParentID == nil || categoryMap[*categories[i].ParentID] == nil {
			root := m.buildTreeNode(&categories[i], categoryMap, processed)
			roots = append(roots, *root)
		}
	}

	return roots
}

func (m *PoiCategoryMapper) buildTreeNode(category *domain.PoiCategory,
	categoryMap map[uint64]*domain.PoiCategory, processed map[uint64]bool) *dto.PoiCategoryTreeResponse {

	processed[category.ID] = true
	node := m.ToTreeResponse(category)

	// Find children
	for i := range categoryMap {
		child := categoryMap[i]
		if child.ParentID != nil && *child.ParentID == category.ID && !processed[child.ID] {
			childNode := m.buildTreeNode(child, categoryMap, processed)
			node.Children = append(node.Children, *childNode)
		}
	}

	return node
}

func (m *PoiCategoryMapper) ToDomainFromCreate(req *dto.CreatePoiCategoryRequest, userID uint64) (*domain.PoiCategory, error) {
	if req == nil {
		return nil, nil
	}

	now := time.Now()
	return &domain.PoiCategory{
		BaseEntity: _entity.BaseEntity{
			CreatedAt: &now,
			UpdatedAt: &now,
			AuditBase: _entity.AuditBase{
				CreatedBy: &userID,
				UpdatedBy: &userID,
			},
		},
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
		ParentID:    normalizeParentID(req.ParentID),
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
		Level:       1, // Sẽ được tính toán lại trong repository
	}, nil
}

// normalizeParentID: parent_id=0 không hợp lệ với FK → lưu NULL (root category).
func normalizeParentID(parentID *uint64) *uint64 {
	if parentID == nil || *parentID == 0 {
		return nil
	}
	return parentID
}

func (m *PoiCategoryMapper) ToDomainFromUpdate(req *dto.UpdatePoiCategoryRequest, existing *domain.PoiCategory, userID uint64) (*domain.PoiCategory, error) {
	if req == nil {
		return nil, nil
	}

	now := time.Now()
	updated := &domain.PoiCategory{
		BaseEntity: _entity.BaseEntity{
			ID:        existing.ID,
			CreatedAt: existing.CreatedAt,
			UpdatedAt: &now,
			AuditBase: _entity.AuditBase{
				CreatedBy: &existing.CreatedBy,
				UpdatedBy: &userID,
			},
		},
		Code:        existing.Code,
		Name:        existing.Name,
		Description: existing.Description,
		Icon:        existing.Icon,
		Color:       existing.Color,
		ParentID:    existing.ParentID,
		Level:       existing.Level,
		Path:        existing.Path,
		SortOrder:   existing.SortOrder,
		IsActive:    existing.IsActive,
		POICount:    existing.POICount,
	}

	// Update only provided fields
	if req.Code != nil {
		updated.Code = *req.Code
	}
	if req.Name != nil {
		updated.Name = *req.Name
	}
	if req.Description != nil {
		updated.Description = *req.Description
	}
	if req.Icon != nil {
		updated.Icon = *req.Icon
	}
	if req.Color != nil {
		updated.Color = *req.Color
	}
	if req.ParentID != nil {
		updated.ParentID = normalizeParentID(req.ParentID)
	}
	if req.SortOrder != nil {
		updated.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		updated.IsActive = *req.IsActive
	}

	return updated, nil
}

func (m *PoiCategoryMapper) ToProto(pc *domain.PoiCategory) *tqdpb.PoiCategory {
	if pc == nil {
		return nil
	}

	return &tqdpb.PoiCategory{
		Id:          pc.ID,
		Name:        pc.Name,
		Description: pc.Description,
		Code:        pc.Code,
		Icon:        pc.Icon,
		IconPath:    m.resolveIconPath(pc.Icon),
		Color:       pc.Color,
		IsActive:    pc.IsActive,
		SortOrder:   uint32(pc.SortOrder),
		ParentId:    m.uint64PtrToValue(pc.ParentID),
		Level:       uint32(pc.Level),
		PoiCount:    uint32(pc.POICount),
		Path:        pc.Path,
		CreatedAt:   _utils.FormatTimeToString(pc.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(pc.UpdatedAt),
	}
}

func (m *PoiCategoryMapper) ToProtoFromResponse(resp *dto.PoiCategoryResponse) *tqdpb.PoiCategory {
	if resp == nil {
		return nil
	}

	iconPath := resp.IconPath
	if iconPath == "" {
		iconPath = m.resolveIconPath(resp.Icon)
	}

	return &tqdpb.PoiCategory{
		Id:          resp.ID,
		Name:        resp.Name,
		Description: resp.Description,
		Code:        resp.Code,
		Icon:        resp.Icon,
		IconPath:    iconPath,
		Color:       resp.Color,
		IsActive:    resp.IsActive,
		SortOrder:   uint32(resp.SortOrder),
		ParentId:    m.uint64PtrToValue(resp.ParentID),
		Level:       uint32(resp.Level),
		PoiCount:    uint32(resp.POICount),
		Path:        resp.Path,
		CreatedAt:   _utils.FormatTimeToString(&resp.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(&resp.UpdatedAt),
	}
}

func (m *PoiCategoryMapper) ToProtoList(categories []domain.PoiCategory) []*tqdpb.PoiCategory {
	if categories == nil {
		return nil
	}
	result := make([]*tqdpb.PoiCategory, len(categories))
	for i, pc := range categories {
		result[i] = m.ToProto(&pc)
	}
	return result
}

func (m *PoiCategoryMapper) ToProtoListFromResponses(responses []dto.PoiCategoryResponse) []*tqdpb.PoiCategory {
	if responses == nil {
		return nil
	}
	result := make([]*tqdpb.PoiCategory, len(responses))
	for i, resp := range responses {
		result[i] = m.ToProtoFromResponse(&resp)
	}
	return result
}

func (m *PoiCategoryMapper) BuildTreeNodes(categories []domain.PoiCategory) []*tqdpb.PoiCategoryTreeNode {
	categoryMap := make(map[uint64]*domain.PoiCategory)
	for i := range categories {
		categoryMap[categories[i].ID] = &categories[i]
	}

	// Build tree nodes
	var rootNodes []*tqdpb.PoiCategoryTreeNode
	processed := make(map[uint64]bool)

	for i := range categories {
		if processed[categories[i].ID] {
			continue
		}
		if categories[i].ParentID == nil || categoryMap[*categories[i].ParentID] == nil {
			node := m.buildProtoTreeNode(&categories[i], categoryMap, processed)
			rootNodes = append(rootNodes, node)
		}
	}

	return rootNodes
}

func (m *PoiCategoryMapper) buildProtoTreeNode(category *domain.PoiCategory,
	categoryMap map[uint64]*domain.PoiCategory, processed map[uint64]bool) *tqdpb.PoiCategoryTreeNode {

	processed[category.ID] = true
	node := &tqdpb.PoiCategoryTreeNode{
		Id:       category.ID,
		Name:     category.Name,
		Icon:     category.Icon,
		IconPath: m.resolveIconPath(category.Icon),
		Path:     category.Path,
		PoiCount: uint32(category.POICount),
	}

	// Find children
	for i := range categoryMap {
		child := categoryMap[i]
		if child.ParentID != nil && *child.ParentID == category.ID && !processed[child.ID] {
			childNode := m.buildProtoTreeNode(child, categoryMap, processed)
			node.Children = append(node.Children, childNode)
		}
	}

	return node
}

func (m *PoiCategoryMapper) FromProto(proto *tqdpb.PoiCategory) (*dto.CreatePoiCategoryRequest, error) {
	if proto == nil {
		return nil, nil
	}

	var parentID *uint64
	if proto.ParentId > 0 {
		id := proto.ParentId
		parentID = &id
	}

	return &dto.CreatePoiCategoryRequest{
		Code:        proto.Code,
		Name:        proto.Name,
		Description: proto.Description,
		Icon:        proto.Icon,
		Color:       proto.Color,
		ParentID:    parentID,
		SortOrder:   int32(proto.SortOrder),
		IsActive:    proto.IsActive,
	}, nil
}

func (m *PoiCategoryMapper) uint64PtrToValue(ptr *uint64) uint64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func (m *PoiCategoryMapper) resolveIconPath(icon string) string {
	return icon
	// if config.AppConfig != nil {
	// 	domainFile = config.AppConfig.ConfigApp.Qhpro.DomainFile
	// }
	// return enums.ResolvePoiCategoryIconPath(domainFile, icon)
}
