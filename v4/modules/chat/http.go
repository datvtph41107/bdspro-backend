package chat

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	commonredis "common/redis"
	commonutils "common/utils"
	v4config "mdn4/config"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gorm.io/gorm"

	bdsproenum "mdn4/modules/bdspro/enums"
	bdsprorepository "mdn4/modules/bdspro/repository"
	chatconfig "mdn4/modules/chat/config"
	"mdn4/modules/chat/infrastructure/client"
	"mdn4/modules/chat/infrastructure/handlers"
	"mdn4/modules/chat/infrastructure/mapper"
	middlewares "mdn4/modules/chat/infrastructure/middleware"
	chatdto "mdn4/modules/chat/internal/dto"
	"mdn4/modules/chat/internal/repository"
	"mdn4/modules/chat/internal/usecases"
	userusecases "mdn4/modules/user/services/user"
	bdspropb "pb/types/bdspro"
)

// NewDelivery wires the migrated chat API into v4 using its database and shared configuration.
func NewDelivery(db *gorm.DB, profileUsecase *userusecases.ProfileUsecase, sharedRedis *commonredis.RedisService) (*handlers.Delivery, func(), error) {
	if db == nil {
		return nil, nil, fmt.Errorf("chat: database is required")
	}
	if profileUsecase == nil {
		return nil, nil, fmt.Errorf("chat: v4 user profile usecase is required")
	}
	chatconfig.AppProperties.Database.Type = "POSTGRES"
	chatconfig.AppProperties.MaxParticipant = 10
	chatconfig.AppProperties.Worker.Number = 10
	chatconfig.AppProperties.Redis.Host = v4config.Properties.Redis.Host
	chatconfig.AppProperties.Redis.Password = v4config.Properties.Redis.Password

	cleanup := func() {}

	chatRepository := repository.NewRepository(db)
	userClient := client.NewUserClient(func(ctx context.Context, ids []uint64) (map[uint64]*chatdto.UserProfile, error) {
		profiles, err := profileUsecase.ProfileFindByIds(ctx, ids)
		if err != nil {
			return nil, err
		}
		mapped := make(map[uint64]*chatdto.UserProfile)
		if profiles == nil {
			return mapped, nil
		}
		for _, profile := range *profiles {
			mapped[profile.ProfileID] = &chatdto.UserProfile{
				ID:           profile.ProfileID,
				FullName:     profile.FullName,
				Avatar:       profile.Avatar,
				TickVerified: profile.TickVerified,
			}
		}
		return mapped, nil
	})
	productRepo := bdsprorepository.NewProductPostgres(db)
	bdsproClient := client.NewBdsproClient(func(ctx context.Context, ids []uint64) (map[uint64]*bdspropb.ProductAttachment, error) {
		products, err := productRepo.ListByIDs(ctx, ids)
		if err != nil {
			return nil, err
		}
		attachments := make(map[uint64]*bdspropb.ProductAttachment, len(products))
		for _, product := range products {
			if product == nil {
				continue
			}
			attachment := &bdspropb.ProductAttachment{
				Id: product.ID, Name: product.Name, Code: product.Code,
				Area: product.Area, AssetId: product.AssetID,
				RentStatus: product.RentStatus, SaleStatus: product.SaleStatus,
				SourceType:     int32(product.SourceType),
				SourceTypeName: bdsproenum.SourceTypeMap[product.SourceType],
			}
			if product.Price != nil && product.Price.SalePrice != nil {
				attachment.Price = *product.Price.SalePrice
			}
			if product.Property != nil && product.Property.Info != nil && product.Property.Info.Avatar != nil {
				attachment.Image = product.Property.Info.Avatar.MediaURL
			}
			attachments[product.ID] = attachment
		}
		return attachments, nil
	})
	conversationUsecase := usecases.NewConversationUsecases(chatRepository, userClient)
	syncProvider := commonutils.NewSyncUtil(sharedRedis)
	messageUsecase := usecases.NewMessageUsecases(chatRepository, conversationUsecase, syncProvider)
	participantUsecase := usecases.NewParticipantUsecases(chatRepository)
	readReceiptUsecase := usecases.NewReadReceptUsecases(chatRepository)
	reactionUsecase := usecases.NewMessageReactionUsecases(chatRepository)
	backgroundUsecase := usecases.NewBackgroundImageUsecases(chatRepository)
	backgroundMapper := mapper.NewBackgroundMapper()
	messageMapper := mapper.NewMessageMapper()
	conversationMapper := mapper.NewConversationMapper(backgroundMapper, messageMapper)
	delivery := handlers.NewDelivery(conversationUsecase, messageUsecase, readReceiptUsecase,
		participantUsecase, reactionUsecase, backgroundUsecase, userClient,
		bdsproClient, conversationMapper, messageMapper, backgroundMapper,
		syncProvider, chatRepository,
		usecases.NewSystemMessageUsecase(chatRepository, userClient))
	return delivery, cleanup, nil
}

// NewHTTPHandler registers every chat-v1 HTTP transcoding route under /v2/chat.
func NewHTTPHandler(db *gorm.DB, profileUsecase *userusecases.ProfileUsecase) (http.Handler, func(), error) {
	sharedRedis := commonredis.NewRedisService()
	delivery, deliveryCleanup, err := NewDelivery(db, profileUsecase, sharedRedis)
	if err != nil {
		_ = sharedRedis.Close()
		return nil, nil, err
	}
	moduleCleanup := func() {
		deliveryCleanup()
		_ = sharedRedis.Close()
	}
	mux := runtime.NewServeMux()
	if err := delivery.RegisterHTTPServer(context.Background(), mux); err != nil {
		moduleCleanup()
		return nil, nil, fmt.Errorf("register chat HTTP API: %w", err)
	}
	httpMux := http.NewServeMux()
	httpMux.Handle("POST /v2/chat/session", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		profileID := commonutils.GetProfileIdWithContext(r.Context())
		if profileID == 0 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var randomBytes [32]byte
		if _, err := rand.Read(randomBytes[:]); err != nil {
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}
		key := base64.RawURLEncoding.EncodeToString(randomBytes[:])
		redisClient, err := sharedRedis.Instance()
		if err != nil {
			http.Error(w, "session service unavailable", http.StatusServiceUnavailable)
			return
		}
		if err := redisClient.Set(r.Context(), "chat:ss:"+key, strconv.FormatUint(profileID, 10), 5*time.Second).Err(); err != nil {
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"key": key})
	}))
	httpMux.Handle("/", mux)
	return middlewares.PresenterMiddleware(httpMux), moduleCleanup, nil
}
