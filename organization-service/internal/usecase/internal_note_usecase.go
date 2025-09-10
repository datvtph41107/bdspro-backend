package usecase

import (
	_utils "common/utils"
	"context"
	"fmt"
	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/enums"
)

type InternalNoteUsecase interface {
	AddInternalNote(ctx context.Context, dealID uint64, content string, actionType enums.ActionType) (*entity.InternalNote, error)
	GetDealHistory(ctx context.Context, dealID uint64, actionTypes []enums.ActionType, page, size uint32) ([]*entity.InternalNote, uint32, error)
}

type internalNoteUsecase struct {
	internalNoteRepository repository.InternalNoteRepository
	dealRepository         repository.DealRepository
	groupAuthUsecase       GroupAuthUsecase
	logWorker              *LogWorker
}

func NewInternalNoteUsecase(
	internalNoteRepository repository.InternalNoteRepository,
	dealRepository repository.DealRepository,
	groupAuthUsecase GroupAuthUsecase,
	logWorker *LogWorker,
) InternalNoteUsecase {
	return &internalNoteUsecase{
		internalNoteRepository: internalNoteRepository,
		dealRepository:         dealRepository,
		groupAuthUsecase:       groupAuthUsecase,
		logWorker:              logWorker,
	}
}

func (u *internalNoteUsecase) AddInternalNote(ctx context.Context, dealID uint64, content string, actionType enums.ActionType) (*entity.InternalNote, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra deal có tồn tại không
	deal, err := u.dealRepository.GetByID(ctx, dealID)
	if err != nil {
		return nil, err
	}
	if deal == nil {
		return nil, custom_error.RecordNotFound("deal not found")
	}

	// Kiểm tra quyền - chỉ leader mới được thêm ghi chú nội bộ
	err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to add internal note")
	}

	// Tạo ghi chú nội bộ
	note := &entity.InternalNote{
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
	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		ActorId: uint32(currentUserId),
		LogType: "ADD_INTERNAL_NOTE",
		LogData: fmt.Sprintf("Added internal note for deal %d", dealID),
	})

	return createdNote, nil
}

func (u *internalNoteUsecase) GetDealHistory(ctx context.Context, dealID uint64, actionTypes []enums.ActionType, page, size uint32) ([]*entity.InternalNote, uint32, error) {
	// Kiểm tra deal có tồn tại không
	deal, err := u.dealRepository.GetByID(ctx, dealID)
	if err != nil {
		return nil, 0, err
	}
	if deal == nil {
		return nil, 0, custom_error.RecordNotFound("deal not found")
	}

	// Lấy lịch sử
	notes, total, err := u.internalNoteRepository.GetByDealID(ctx, dealID, actionTypes, page, size)
	if err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}
