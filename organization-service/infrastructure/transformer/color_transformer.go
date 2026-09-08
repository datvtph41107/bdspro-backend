package transformer

import (
	_models "common/models"
	"organization/internal/domain/entity"
	organizationpb "pb/types/organization"
	"time"
)

type ColorTransformer struct{}

func NewColorTransformer() *ColorTransformer {
	return &ColorTransformer{}
}

func (t *ColorTransformer) ToProto(color *entity.Color) *organizationpb.Color {
	if color == nil {
		return nil
	}

	protoColor := &organizationpb.Color{
		Id:              uint32(color.ID),
		Name:            color.Name,
		ContentColor:    color.ContentColor,
		BackgroundColor: color.BackgroundColor,
		Description:     color.Description,
		IsActive:        color.IsActive,
	}

	if color.CreatedAt != nil {
		protoColor.CreatedAt = color.CreatedAt.Format(time.RFC3339)
	}
	if color.UpdatedAt != nil {
		protoColor.UpdatedAt = color.UpdatedAt.Format(time.RFC3339)
	}
	if color.CreatedBy != nil {
		protoColor.CreatedBy = uint32(*color.CreatedBy)
	}
	if color.UpdatedBy != nil {
		protoColor.UpdatedBy = uint32(*color.UpdatedBy)
	}

	return protoColor
}

func (t *ColorTransformer) ToEntity(protoColor *organizationpb.Color) *entity.Color {
	if protoColor == nil {
		return nil
	}

	color := &entity.Color{
		Name:            protoColor.Name,
		ContentColor:    protoColor.ContentColor,
		BackgroundColor: protoColor.BackgroundColor,
		Description:     protoColor.Description,
		IsActive:        protoColor.IsActive,
	}

	if protoColor.Id > 0 {
		color.ID = uint64(protoColor.Id)
	}

	return color
}

func (t *ColorTransformer) ToProtoList(colors []*entity.Color) []*organizationpb.Color {
	if colors == nil {
		return nil
	}

	result := make([]*organizationpb.Color, len(colors))
	for i, color := range colors {
		result[i] = t.ToProto(color)
	}
	return result
}

func (t *ColorTransformer) ToCreateRequest(req *organizationpb.CreateColorRequest) *entity.Color {
	if req == nil {
		return nil
	}

	return &entity.Color{
		Name:            req.Name,
		ContentColor:    req.ContentColor,
		BackgroundColor: req.BackgroundColor,
		Description:     req.Description,
		IsActive:        req.IsActive,
	}
}

func (t *ColorTransformer) ToUpdateRequest(req *organizationpb.UpdateColorRequest) *entity.Color {
	if req == nil {
		return nil
	}

	return &entity.Color{
		BaseEntity: _models.BaseEntity{
			ID: uint64(req.Id),
		},
		Name:            req.Name,
		ContentColor:    req.ContentColor,
		BackgroundColor: req.BackgroundColor,
		Description:     req.Description,
		IsActive:        req.IsActive,
	}
} 