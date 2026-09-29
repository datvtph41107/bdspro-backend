package mapper

import (
	_utils "common/utils"
	"strings"
	"user/enums"
	"user/internal/dto"
	models "user/internal/models"

	userpb "pb/types/user"
)

type TagMapper struct{}

func NewTagMapper() *TagMapper {
	return &TagMapper{}
}

func (m *TagMapper) AdminTagRequestToDTO(req *userpb.AdminTagRequest) *dto.AdminTagRequest {
	if req == nil {
		return nil
	}

	return &dto.AdminTagRequest{
		ID:        req.Id,
		Name:      strings.TrimSpace(req.Name),
		TagType:   enums.TagType(req.TagType),
		IsActive:  req.GetIsActive(),
		IsDefault: req.GetIsDefault(),
	}
}

func (m *TagMapper) TagListRequestToDTO(req *userpb.TagListRequest) *dto.TagListRequest {
	if req == nil {
		return nil
	}

	return &dto.TagListRequest{
		TagType: enums.TagType(req.TagType),
		Page:    req.Page,
		Size:    req.Size,
		Text:    req.Text,
		Sort:    req.Sort,
	}
}

func (m *TagMapper) AdminTagListRequestToDTO(req *userpb.AdminTagListRequest) *dto.AdminTagListRequest {
	if req == nil {
		return nil
	}

	tagTypes := make([]enums.TagType, 0, len(req.TagTypes))
	for _, t := range req.TagTypes {
		tagTypes = append(tagTypes, enums.TagType(t))
	}

	var isDefault *bool
	if req.IsDefault != nil {
		isDefault = req.IsDefault
	}

	return &dto.AdminTagListRequest{
		TagTypes:  tagTypes,
		Page:      req.Page,
		Size:      req.Size,
		Text:      req.Text,
		Sort:      req.Sort,
		IsDefault: isDefault,
	}
}

func (m *TagMapper) GetProfileTagsByTypeRequestToDTO(req *userpb.GetProfileTagsByTypeRequest) (uint64, enums.TagType) {
	if req == nil {
		return 0, enums.TagTypeMainArea // Default
	}
	tagType := enums.TagType(req.TagType)
	if tagType == 0 {
		tagType = enums.TagTypeMainArea // Default là mainArea
	}
	return req.ProfileId, tagType
}

func (m *TagMapper) TagEntityToPb(tag *models.TagEntity) *userpb.TagItem {
	if tag == nil {
		return nil
	}

	return &userpb.TagItem{
		Id:        tag.ID,
		Name:      tag.Name,
		TagType:   uint32(tag.TagType),
		IsDefault: tag.IsDefault,
		IsActive:  tag.IsActive,
		CreatedAt: _utils.FormatTimeToString(tag.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(tag.UpdatedAt),
	}
}

func (m *TagMapper) TagsToPb(tags []models.TagEntity) []*userpb.TagItem {
	if len(tags) == 0 {
		return []*userpb.TagItem{}
	}

	result := make([]*userpb.TagItem, 0, len(tags))
	for i := range tags {
		result = append(result, m.TagEntityToPb(&tags[i]))
	}
	return result
}
