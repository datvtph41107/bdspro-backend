package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	_utils "common/utils"
	"context"
)

type DealInternalNoteUsecase interface {
	AddInternalNote(ctx context.Context, dealID uint64, content string, actionType enums.DealActionType) (*domain.DealInternalNote, error)
	GetDealHistory(ctx context.Context, dealID uint64, actionTypes []enums.DealActionType, page, size uint32) ([]*domain.DealInternalNote, uint32, error)
}

type dealInternalNoteUsecase struct {
	internalNoteRepository repo.DealInternalNoteRepo
	dealRepository         repo.DealRepository
	// groupAuthUsecase       GroupAuthUsecase
	// logWorker              *LogWorker
}

func NewInternalNoteUsecase(
	internalNoteRepository repo.DealInternalNoteRepo,
	dealRepository repo.DealRepository,
	// groupAuthUsecase GroupAuthUsecase,
	// logWorker *LogWorker,
) DealInternalNoteUsecase {
	return &dealInternalNoteUsecase{
		internalNoteRepository: internalNoteRepository,
		dealRepository:         dealRepository,
		// groupAuthUsecase:       groupAuthUsecase,
		// logWorker:              logWorker,
	}
}

func (u *dealInternalNoteUsecase) AddInternalNote(ctx context.Context, dealID uint64, content string, actionType enums.DealActionType) (*domain.DealInternalNote, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra deal có tồn tại không
	deal, err := u.dealRepository.GetByID(ctx, dealID)
	if err != nil {
		return nil, err
	}
	if deal == nil {
		return nil, domain.ErrNotFound
	}

	// Kiểm tra quyền - chỉ leader mới được thêm ghi chú nội bộ
	// todo:
	// err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	// if err != nil {
	// 	return nil, custom_error.Forbidden("you are not allowed to add internal note")
	// }

	// Tạo ghi chú nội bộ
	note := &domain.DealInternalNote{
		DealID:     dealID,
		ActorID:    currentUserId,
		Content:    content,
		ActionType: actionType,
	}

	createdNote, err := u.internalNoteRepository.Create(ctx, note)
	if err != nil {
		return nil, err
	}

	// Ghi log
	// u.logWorker.Push(LogEvent{
	// 	GroupId: uint32(deal.OwnerId),
	// 	ActorId: uint32(currentUserId),
	// 	LogType: "ADD_INTERNAL_NOTE",
	// 	LogData: fmt.Sprintf("Added internal note for deal %d", dealID),
	// })

	return createdNote, nil
}

func (u *dealInternalNoteUsecase) GetDealHistory(ctx context.Context, dealID uint64, actionTypes []enums.DealActionType, page, size uint32) ([]*domain.DealInternalNote, uint32, error) {
	// Kiểm tra deal có tồn tại không
	deal, err := u.dealRepository.GetByID(ctx, dealID)
	if err != nil {
		return nil, 0, err
	}
	if deal == nil {
		return nil, 0, domain.ErrNotFound
	}

	// Lấy lịch sử
	notes, total, err := u.internalNoteRepository.GetByDealID(ctx, dealID, actionTypes, page, size)
	if err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}
