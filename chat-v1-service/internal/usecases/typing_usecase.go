package usecases

import (
	"context"

	_provider "chat/internal/interface/provider"
	"chat/utils"
)

type conversationTypingValidator interface {
	ValidateConversationAndCurrentUser(ctx context.Context, conversationId uint64) error
}

type TypingUsecase struct {
	validator conversationTypingValidator
	counter   _provider.ITypingCounter
}

func NewTypingUsecase(validator conversationTypingValidator, counter _provider.ITypingCounter) *TypingUsecase {
	return &TypingUsecase{validator: validator, counter: counter}
}

// ApplyTypingState kiểm tra quyền trong hội thoại rồi INCR/DECR Redis theo isTyping.
func (u *TypingUsecase) ApplyTypingState(ctx context.Context, conversationID uint64, isTyping bool) (int64, error) {
	if err := u.validator.ValidateConversationAndCurrentUser(ctx, conversationID); err != nil {
		return 0, err
	}
	userID := utils.GetCurrentUserID(ctx)
	return u.counter.AdjustTypingRef(ctx, conversationID, userID, isTyping)
}
