package postgre

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	err := s.db.WithContext(c).
		Model(&domain.SharingEntity{}).
		Table("sharing_access").
		Where("id IN (?)", deleteIds).
		Update("deleted_at", time.Now()).Error
	if err != nil {
		return nil, err
	}
	if len(sharing) == 0 {
		return sharing, nil
	}

	result := s.db.WithContext(c).
		Model(&domain.SharingEntity{}).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "contact_id"},
					{Name: "receiver_id"},
					{Name: "receiver_type"},
				},
				DoUpdates: clause.AssignmentColumns([]string{"permissions"}),
			},
		).Create(&sharing)

	if result.Error != nil {
		return nil, result.Error
	}

	return sharing, nil
}