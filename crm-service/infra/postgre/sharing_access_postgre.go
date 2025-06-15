package postgre

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"

	"gorm.io/gorm"
)

// @bind: crm/internal/repo.SharingAccessRepo
type SharingAccessPostgre struct {
	db *gorm.DB
}

func NewSharingAccessPostgre(db *gorm.DB) *SharingAccessPostgre {
	return &SharingAccessPostgre{db: db}
}

func (s *SharingAccessPostgre) List(c context.Context, contactId uint64, dto *dto.SharingAccessSearchDTO) ([]domain.SharingEntity, int64, error) {
	sharing := []domain.SharingEntity{}
	err := s.db.WithContext(c).Where("contact_id = ?", contactId).Find(&sharing).Error
	return sharing, int64(len(sharing)), err
}

func (s *SharingAccessPostgre) Bulk(c context.Context, deleteIds []uint64, sharing []domain.SharingEntity) ([]domain.SharingEntity, error) {
	err := s.db.WithContext(c).Where("id IN (?)", deleteIds).Delete(&domain.SharingEntity{}).Error
	if err != nil {
		return nil, err
	}
	return sharing, nil
}
