package usecase

import (
	_errors "common/errors"
	"context"
	"crm/internal"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/interface/provider"
	"crm/internal/repo"
)

type SharingAccessUsecase struct {
	repo        repo.SharingAccessRepo
	contactRepo repo.ContactRepo
	userClient  provider.UserClient
}

func NewSharingAccessUsecase(
	repo repo.SharingAccessRepo,
	contactRepo repo.ContactRepo,
	userClient provider.UserClient,
) *SharingAccessUsecase {
	return &SharingAccessUsecase{
		repo:        repo,
		contactRepo: contactRepo,
		userClient:  userClient,
	}
}

func (s *SharingAccessUsecase) CheckPermission(c context.Context, contactId uint64) error {

	return nil
}

func (s *SharingAccessUsecase) List(c context.Context, contactId uint64, dto *dto.SharingAccessSearchDTO) ([]domain.SharingEntity, int64, error) {
	if err := s.CheckPermission(c, contactId); err != nil {
		return nil, 0, err
	}
	sharing, total, err := s.repo.List(c, contactId, dto)
	if err != nil {
		return nil, 0, err
	}
	return sharing, total, nil
}

func (s *SharingAccessUsecase) BulkShare(c context.Context,
	contactId uint64,
	permissions []int32,
	note string,
	ownerType enums.EOwnerOf,
	deletedIds []uint64,
	receivers []domain.SharingEntity,
) error {
	contact, err := s.contactRepo.GetByID(c, contactId)
	if err != nil {
		return err
	}
	if contact == nil || contact.ID == 0 {
		return _errors.ReturnError(service.ContactNotFound, _errors.WithPublicMessage("Liên hệ không tồn tại"), _errors.WithLegacyCode(400))
	}

	profileIds := make([]uint64, len(receivers))
	for i := range receivers {
		profileIds[i] = receivers[i].ReceiverID
	}
	validated, err := s.userClient.ValidateProfileIds(c, profileIds)
	if err != nil {
		return err
	}
	if !validated {
		return _errors.ReturnError(service.UserNotFound, _errors.WithPublicMessage("Người dùng không tồn tại vui lòng kiểm tra lại"), _errors.WithLegacyCode(400))
	}

	if err := s.CheckPermission(c, contactId); err != nil {
		return err
	}

	for i := range receivers {
		receivers[i].ContactID = contactId
		receivers[i].Permissions = permissions
		receivers[i].Note = note
	}

	_, err = s.repo.Bulk(c, deletedIds, receivers)
	if err != nil {
		return err
	}
	return nil
}
