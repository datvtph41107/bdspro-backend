package usecase

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/repo"
)

type SharingAccessUsecase struct {
	repo repo.SharingAccessRepo
}

func NewSharingAccessUsecase(repo repo.SharingAccessRepo) *SharingAccessUsecase {
	return &SharingAccessUsecase{repo: repo}
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

func (s *SharingAccessUsecase) BulkShare(c context.Context, contactId uint64, deletedIds []uint64, receivers []domain.SharingEntity) error {
	if err := s.CheckPermission(c, contactId); err != nil {
		return err
	}
	_, err := s.repo.Bulk(c, deletedIds, receivers)
	if err != nil {
		return err
	}
	return nil
}
