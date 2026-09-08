package usecase

import (
	"context"
	"fmt"
	"notification/internal/domain"
	)

type DealHistoryUsecase struct {
	dealHistoryRepo DealHistoryStore
}

func NewDealHistoryUsecase(dealHistoryRepo DealHistoryStore) *DealHistoryUsecase {
	return &DealHistoryUsecase{
		dealHistoryRepo: dealHistoryRepo,
	}
}

func (u *DealHistoryUsecase) CreateDealHistory(ctx context.Context, req *domain.DealHistoryRequest) (*domain.DealHistoryResponse, error) {
	// Convert request to entity
	history := &domain.DealHistoryEntity{
		DealID:      req.DealID,
		ActorID:     req.ActorID,
		ActorName:   req.ActorName,
		ActorAvatar: req.ActorAvatar,
		ActionType:  req.ActionType,
		ActionName:  req.ActionName,
		Content:     req.Content,
		Metadata:    req.Metadata,
	}

	// Create in database
	createdHistory, err := u.dealHistoryRepo.Create(ctx, history)
	if err != nil {
		return &domain.DealHistoryResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create deal history: %v", err),
		}, err
	}

	return &domain.DealHistoryResponse{
		Success: true,
		Message: "Deal history created successfully",
		ID:      createdHistory.ID,
	}, nil
}

func (u *DealHistoryUsecase) GetDealHistory(ctx context.Context, dealID uint64, page, size int) ([]*domain.DealHistoryEntity, uint32, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20 // Default page size
	}

	histories, total, err := u.dealHistoryRepo.GetByDealID(ctx, dealID, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get deal history: %w", err)
	}

	return histories, total, nil
}

func (u *DealHistoryUsecase) GetDealHistoryByID(ctx context.Context, id uint64) (*domain.DealHistoryEntity, error) {
	history, err := u.dealHistoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get deal history by id: %w", err)
	}

	return history, nil
}
