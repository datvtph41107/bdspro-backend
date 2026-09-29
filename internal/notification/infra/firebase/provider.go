package firebase

import (
	"common/logging"
	"context"
	"fmt"
	"log/slog"

	"notification/config"
	"notification/internal/domain"
	deliverydomain "notification/internal/domain/delivery"
	"notification/internal/dto"

	"gorm.io/gorm"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// FirebaseProvider xử lý gửi thông báo đến Firebase
// @bind: notification/internal/usecase.PushSender
type FirebaseProvider struct {
	firebaseApp *firebase.App
	db          *gorm.DB
}

// NewFirebaseService tạo một instance mới của FirebaseService
func NewFirebaseService(firebaseConfig *config.FirebaseConfig, db *gorm.DB) *FirebaseProvider {
	var app *firebase.App
	if firebaseConfig != nil {
		app = firebaseConfig.App
	}
	return &FirebaseProvider{firebaseApp: app, db: db}
}

// SendToTopic gửi thông báo đến một topic
func (s *FirebaseProvider) SendToTopic(c *gin.Context, request *dto.SendNotiRequest) (*dto.SendNotiResponse, error) {
	if s == nil || s.firebaseApp == nil {
		return nil, fmt.Errorf("firebase push delivery is not configured")
	}
	var entity domain.HistoryFcmEntity
	copier.Copy(&entity, &request)
	if c == nil || c.Request == nil {
		return nil, fmt.Errorf("request context is required")
	}
	ctx := c.Request.Context()

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
		logging.WithComponent(ctx, "firebase").Error(
			"Firebase message send failed",
			slog.Any("error", err),
		)
		entity.Error = err.Error()
	} else {
		entity.Response = result
	}

	if s.db != nil {
		_ = s.db.WithContext(ctx).Create(&entity).Error
	}
	return &dto.SendNotiResponse{Response: result}, nil
}

// SendToToken gửi thông báo đến một thiết bị dựa vào token
func (s *FirebaseProvider) SendToToken(c *gin.Context, request *dto.SendNotiRequest) (*dto.SendNotiResponse, error) {
	if s == nil || s.firebaseApp == nil {
		return nil, fmt.Errorf("firebase push delivery is not configured")
	}
	var entity domain.HistoryFcmEntity
	copier.Copy(&entity, &request)
	if c == nil || c.Request == nil {
		return nil, fmt.Errorf("request context is required")
	}
	ctx := c.Request.Context()

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
		logging.WithComponent(ctx, "firebase").Error(
			"Firebase message send failed",
			slog.Any("error", err),
		)
		entity.Error = err.Error()
	} else {
		entity.Response = result
	}

	if s.db != nil {
		_ = s.db.WithContext(ctx).Create(&entity).Error
	}
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

// SendPushNotification gửi push notification đến nhiều tokens
func (s *FirebaseProvider) SendPushNotification(ctx context.Context, tokens []string, title string, message string, data map[string]string) error {
	if len(tokens) == 0 {
		return &deliverydomain.EffectError{Kind: deliverydomain.EffectFailurePermanent, Err: fmt.Errorf("no tokens provided")}
	}

	if s.firebaseApp == nil {
		return &deliverydomain.EffectError{Kind: deliverydomain.EffectFailurePermanent, Err: fmt.Errorf("firebase app is not initialized")}
	}

	client, err := s.firebaseApp.Messaging(ctx)
	if err != nil {
		return &deliverydomain.EffectError{Kind: deliverydomain.EffectFailureRetryable, Err: fmt.Errorf("get Firebase messaging client: %w", err)}
	}

	// Tạo notification
	notification := &messaging.Notification{
		Title: title,
		Body:  message,
	}

	// Tạo message với data
	messageBody := &messaging.Message{
		Notification: notification,
		Data:         data,
	}

	// Nếu có nhiều tokens, dùng multicast
	if len(tokens) > 1 {
		messages := make([]*messaging.Message, 0, len(tokens))
		for _, token := range tokens {
			msg := &messaging.Message{
				Token:        token,
				Notification: notification,
				Data:         data,
			}
			messages = append(messages, msg)
		}

		// Gửi multicast
		br, err := client.SendMulticast(ctx, &messaging.MulticastMessage{
			Tokens:       tokens,
			Notification: notification,
			Data:         data,
		})
		if err != nil {
			return &deliverydomain.EffectError{Kind: deliverydomain.EffectFailureUnknown, Err: fmt.Errorf("Firebase multicast outcome unknown: %w", err)}
		}

		// Log kết quả
		logging.WithComponent(ctx, "firebase").Info(
			"Firebase multicast sent",
			slog.Int("firebase.success_count", br.SuccessCount),
			slog.Int("firebase.failure_count", br.FailureCount),
		)
		if br.FailureCount > 0 {
			for i, resp := range br.Responses {
				if !resp.Success {
					logging.WithComponent(ctx, "firebase").Warn(
						"Firebase multicast token delivery failed",
						slog.Int("firebase.token_index", i),
						slog.Any("error", resp.Error),
					)
				}
			}
		}
		return nil
	}

	// Nếu chỉ có 1 token, gửi đơn lẻ
	messageBody.Token = tokens[0]
	_, err = client.Send(ctx, messageBody)
	if err != nil {
		return &deliverydomain.EffectError{Kind: deliverydomain.EffectFailureUnknown, Err: fmt.Errorf("Firebase send outcome unknown: %w", err)}
	}

	logging.WithComponent(ctx, "firebase").Info(
		"Firebase push notification sent",
	)
	return nil
}
