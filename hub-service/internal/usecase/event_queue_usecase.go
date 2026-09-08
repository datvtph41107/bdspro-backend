package usecase

import (
	_err "common/domain/err"
	_utils "common/utils"
	"context"
	"hub/internal/domain"
	"hub/internal/dto"
	"hub/internal/enums"
	"hub/internal/repo"
)

type EventQueueUsecase struct {
	eventQueueRepo repo.IEventQueueRepo
}

func NewEventQueueUsecase(eventQueueRepo repo.IEventQueueRepo) *EventQueueUsecase {
	return &EventQueueUsecase{
		eventQueueRepo: eventQueueRepo,
	}
}

// Search tìm kiếm event queue với phân trang
func (uc *EventQueueUsecase) Search(ctx context.Context, searchDTO dto.EventQueueSearchDTO) ([]domain.EventQueueEntity, int64, *_err.ErrorDTO) {
	events, total, err := uc.eventQueueRepo.Search(ctx, searchDTO)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: "Lỗi khi tìm kiếm event queue: " + err.Error(),
		}
	}

	return events, total, nil
}

// Detail lấy chi tiết event queue
func (uc *EventQueueUsecase) Detail(ctx context.Context, id uint64) (*domain.EventQueueEntity, *_err.ErrorDTO) {
	event, err := uc.eventQueueRepo.GetByID(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy event queue",
		}
	}

	return event, nil
}

// Create tạo mới event queue
func (uc *EventQueueUsecase) Create(ctx context.Context, saveDTO dto.EventQueueSaveDTO) (*domain.EventQueueEntity, *_err.ErrorDTO) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	organizationID := _utils.GetOrganizationIdFromContext(ctx)

	entity := &domain.EventQueueEntity{
		EventType:      saveDTO.EventType,
		EventName:      saveDTO.EventName,
		EventData:      saveDTO.EventData,
		Metadata:       saveDTO.Metadata,
		Status:         enums.EEventQueueStatusPending,
		ProfileID:      &profileID,
		OrganizationID: &organizationID,
		RetryCount:     0,
		MaxRetries:     3,
	}

	// Override defaults if provided
	if saveDTO.Status != nil {
		entity.Status = *saveDTO.Status
	}
	if saveDTO.ProfileID != nil {
		entity.ProfileID = saveDTO.ProfileID
	}
	if saveDTO.OrganizationID != nil {
		entity.OrganizationID = saveDTO.OrganizationID
	}
	if saveDTO.MaxRetries != nil {
		entity.MaxRetries = *saveDTO.MaxRetries
	}
	if saveDTO.ScheduledAt != nil {
		scheduledAt := _utils.ParseStringToTime(*saveDTO.ScheduledAt)
		if scheduledAt != nil {
			entity.ScheduledAt = scheduledAt
		}
	}

	err := uc.eventQueueRepo.Create(ctx, entity)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Lỗi khi tạo event queue: " + err.Error(),
		}
	}

	return entity, nil
}

// Update cập nhật event queue
func (uc *EventQueueUsecase) Update(ctx context.Context, id uint64, saveDTO dto.EventQueueSaveDTO) (*domain.EventQueueEntity, *_err.ErrorDTO) {
	// Check if exists
	existingEvent, err := uc.eventQueueRepo.GetByID(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy event queue",
		}
	}

	// Update fields
	existingEvent.EventType = saveDTO.EventType
	existingEvent.EventName = saveDTO.EventName
	existingEvent.EventData = saveDTO.EventData
	existingEvent.Metadata = saveDTO.Metadata

	if saveDTO.Status != nil {
		existingEvent.Status = *saveDTO.Status
	}
	if saveDTO.ProfileID != nil {
		existingEvent.ProfileID = saveDTO.ProfileID
	}
	if saveDTO.OrganizationID != nil {
		existingEvent.OrganizationID = saveDTO.OrganizationID
	}
	if saveDTO.MaxRetries != nil {
		existingEvent.MaxRetries = *saveDTO.MaxRetries
	}
	if saveDTO.RetryCount != nil {
		existingEvent.RetryCount = *saveDTO.RetryCount
	}
	if saveDTO.ErrorMessage != nil {
		existingEvent.ErrorMessage = saveDTO.ErrorMessage
	}
	if saveDTO.ScheduledAt != nil {
		scheduledAt := _utils.ParseStringToTime(*saveDTO.ScheduledAt)
		if scheduledAt != nil {
			existingEvent.ScheduledAt = scheduledAt
		}
	}

	updateErr := uc.eventQueueRepo.Update(ctx, existingEvent.ID, existingEvent)
	if updateErr != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Lỗi khi cập nhật event queue: " + updateErr.Error(),
		}
	}

	return existingEvent, nil
}

// Delete xóa mềm event queue
func (uc *EventQueueUsecase) Delete(ctx context.Context, id uint64) *_err.ErrorDTO {
	// Check if exists
	existingEvent, err := uc.eventQueueRepo.GetByID(ctx, id)
	if err != nil {
		return &_err.ErrorDTO{
			Code:    404,
			Message: "Không tìm thấy event queue",
		}
	}

	deleteErr := uc.eventQueueRepo.Delete(ctx, existingEvent.ID)
	if deleteErr != nil {
		return &_err.ErrorDTO{
			Code:    500,
			Message: "Lỗi khi xóa event queue: " + deleteErr.Error(),
		}
	}

	return nil
}

// UpdateStatus cập nhật trạng thái event queue
func (uc *EventQueueUsecase) UpdateStatus(ctx context.Context, statusDTO dto.EventQueueStatusDTO) (*domain.EventQueueEntity, *_err.ErrorDTO) {
	// Validate status
	if !statusDTO.Status.IsValid() {
		return nil, &_err.ErrorDTO{
			Code:    400,
			Message: "Trạng thái không hợp lệ",
		}
	}

	updatedEvent, err := uc.eventQueueRepo.UpdateStatus(ctx, statusDTO.ID, uint32(statusDTO.Status), statusDTO.ErrorMessage)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Lỗi khi cập nhật trạng thái: " + err.Error(),
		}
	}

	return updatedEvent, nil
}

// GetPendingEvents lấy danh sách event đang chờ xử lý
func (uc *EventQueueUsecase) GetPendingEvents(ctx context.Context, limit int) ([]domain.EventQueueEntity, *_err.ErrorDTO) {
	events, err := uc.eventQueueRepo.GetPendingEvents(ctx, limit)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Lỗi khi lấy danh sách event pending: " + err.Error(),
		}
	}

	return events, nil
}
