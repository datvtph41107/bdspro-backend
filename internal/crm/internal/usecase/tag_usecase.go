package usecase

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"crm/internal"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/repo"
	"strings"

	"gorm.io/gorm"
)

type TagUsecase struct {
	tagRepo repo.TagRepo
}

func NewTagUsecase(tagRepo repo.TagRepo) *TagUsecase {
	return &TagUsecase{
		tagRepo: tagRepo,
	}
}

func (u *TagUsecase) AdminCreateTag(ctx context.Context, req *dto.AdminTagRequest) (*domain.TagEntity, error) {
	if err := validateAdminTagRequest(req); err != nil {
		return nil, err
	}

	tag := &domain.TagEntity{
		Name:      strings.TrimSpace(req.Name),
		IsDefault: true,
		IsActive:  req.IsActive,
		OwnerID:   nil,
		OwnerOf:   nil,
	}

	if err := u.tagRepo.Create(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (u *TagUsecase) AdminUpdateTag(ctx context.Context, req *dto.AdminTagRequest) (*domain.TagEntity, error) {
	if req == nil || req.ID == 0 {
		return nil, _errors.ReturnError(service.TagPayloadRequired)
	}
	if err := validateAdminTagRequest(req); err != nil {
		return nil, err
	}

	existed, err := u.tagRepo.GetByID(ctx, req.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, _errors.ReturnError(service.TagNotFound)
		}
		return nil, err
	}

	existed.Name = strings.TrimSpace(req.Name)
	existed.IsDefault = req.IsDefault
	existed.IsActive = req.IsActive

	if err := u.tagRepo.Update(ctx, existed.ID, existed); err != nil {
		return nil, err
	}
	return existed, nil
}

func (u *TagUsecase) AdminDeleteTag(ctx context.Context, id uint64) error {
	if id == 0 {
		return _errors.ReturnError(service.TagPayloadRequired)
	}
	return u.tagRepo.Delete(ctx, id)
}

func (u *TagUsecase) ListTagsForContact(ctx context.Context, req *dto.TagListRequest) ([]domain.TagEntity, uint32, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	ownerOf := _utils.GetOrganizationIdFromContext(ctx)

	var ownerOfInt32 int32
	if ownerOf != 0 {
		ownerOfInt32 = 30 // Organization
	} else {
		ownerOfInt32 = 10 // Member
	}

	pagable := &_dto.Pagable{
		Page: req.Page,
		Size: req.Size,
	}
	return u.tagRepo.ListForContact(ctx, profileID, ownerOfInt32, pagable, strings.TrimSpace(req.Text))
}

func (u *TagUsecase) AdminListTags(ctx context.Context, req *dto.AdminTagListRequest) ([]domain.TagEntity, uint32, error) {
	pagable := &_dto.Pagable{
		Page: req.Page,
		Size: req.Size,
	}
	return u.tagRepo.ListAdmin(ctx, pagable, strings.TrimSpace(req.Text))
}

func (u *TagUsecase) ListTagsByContactID(ctx context.Context, contactID uint64) ([]domain.TagEntity, error) {
	return u.tagRepo.ListByContactID(ctx, contactID)
}

func (u *TagUsecase) ReplaceContactTags(ctx context.Context, contactID uint64, tags []domain.TagEntity) error {
	profileID := _utils.GetProfileIdWithContext(ctx)
	organizationID := _utils.GetOrganizationIdFromContext(ctx)

	var ownerID uint64
	var ownerOf *int32

	if organizationID != 0 {
		ownerID = organizationID
		ownerOfVal := int32(30) // Organization
		ownerOf = &ownerOfVal
	} else {
		ownerID = profileID
		ownerOfVal := int32(10) // Member
		ownerOf = &ownerOfVal
	}

	// Set owner cho các tag mới
	for i := range tags {
		if tags[i].ID == 0 {
			tags[i].OwnerID = &ownerID
			tags[i].OwnerOf = ownerOf
		}
	}

	return u.tagRepo.ReplaceContactTags(ctx, contactID, tags)
}

func validateAdminTagRequest(req *dto.AdminTagRequest) error {
	if req == nil {
		return _errors.ReturnError(service.TagPayloadRequired)
	}
	if strings.TrimSpace(req.Name) == "" {
		return _errors.ReturnError(service.TagNameRequired)
	}
	return nil
}
