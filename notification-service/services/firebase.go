package services

import (
	_db "common/db"
	"context"
	"log"

	"noti/config"
	"noti/data"
	"noti/models"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// FirebaseService xử lý gửi thông báo đến Firebase
type FirebaseService struct {
	firebaseApp *firebase.App
}

// NewFirebaseService tạo một instance mới của FirebaseService
func NewFirebaseService(firebaseConfig *config.FirebaseConfig) *FirebaseService {
	// firebaseApp.Messaging(ctx)
	return &FirebaseService{
		firebaseApp: firebaseConfig.App,
	}
}

// SendToTopic gửi thông báo đến một topic
func (s *FirebaseService) SendToTopic(c *gin.Context, request *data.SendNotiRequest) (*data.SendNotiResponse, error) {
	var entity models.HistoryEntity
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
	return &data.SendNotiResponse{Response: result}, nil
}

// SendToToken gửi thông báo đến một thiết bị dựa vào token
func (s *FirebaseService) SendToToken(c *gin.Context, request *data.SendNotiRequest) (*data.SendNotiResponse, error) {
	var entity models.HistoryEntity
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
	return &data.SendNotiResponse{Response: result}, nil
}

// dtoToNotification chuyển đổi SendNotiRequest thành Firebase Notification
func (s *FirebaseService) dtoToNotification(request *data.SendNotiRequest) *messaging.Notification {
	return &messaging.Notification{
		Title:    request.Title,
		Body:     request.Body,
		ImageURL: request.Image,
	}
}
