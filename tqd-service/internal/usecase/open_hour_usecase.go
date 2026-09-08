package usecase

import (
	"context"
	"fmt"
	"time"

	_utils "common/utils"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
)

// OpenHourUsecase defines business logic for open hour
type OpenHourUsecase interface {
	// CRUD operations
	Create(ctx context.Context, req *dto.CreateOpenHourRequest) (*dto.OpenHourResponse, error)
	GetByID(ctx context.Context, id uint64) (*dto.OpenHourResponse, error)
	Update(ctx context.Context, id uint64, req *dto.UpdateOpenHourRequest) (*dto.OpenHourResponse, error)
	Delete(ctx context.Context, id uint64) error

	// List operations
	List(ctx context.Context, filter *dto.OpenHourFilter) (*dto.ListOpenHoursResponse, error)

	// Business operations
	GetTimesByDay(ctx context.Context, req *dto.GetTimesByDayRequest) ([]dto.OpenHourTimesByDay, error)
	BulkSave(ctx context.Context, req *dto.BulkSaveOpenHourRequest) (*dto.BulkSaveOpenHourResponse, error)

	// Validation
	ValidateTimeSlot(ctx context.Context, poiID uint64, dayOfWeek int, startTime, endTime string, excludeID *uint64) error
}

// openHourUsecase implements OpenHourUsecase
type openHourUsecase struct {
	repo   repo.IOpenHourRepo
	mapper *mapper.OpenHourMapper
}

// NewOpenHourUsecase creates new open hour usecase
func NewOpenHourUsecase(
	repo repo.IOpenHourRepo,
	mapper *mapper.OpenHourMapper,
) OpenHourUsecase {
	return &openHourUsecase{
		repo:   repo,
		mapper: mapper,
	}
}

// Create implements OpenHourUsecase.Create
func (u *openHourUsecase) Create(ctx context.Context, req *dto.CreateOpenHourRequest) (*dto.OpenHourResponse, error) {
	// Get user ID from context
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized")
	}

	// Validate time format
	if err := u.validateTimeFormat(req.OpenTime, req.CloseTime); err != nil {
		return nil, err
	}

	// Check time overlap
	if err := u.ValidateTimeSlot(ctx, req.POIID, req.DayOfWeek, req.OpenTime, req.CloseTime, nil); err != nil {
		return nil, err
	}

	// Convert to domain
	openHour, err := u.mapper.ToDomainFromCreate(req, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to map request: %w", err)
	}

	// Save to database
	if err := u.repo.Create(ctx, openHour); err != nil {
		return nil, fmt.Errorf("failed to create open hour: %w", err)
	}

	// Get created entity
	created, err := u.repo.GetByID(ctx, openHour.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created open hour: %w", err)
	}

	return u.mapper.ToResponse(created), nil
}

// GetByID implements OpenHourUsecase.GetByID
func (u *openHourUsecase) GetByID(ctx context.Context, id uint64) (*dto.OpenHourResponse, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid open hour ID")
	}

	openHour, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get open hour: %w", err)
	}
	if openHour == nil {
		return nil, fmt.Errorf("open hour not found")
	}

	return u.mapper.ToResponse(openHour), nil
}

// Update implements OpenHourUsecase.Update
func (u *openHourUsecase) Update(ctx context.Context, id uint64, req *dto.UpdateOpenHourRequest) (*dto.OpenHourResponse, error) {
	// Get user ID
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized")
	}

	// Get existing
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get open hour: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("open hour not found")
	}

	// Validate time format if updating times
	if req.OpenTime != nil || req.CloseTime != nil {
		openTime := existing.OpenTime
		if req.OpenTime != nil {
			openTime = *req.OpenTime
		}
		closeTime := existing.CloseTime
		if req.CloseTime != nil {
			closeTime = *req.CloseTime
		}
		if err := u.validateTimeFormat(openTime, closeTime); err != nil {
			return nil, err
		}
	}

	// Check time overlap if times or day changed
	dayOfWeek := existing.DayOfWeek
	if req.DayOfWeek != nil {
		dayOfWeek = *req.DayOfWeek
	}
	openTime := existing.OpenTime
	if req.OpenTime != nil {
		openTime = *req.OpenTime
	}
	closeTime := existing.CloseTime
	if req.CloseTime != nil {
		closeTime = *req.CloseTime
	}

	if err := u.ValidateTimeSlot(ctx, *existing.POIID, dayOfWeek, openTime, closeTime, &id); err != nil {
		return nil, err
	}

	// Convert to domain
	updated, err := u.mapper.ToDomainFromUpdate(req, existing, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to map request: %w", err)
	}

	// Update in database
	if err := u.repo.Update(ctx, id, updated); err != nil {
		return nil, fmt.Errorf("failed to update open hour: %w", err)
	}

	// Get updated entity
	result, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated open hour: %w", err)
	}

	return u.mapper.ToResponse(result), nil
}

// Delete implements OpenHourUsecase.Delete
func (u *openHourUsecase) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return fmt.Errorf("invalid open hour ID")
	}

	// Check existence
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get open hour: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("open hour not found")
	}

	// Delete
	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete open hour: %w", err)
	}

	return nil
}

// List implements OpenHourUsecase.List
func (u *openHourUsecase) List(ctx context.Context, filter *dto.OpenHourFilter) (*dto.ListOpenHoursResponse, error) {
	// Validate pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Size < 1 {
		filter.Size = 10
	}
	if filter.Size > 100 {
		filter.Size = 100
	}

	// Get from repository
	openHours, total, err := u.repo.ListWithFilter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list open hours: %w", err)
	}

	return &dto.ListOpenHoursResponse{
		Data:  u.mapper.ToResponseList(openHours),
		Total: total,
		Page:  filter.Page,
		Size:  filter.Size,
	}, nil
}

// GetTimesByDay implements OpenHourUsecase.GetTimesByDay
func (u *openHourUsecase) GetTimesByDay(ctx context.Context, req *dto.GetTimesByDayRequest) ([]dto.OpenHourTimesByDay, error) {
	dayOfWeek, err := dto.StringToDayOfWeek(req.DayOfWeek)
	if err != nil {
		return nil, fmt.Errorf("invalid day of week: %w", err)
	}

	var poiID *uint64
	if req.POIID > 0 {
		poiID = &req.POIID
	}

	// Get from repository
	openHours, err := u.repo.FindByDayOfWeek(ctx, dayOfWeek, poiID)
	if err != nil {
		return nil, fmt.Errorf("failed to get open hours: %w", err)
	}

	// Group by day (though we already filtered by day)
	dayMap := make(map[int][]dto.OpenHourTimeSlot)
	for _, oh := range openHours {
		dayMap[oh.DayOfWeek] = append(dayMap[oh.DayOfWeek], dto.OpenHourTimeSlot{
			ID:        oh.ID,
			POIID:     *oh.POIID,
			OpenTime:  oh.OpenTime,
			CloseTime: oh.CloseTime,
			Note:      oh.Note,
			IsOpen:    oh.IsOpen,
		})
	}

	// Build response
	var result []dto.OpenHourTimesByDay
	for day, times := range dayMap {
		result = append(result, dto.OpenHourTimesByDay{
			DayOfWeek: day,
			DayName:   dto.DayOfWeekToString(day),
			Times:     times,
		})
	}

	return result, nil
}

// BulkSave implements OpenHourUsecase.BulkSave
func (u *openHourUsecase) BulkSave(ctx context.Context, req *dto.BulkSaveOpenHourRequest) (*dto.BulkSaveOpenHourResponse, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized")
	}

	response := &dto.BulkSaveOpenHourResponse{
		Results: make([]dto.BulkSaveOpenHourResult, 0, len(req.Items)+len(req.DeleteIDs)),
	}

	// Process creates/updates
	for _, item := range req.Items {
		result := dto.BulkSaveOpenHourResult{ID: item.ID}

		// Validate time format
		if err := u.validateTimeFormat(item.OpenTime, item.CloseTime); err != nil {
			result.Status = "failed"
			result.Message = err.Error()
			response.Failed++
			response.Results = append(response.Results, result)
			continue
		}

		// Convert to domain
		openHour, err := u.mapper.ToDomainFromBulkItem(&item, userID)
		if err != nil {
			result.Status = "failed"
			result.Message = err.Error()
			response.Failed++
			response.Results = append(response.Results, result)
			continue
		}

		// Check time overlap
		var excludeID *uint64
		if item.ID != nil {
			excludeID = item.ID
		}
		if err := u.ValidateTimeSlot(ctx, item.POIID, item.DayOfWeek,
			item.OpenTime, item.CloseTime, excludeID); err != nil {
			result.Status = "failed"
			result.Message = err.Error()
			response.Failed++
			response.Results = append(response.Results, result)
			continue
		}

		// Save
		if item.ID != nil {
			// Update
			if err := u.repo.Update(ctx, *item.ID, openHour); err != nil {
				result.Status = "failed"
				result.Message = err.Error()
				response.Failed++
			} else {
				result.Status = "success"
				response.Updated++
			}
		} else {
			// Create
			if err := u.repo.Create(ctx, openHour); err != nil {
				result.Status = "failed"
				result.Message = err.Error()
				response.Failed++
			} else {
				result.Status = "success"
				response.Created++
			}
		}
		response.Results = append(response.Results, result)
	}

	// Process deletes
	if len(req.DeleteIDs) > 0 {
		if err := u.repo.BulkDelete(ctx, req.DeleteIDs); err != nil {
			// Add failed results for each delete
			for _, id := range req.DeleteIDs {
				response.Results = append(response.Results, dto.BulkSaveOpenHourResult{
					ID:      &id,
					Status:  "failed",
					Message: err.Error(),
				})
				response.Failed++
			}
		} else {
			response.Deleted = len(req.DeleteIDs)
			for _, id := range req.DeleteIDs {
				response.Results = append(response.Results, dto.BulkSaveOpenHourResult{
					ID:     &id,
					Status: "success",
				})
			}
		}
	}

	response.Total = len(req.Items) + len(req.DeleteIDs)
	return response, nil
}

// ValidateTimeSlot implements OpenHourUsecase.ValidateTimeSlot
func (u *openHourUsecase) ValidateTimeSlot(ctx context.Context, poiID uint64, dayOfWeek int,
	startTime, endTime string, excludeID *uint64) error {

	exists, err := u.repo.CheckTimeOverlap(ctx, poiID, dayOfWeek, startTime, endTime, excludeID)
	if err != nil {
		return fmt.Errorf("failed to check time overlap: %w", err)
	}
	if exists {
		return fmt.Errorf("time slot overlaps with existing open hour")
	}
	return nil
}

// validateTimeFormat validates time strings
func (u *openHourUsecase) validateTimeFormat(openTime, closeTime string) error {
	if _, err := time.Parse("15:04:05", openTime); err != nil {
		return fmt.Errorf("invalid open time format, must be HH:MM:SS")
	}
	if _, err := time.Parse("15:04:05", closeTime); err != nil {
		return fmt.Errorf("invalid close time format, must be HH:MM:SS")
	}
	return nil
}
