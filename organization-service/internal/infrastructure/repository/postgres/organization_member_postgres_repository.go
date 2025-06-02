package postgres

import (
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type OrganizationMemberModel struct {
	ID             uint32    `gorm:"primaryKey;autoIncrement"`
	OrganizationID uint32    `gorm:"not null"`
	UserID         uint32    `gorm:"not null"`
	Role           string    `gorm:"not null"`
	Status         string    `gorm:"default:active"`
	JoinedAt       time.Time `gorm:"not null"`
	RemovedAt      *time.Time

	Organization *OrganizationModel
}

func (OrganizationMemberModel) TableName() string {
	return "organization_members"
}

func OrganizationMemberModelToEntity(model OrganizationMemberModel) entity.OrganizationMember {
	return entity.OrganizationMember{
		ID:             model.ID,
		OrganizationID: model.OrganizationID,
		UserID:         model.UserID,
		Role:           entity.OrganizationMemberRole(model.Role),
		Status:         entity.OrganizationMemberStatus(model.Status),
		JoinedAt:       model.JoinedAt,
		RemovedAt:      model.RemovedAt,
	}
}

func OrganizationMemberEntityToModel(entity entity.OrganizationMember) OrganizationMemberModel {
	return OrganizationMemberModel{
		ID:             entity.ID,
		OrganizationID: entity.OrganizationID,
		UserID:         entity.UserID,
		Role:           string(entity.Role),
		Status:         string(entity.Status),
		JoinedAt:       entity.JoinedAt,
		RemovedAt:      entity.RemovedAt,
	}
}
