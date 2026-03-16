package mapper

import (
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"tqd/internal/domain"
)

// PoiCategoryMapper handles mapping between domain and proto
type PoiCategoryMapper struct{}

// NewPoiCategoryMapper creates a new PoiCategoryMapper
func NewPoiCategoryMapper() *PoiCategoryMapper {
	return &PoiCategoryMapper{}
}

// ToProto converts domain PoiCategory to proto PoiCategory
func (m *PoiCategoryMapper) ToProto(poiCategory *domain.PoiCategory) *tqdpb.PoiCategory {
	if poiCategory == nil {
		return nil
	}

	return &tqdpb.PoiCategory{
		Id:          poiCategory.ID,
		Name:        poiCategory.Name,
		Description: poiCategory.Description,
		Code:        poiCategory.Code,
		Icon:        poiCategory.Icon,
		IsActive:    poiCategory.IsActive,
		CreatedAt:   _utils.FormatTimeToString(poiCategory.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(poiCategory.UpdatedAt),
	}
}

// ToProtoList converts slice of domain PoiCategory to slice of proto PoiCategory
func (m *PoiCategoryMapper) ToProtoList(poiCategories []*domain.PoiCategory) []*tqdpb.PoiCategory {
	result := make([]*tqdpb.PoiCategory, len(poiCategories))
	for i, poiCategory := range poiCategories {
		result[i] = m.ToProto(poiCategory)
	}
	return result
}

// BuildTreeNodes builds a tree structure from flat domain PoiCategory list
func (m *PoiCategoryMapper) BuildTreeNodes(categories []*domain.PoiCategory) []*tqdpb.PoiCategoryTreeNode {
	// Map categories by ID
	categoryMap := make(map[uint64]*domain.PoiCategory)
	for _, cat := range categories {
		categoryMap[cat.ID] = cat
	}

	// Build tree nodes
	var rootNodes []*tqdpb.PoiCategoryTreeNode
	processedIDs := make(map[uint64]bool)

	for _, cat := range categories {
		if processedIDs[cat.ID] {
			continue
		}
		if cat.ParentID == nil || categoryMap[*cat.ParentID] == nil {
			// This is a root node
			node := m.buildTreeNode(cat, categoryMap, processedIDs)
			rootNodes = append(rootNodes, node)
		}
	}

	return rootNodes
}

// buildTreeNode recursively builds tree nodes
func (m *PoiCategoryMapper) buildTreeNode(category *domain.PoiCategory, categoryMap map[uint64]*domain.PoiCategory, processedIDs map[uint64]bool) *tqdpb.PoiCategoryTreeNode {
	processedIDs[category.ID] = true

	node := &tqdpb.PoiCategoryTreeNode{
		Data: m.ToProto(category),
	}

	// Find and add children
	for _, potentialChild := range categoryMap {
		if potentialChild.ParentID != nil && *potentialChild.ParentID == category.ID && !processedIDs[potentialChild.ID] {
			childNode := m.buildTreeNode(potentialChild, categoryMap, processedIDs)
			node.Children = append(node.Children, childNode)
		}
	}

	return node
}
