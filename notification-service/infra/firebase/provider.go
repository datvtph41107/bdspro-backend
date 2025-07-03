package firebase

import (
	_db "common/db"
	"context"
	"log"

	"noti/config"
	"noti/internal/domain"
	"noti/internal/dto"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// FirebaseProvider xử lý gửi thông báo đến Firebase
type FirebaseProvider struct {
	firebaseApp *firebase.App
}

// NewFirebaseService tạo một instance mới của FirebaseService
func NewFirebaseService(firebaseConfig *config.FirebaseConfig) *FirebaseProvider {
	// firebaseApp.Messaging(ctx)
	return &FirebaseProvider{
		firebaseApp: firebaseConfig.App,
	}
}

// SendToTopic gửi thông báo đến một topic
func (s *FirebaseProvider) SendToTopic(c *gin.Context, request *dto.SendNotiRequest) (*dto.SendNotiResponse, error) {
	var entity domain.HistoryFcmEntity
	copier.Copy(&entity, &request)
	ctx := context.Background()

	notification := s.dtoToNotification(request)
	messageBody := &messaging.Message{
		Topic:        request.Target,
		Notification: notification,
	}
	client, err := s.firebaseApp.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	result, err := client.Send(ctx, messageBody)
	if err != nil {
		log.Println("Error sending message:", err)
		entity.Error = err.Error()
	} else {
		entity.Response = result
	}

	_db.SaveWithAudit(c, entity)
	return &dto.SendNotiResponse{Response: result}, nil
}

// SendToToken gửi thông báo đến một thiết bị dựa vào token
func (s *FirebaseProvider) SendToToken(c *gin.Context, request *dto.SendNotiRequest) (*dto.SendNotiResponse, error) {
	var entity domain.HistoryFcmEntity
	copier.Copy(&entity, &request)
	ctx := context.Background()

	notification := s.dtoToNotification(request)
	message := &messaging.Message{
		Token:        request.Target,
		Notification: notification,
	}

	client, err := s.firebaseApp.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	result, err := client.Send(ctx, message)
	if err != nil {
		log.Println("Error sending message:", err)
		entity.Error = err.Error()
	} else {
		entity.Response = result
	}

	_db.SaveWithAudit(c, &entity)
	return &dto.SendNotiResponse{Response: result}, nil
}

// dtoToNotification chuyển đổi SendNotiRequest thành Firebase Notification
func (s *FirebaseProvider) dtoToNotification(request *dto.SendNotiRequest) *messaging.Notification {
	return &messaging.Notification{
		Title:    request.Title,
		Body:     request.Body,
		ImageURL: request.Image,
	}
}
