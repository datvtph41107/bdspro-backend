package usecase

import (
	"context"
	"fmt"
	"time"

	_utils "common/utils"
	"organization/internal/dto"
	"organization/internal/enums"
	iusecase "organization/internal/interface"
)

type EventHistoryUsecase struct {
	notiClient iusecase.INotiClient
}

func NewEventHistoryUsecase(notiClient iusecase.INotiClient) *EventHistoryUsecase {
	return &EventHistoryUsecase{
		notiClient: notiClient,
	}
}

// logDealHistory ghi lại lịch sử thương vụ cho dealUsecase
func (u *EventHistoryUsecase) LogDealHistory(ctx context.Context, dealID uint64, eventType enums.DealHistoryEventType, content string, metadata map[string]interface{}) {
	// Lấy thông tin user hiện tại
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Tạo deal history DTO
	history := &dto.DealHistoryDTO{
		DealID:      dealID,
		ActorID:     currentUserId,
		ActorName:   "", // Sẽ được lấy từ user service nếu cần
		ActorAvatar: "", // Sẽ được lấy từ user service nếu cần
		ActionType:  eventType,
		ActionName:  eventType.GetEventName(),
		Content:     content,
		CreatedAt:   time.Now(),
		Metadata:    metadata,
	}

	// Gọi deal history client để lưu lịch sử
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		if err := u.notiClient.HistoryDeal(cloneCtx, history); err != nil {
			// Log error nhưng không ảnh hưởng đến flow chính
			fmt.Printf("Failed to log deal history: %v\n", err)
		}
	}()
}

// // logDealHistory ghi lại lịch sử thương vụ cho dealInvitationUsecase
// func (u *dealInvitationUsecase) logDealHistory(ctx context.Context, dealID uint64, eventType enums.DealHistoryEventType, content string, metadata map[string]interface{}) {
// 	// Lấy thông tin user hiện tại
// 	currentUserId := _utils.GetProfileIdWithContext(ctx)

// 	// Tạo deal history DTO
// 	history := &dto.DealHistoryDTO{
// 		DealID:      dealID,
// 		ActorID:     currentUserId,
// 		ActorName:   "", // Sẽ được lấy từ user service nếu cần
// 		ActorAvatar: "", // Sẽ được lấy từ user service nếu cần
// 		ActionType:  string(eventType),
// 		ActionName:  eventType.GetEventName(),
// 		Content:     content,
// 		CreatedAt:   time.Now(),
// 		Metadata:    metadata,
// 	}

// 	// Gọi notification client để lưu lịch sử
// 	go func() {
// 		cloneCtx := _utils.CloneContext(ctx)
// 		if err := u.notiClient.HistoryDeal(cloneCtx, history); err != nil {
// 			// Log error nhưng không ảnh hưởng đến flow chính
// 			fmt.Printf("Failed to log deal history: %v\n", err)
// 		}
// 	}()
// }
