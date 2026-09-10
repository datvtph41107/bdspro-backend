package initial

import (
	"hub/infra/handler"
	"hub/internal/usecase"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type InitialApp struct {
	DB                        *gorm.DB
	ApiKeyInterceptor         grpc.UnaryServerInterceptor
	EventQueueService         *handler.EventQueueService
	InternalHandler           *handler.InternalHandler
	UpdateDataHandler         *handler.UpdateDataHandler
	LocationHandler           *handler.LocationHandler
	LocationV2Handler         *handler.LocationV2Handler
	UserGuideHandler          *handler.UserGuideHandler
	FAQHandler                *handler.FAQHandler
	SystemConfigHandler       *handler.SystemConfigHandler
	VersionHandler            *handler.VersionHandler
	ApiKeyHandler             *handler.ApiKeyHandler
	SystemConfigUsecase       *usecase.SystemConfigUsecase
	Logger                    *flogging.FabricLogger
	InteractiveEventHandler   *handler.InteractiveEventHandler
	ClickURLHandler           *handler.ClickURLHandler
	ErrorLogHandler           *handler.ErrorLogHandler
	ApplinkHandler            *handler.ApplinkHandler
	RequestStatsHandler       *handler.RequestStatsHandler
	KeywordSearchStatsHandler *handler.KeywordSearchStatsHandler
	SearchLatencyStatsHandler *handler.SearchLatencyStatsHandler
}

func NewInitialApp(
	db *gorm.DB,
	apiKeyInterceptor grpc.UnaryServerInterceptor,
	eventQueueService *handler.EventQueueService,
	internalHandler *handler.InternalHandler,
	updateDataHandler *handler.UpdateDataHandler,
	locationHandler *handler.LocationHandler,
	locationV2Handler *handler.LocationV2Handler,
	userGuideHandler *handler.UserGuideHandler,
	faqHandler *handler.FAQHandler,
	systemConfigHandler *handler.SystemConfigHandler,
	versionHandler *handler.VersionHandler,
	apiKeyHandler *handler.ApiKeyHandler,
	systemConfigUsecase *usecase.SystemConfigUsecase,
	interactiveEventHandler *handler.InteractiveEventHandler,
	clickURLHandler *handler.ClickURLHandler,
	errorLogHandler *handler.ErrorLogHandler,
	appLinkHandler *handler.ApplinkHandler,
	requestStatsHandler *handler.RequestStatsHandler,
	keywordSearchStatsHandler *handler.KeywordSearchStatsHandler,
	searchLatencyStatsHandler *handler.SearchLatencyStatsHandler,
) *InitialApp {
	return &InitialApp{
		DB:                        db,
		ApiKeyInterceptor:         apiKeyInterceptor,
		EventQueueService:         eventQueueService,
		InternalHandler:           internalHandler,
		UpdateDataHandler:         updateDataHandler,
		LocationHandler:           locationHandler,
		LocationV2Handler:         locationV2Handler,
		UserGuideHandler:          userGuideHandler,
		FAQHandler:                faqHandler,
		SystemConfigHandler:       systemConfigHandler,
		VersionHandler:            versionHandler,
		ApiKeyHandler:             apiKeyHandler,
		SystemConfigUsecase:       systemConfigUsecase,
		Logger:                    flogging.MustGetLogger(viper.GetString("server.name")),
		InteractiveEventHandler:   interactiveEventHandler,
		ClickURLHandler:           clickURLHandler,
		ErrorLogHandler:           errorLogHandler,
		ApplinkHandler:            appLinkHandler,
		RequestStatsHandler:       requestStatsHandler,
		KeywordSearchStatsHandler: keywordSearchStatsHandler,
		SearchLatencyStatsHandler: searchLatencyStatsHandler,
	}
}
