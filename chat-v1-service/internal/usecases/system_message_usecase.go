// internal/usecases/system_message.go
package usecases

import (
	"context"
	"fmt"
	"time"

	iusecase "chat/internal/interface"
	"chat/models"
	"chat/utils"
)

type SystemMessageUsecase struct {
	repo       Repository
	userClient iusecase.IUserClient
}

func NewSystemMessageUsecase(repo Repository, userClient iusecase.IUserClient) *SystemMessageUsecase {
	return &SystemMessageUsecase{
		repo:       repo,
		userClient: userClient,
	}
}

func (s *SystemMessageUsecase) PinSystemMessage(ctx context.Context, conversationID uint64, messageID uint64, isPin bool) (*models.MessageModel, error) {
	actorID := utils.GetCurrentUserID(ctx)

	actorName, err := s.getUserName(ctx, actorID)
	if err != nil {
		actorName = fmt.Sprintf("User %d", actorID)
	}

	var content string
	if isPin {
		content = fmt.Sprintf("%s đã ghim tin nhắn", actorName)
	} else {
		content = fmt.Sprintf("%s đã bỏ ghim tin nhắn", actorName)
	}

	systemMsg := &models.MessageModel{
		ConversationID: conversationID,
		SenderID:       actorID,
		Content:        content,
		ContentType:    models.ContentTypeSystem,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	message, err := s.repo.CreateMessage(ctx, systemMsg)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *SystemMessageUsecase) AddMemberSystemMessage(ctx context.Context, conversationID uint64, newUserIDs []uint64) error {
	if len(newUserIDs) == 0 {
		return nil
	}

	actorID := utils.GetCurrentUserID(ctx)

	actorName, err := s.getUserName(ctx, actorID)
	if err != nil {
		actorName = fmt.Sprintf("User %d", actorID)
	}

	var content string
	if len(newUserIDs) == 1 {
		newUserName, err := s.getUserName(ctx, newUserIDs[0])
		if err != nil {
			newUserName = fmt.Sprintf("User %d", newUserIDs[0])
		}
		content = fmt.Sprintf("%s đã thêm %s vào nhóm", actorName, newUserName)
	} else {
		content = fmt.Sprintf("%s đã thêm %d thành viên vào nhóm", actorName, len(newUserIDs))
	}

	systemMsg := &models.MessageModel{
		ConversationID: conversationID,
		SenderID:       actorID,
		Content:        content,
		ContentType:    models.ContentTypeSystem,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, err = s.repo.CreateMessage(ctx, systemMsg)
	return err
}

func (s *SystemMessageUsecase) getUserName(ctx context.Context, userID uint64) (string, error) {
	user, err := s.userClient.GetUserById(ctx, userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("user not found")
	}
	return user.FullName, nil
}
