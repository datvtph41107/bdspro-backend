package usecase

import (
	_crud "common/domain/crud"
	"context"
	"fmt"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"
)

// OpenHourUsecase defines the interface for open hour usecase
type OpenHourUsecase interface {
	_crud.IBaseUsecase[domain.OpenHour]
	GetTimesByDayOfWeek(ctx context.Context, poiID *uint64) ([]dto.OpenHourTimesByDayResponse, error)
	BulkSave(ctx context.Context, req *dto.OpenHourBulkSaveRequest) (*dto.OpenHourBulkSaveResponse, error)
}

// openHourUsecase implements OpenHourUsecase interface
type openHourUsecase struct {
	_crud.BaseUsecase[domain.OpenHour, repo.IOpenHourRepo]
}

// NewOpenHourUsecase creates a new OpenHourUsecase
// @bind: tqd/internal/usecase.OpenHourUsecase
func NewOpenHourUsecase(
	openHourRepo repo.IOpenHourRepo,
) OpenHourUsecase {
	return &openHourUsecase{
		BaseUsecase: _crud.BaseUsecase[domain.OpenHour, repo.IOpenHourRepo]{
			Repo: openHourRepo,
		},
	}
}

// GetTimesByDayOfWeek lấy danh sách giờ mở cửa theo từng ngày trong tuần
func (u *openHourUsecase) GetTimesByDayOfWeek(ctx context.Context, poiID *uint64) ([]dto.OpenHourTimesByDayResponse, error) {
	grouped, err := u.Repo.GetGroupedByDayOfWeek(ctx, poiID)
	if err != nil {
		return nil, err
	}

	// Map tên ngày
	dayNames := map[int]string{
		2: "Thứ Hai",
		3: "Thứ Ba",
		4: "Thứ Tư",
		5: "Thứ Năm",
		6: "Thứ Sáu",
		7: "Thứ Bảy",
		8: "Chủ Nhật",
	}

	var result []dto.OpenHourTimesByDayResponse

	// Sắp xếp theo day_of_week từ 2 đến 8
	for day := 2; day <= 8; day++ {
		if hours, exists := grouped[day]; exists {
			var times []dto.OpenHourTimeSlot
			for _, hour := range hours {
				times = append(times, dto.OpenHourTimeSlot{
					ID:        hour.ID,
					Name:      hour.Name,
					POIID:     hour.POIID,
					OpenTime:  hour.OpenTime,
					CloseTime: hour.CloseTime,
					Note:      hour.Note,
					IsOpen:    hour.IsOpen,
					Status:    uint32(hour.Status),
					Type:      uint32(hour.Type),
				})
			}

			result = append(result, dto.OpenHourTimesByDayResponse{
				DayOfWeek: day,
				DayName:   dayNames[day],
				Times:     times,
			})
		}
	}

	return result, nil
}

// BulkSave xử lý nhiều thao tác cùng lúc: create, update, delete
func (u *openHourUsecase) BulkSave(ctx context.Context, req *dto.OpenHourBulkSaveRequest) (*dto.OpenHourBulkSaveResponse, error) {
	response := &dto.OpenHourBulkSaveResponse{
		Total:   len(req.Items) + len(req.DeleteIDs),
		Results: []dto.OpenHourBulkSaveResultItem{},
	}

	// Map tên ngày
	dayNames := map[int]string{
		2: "Thứ 2", 3: "Thứ 3", 4: "Thứ 4", 5: "Thứ 5",
		6: "Thứ 6", 7: "Thứ 7", 8: "Chủ nhật",
	}

	// Xử lý create/update
	for _, item := range req.Items {
		// Map DTO to domain
		openHour := &domain.OpenHour{
			Name:      item.Name,
			POIID:     item.POIID,
			DayOfWeek: item.DayOfWeek,
			OpenTime:  item.OpenTime,
			CloseTime: item.CloseTime,
			Note:      item.Note,
			IsOpen:    item.IsOpen,
		}
		openHour.Status = enums.OpenHourStatus(item.Status)
		openHour.Type = enums.OpenHourType(item.Type)

		// Format thông tin khung giờ
		dayName := dayNames[item.DayOfWeek]
		timeSlotInfo := fmt.Sprintf("Khung giờ(%s-%s, %s)", item.OpenTime, item.CloseTime, dayName)

		result := dto.OpenHourBulkSaveResultItem{
			Item: &item,
		}

		if item.ID != nil && *item.ID > 0 {
			// Update existing record
			err := u.Repo.Update(ctx, *item.ID, openHour)
			if err != nil {
				response.Failed++
				result.Status = "failed"
				result.Message = fmt.Sprintf("%s: %s", timeSlotInfo, err.Error())
			} else {
				response.Updated++
				result.Status = "success"
			}
		} else {
			// Create new record
			err := u.Repo.Create(ctx, openHour)
			if err != nil {
				response.Failed++
				result.Status = "failed"
				result.Message = fmt.Sprintf("%s: %s", timeSlotInfo, err.Error())
			} else {
				response.Created++
				result.Status = "success"
			}
		}

		response.Results = append(response.Results, result)
	}

	// Xử lý delete
	for _, id := range req.DeleteIDs {
		result := dto.OpenHourBulkSaveResultItem{
			Item: &dto.OpenHourBulkItem{ID: &id},
		}

		err := u.Repo.Delete(ctx, id)
		if err != nil {
			response.Failed++
			result.Status = "failed"
			result.Message = fmt.Sprintf("Xóa ID #%d thất bại: %s", id, err.Error())
		} else {
			response.Deleted++
			result.Status = "success"
		}

		response.Results = append(response.Results, result)
	}

	return response, nil
}
