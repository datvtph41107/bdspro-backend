package mapper

import (
	"crm/internal/domain"
	"crm/internal/dto"
	"strings"

	crmpb "pb/types/crm"
)

type TagMapper struct{}

func NewTagMapper() *TagMapper {
	return &TagMapper{}
}

func (m *TagMapper) AdminTagRequestToDTO(req *crmpb.AdminTagRequest) *dto.AdminTagRequest {
	if req == nil {
		return nil
	}

	return &dto.AdminTagRequest{
		ID:        req.Id,
		Name:      strings.TrimSpace(req.Name),
		IsActive:  req.GetIsActive(),
		IsDefault: req.GetIsDefault(),
	}
}

func (m *TagMapper) TagListRequestToDTO(req *crmpb.TagListRequest) *dto.TagListRequest {
	if req == nil {
		return nil
	}

	return &dto.TagListRequest{
		Page: req.Page,
		Size: req.Size,
		Text: req.Text,
	}
}

func (m *TagMapper) AdminTagListRequestToDTO(req *crmpb.AdminTagListRequest) *dto.AdminTagListRequest {
	if req == nil {
		return nil
	}

	return &dto.AdminTagListRequest{
		Page: req.Page,
		Size: req.Size,
		Text: req.Text,
	}
}

func (m *TagMapper) TagEntityToPb(tag *domain.TagEntity) *crmpb.TagItem {
	if tag == nil {
		return nil
	}

	result := &crmpb.TagItem{
		Id:        tag.ID,
		Name:      tag.Name,
		IsDefault: tag.IsDefault,
		IsActive:  tag.IsActive,
	}

	if tag.OwnerID != nil {
		result.OwnerId = tag.OwnerID
	}
	if tag.OwnerOf != nil {
		result.OwnerOf = tag.OwnerOf
	}

	return result
}

func (m *TagMapper) TagsToPb(tags []domain.TagEntity) []*crmpb.TagItem {
	if len(tags) == 0 {
		return []*crmpb.TagItem{}
	}

	result := make([]*crmpb.TagItem, 0, len(tags))
	for i := range tags {
		result = append(result, m.TagEntityToPb(&tags[i]))
	}
	return result
}

func (m *TagMapper) PbToTagEntities(pbTags []*crmpb.TagItem) []domain.TagEntity {
	if len(pbTags) == 0 {
		return []domain.TagEntity{}
	}

	result := make([]domain.TagEntity, 0, len(pbTags))
	for _, pbTag := range pbTags {
		if pbTag == nil {
			continue
		}
		tag := domain.TagEntity{}
		tag.ID = pbTag.Id
		tag.Name = pbTag.Name
		tag.IsDefault = pbTag.IsDefault
		tag.IsActive = pbTag.IsActive
		if pbTag.OwnerId != nil {
			tag.OwnerID = pbTag.OwnerId
		}
		if pbTag.OwnerOf != nil {
			ownerOfVal := *pbTag.OwnerOf
			tag.OwnerOf = &ownerOfVal
		}
		result = append(result, tag)
	}
	return result
}