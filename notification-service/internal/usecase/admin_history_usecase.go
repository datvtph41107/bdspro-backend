package usecase

import (
	"context"
	shared_enum "pb/enums"
	"time"

	"notification/internal/domain"
	"notification/internal/dto"
	"notification/internal/repo"
)

type AdminHistoryUsecase struct {
	adminHistoryRepo repo.AdminHistoryRepo
}

func NewAdminHistoryUsecase(adminHistoryRepo repo.AdminHistoryRepo) *AdminHistoryUsecase {
	return &AdminHistoryUsecase{
		adminHistoryRepo: adminHistoryRepo,
	}
}

// Search: Tìm kiếm lịch sử admin theo adminID
func (h *AdminHistoryUsecase) Search(c context.Context, adminID uint64, dto dto.AdminHistorySearchDTO) ([]domain.AdminHistoryEntity, int64, error) {
	histories, total, err := h.adminHistoryRepo.Search(c, adminID, dto)
	if err != nil {
		return nil, 0, err
	}
	return histories, total, nil
}

// SearchByOwner: Tìm kiếm lịch sử admin theo owner
func (h *AdminHistoryUsecase) SearchByOwner(c context.Context, ownerID uint64, ownerType shared_enum.EOwnerType, dto dto.AdminHistorySearchDTO) ([]domain.AdminHistoryEntity, int64, error) {
	histories, total, err := h.adminHistoryRepo.SearchByOwner(c, ownerID, ownerType, dto)
	if err != nil {
		return nil, 0, err
	}
	return histories, total, nil
}

// CreateAdminHistory: Tạo lịch sử admin mới
func (h *AdminHistoryUsecase) CreateAdminHistory(c context.Context, dto *dto.AdminHistoryCreateDTO) (*domain.AdminHistoryEntity, error) {
	entity := &domain.AdminHistoryEntity{
		TargetId:   dto.TargetID,
		TargetType: dto.TargetType,
		ActionType: dto.ActionType,
		Title:      dto.Title,
		Note:       dto.Note,
		PreStage:   dto.PreStage,
		AfterStage: dto.AfterStage,
		AdminID:    dto.AdminID,
		OwnerID:    dto.OwnerID,
		OwnerType:  dto.OwnerType,
		IsInternal: dto.IsInternal,
		AdminRole:  dto.AdminRole,
		IPAddress:  dto.IPAddress,
		UserAgent:  dto.UserAgent,
	}

	err := h.adminHistoryRepo.CreateAdminHistory(c, entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// CreateAdminHistoryBatch: Tạo hàng loạt lịch sử admin
func (h *AdminHistoryUsecase) CreateAdminHistoryBatch(c context.Context, dtos []*dto.AdminHistoryCreateDTO) ([]uint64, error) {
	entities := make([]*domain.AdminHistoryEntity, len(dtos))

	for i, dto := range dtos {
		entity := &domain.AdminHistoryEntity{
			TargetId:   dto.TargetID,
			TargetType: dto.TargetType,
			ActionType: dto.ActionType,
			Title:      dto.Title,
			Note:       dto.Note,
			PreStage:   dto.PreStage,
			AfterStage: dto.AfterStage,
			AdminID:    dto.AdminID,
			OwnerID:    dto.OwnerID,
			OwnerType:  dto.OwnerType,
			IsInternal: dto.IsInternal,
			AdminRole:  dto.AdminRole,
			IPAddress:  dto.IPAddress,
			UserAgent:  dto.UserAgent,
		}
		entities[i] = entity
	}

	// Tạo batch trong database
	err := h.adminHistoryRepo.CreateAdminHistoryBatch(c, entities)
	if err != nil {
		return nil, err
	}

	// Lấy IDs của các entity đã tạo
	ids := make([]uint64, len(entities))
	for i, entity := range entities {
		ids[i] = entity.ID
	}

	return ids, nil
}

// GetAdminActions: Lấy danh sách hành động của admin trong khoảng thời gian
func (h *AdminHistoryUsecase) GetAdminActions(c context.Context, adminID uint64, fromDate time.Time, toDate time.Time) ([]domain.AdminHistoryEntity, error) {
	actions, err := h.adminHistoryRepo.GetAdminActions(c, adminID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	return actions, nil
}

// GetRecentAdminActions: Lấy danh sách hành động admin gần đây
func (h *AdminHistoryUsecase) GetRecentAdminActions(c context.Context, limit int) ([]domain.AdminHistoryEntity, error) {
	actions, err := h.adminHistoryRepo.GetRecentAdminActions(c, limit)
	if err != nil {
		return nil, err
	}
	return actions, nil
}

// GetAdminActionStats: Lấy thống kê hành động của admin
func (h *AdminHistoryUsecase) GetAdminActionStats(c context.Context, adminID uint64, fromDate time.Time, toDate time.Time) (map[string]int64, error) {
	stats, err := h.adminHistoryRepo.GetAdminActionStats(c, adminID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

// SearchInternal: Tìm kiếm nội bộ (cho admin)
func (h *AdminHistoryUsecase) SearchInternal(c context.Context, dto dto.AdminHistorySearchDTO) ([]domain.AdminHistoryEntity, int64, error) {
	histories, total, err := h.adminHistoryRepo.SearchInternal(c, dto)
	if err != nil {
		return nil, 0, err
	}

	for i := range histories {
		status := "success"
		histories[i].Status = &status

		adminName := "Admin " + string(rune(histories[i].AdminID))
		histories[i].AdminName = &adminName
	}

	return histories, total, nil
}
