package config

import (
	_ "embed"
	"log"
	"user/config/property"
	"user/infra/handler"
	"user/infra/routes"
	"user/infra/services"

	"github.com/spf13/viper"
)

//go:embed app.yml
var ConfigFile []byte

type ConfigApp struct {
	Auth          *routes.AuthRouter
	OAuthRouter   *routes.OAuthRouter
	ProfileRouter *routes.ProfileRouter
	BlockRouter   *routes.BlockRouter
	FollowRouter  *routes.FollowRouter
	FriendRouter  *routes.FriendRouter
	GroupRouter   *routes.GroupRouter
	RateRouter    *routes.RateRouter
	// V2 Handlers (protobuf-based)
	ReportHandler       *handler.ReportHandler
	ReportReasonHandler *handler.ReportReasonHandler
	ReportProofHandler  *handler.ReportProofHandler
	ProfileService      *services.GrpcProfileService
	AuthService         *services.GrpcAuthService
	RateService         *services.RateService
	// GroupRouter   *routes.GroupRouter
}

func NewConfigApp(authService *routes.AuthRouter,
	oAuthRouter *routes.OAuthRouter,
	profile *routes.ProfileRouter,
	block *routes.BlockRouter,
	follow *routes.FollowRouter,
	friend *routes.FriendRouter,
	group *routes.GroupRouter,
	rate *routes.RateRouter,
	reportHandler *handler.ReportHandler,
	reportReasonHandler *handler.ReportReasonHandler,
	reportProofHandler *handler.ReportProofHandler,
	profileService *services.GrpcProfileService,
	grpcAuthService *services.GrpcAuthService,
	rateService *services.RateService,
) *ConfigApp {
	return &ConfigApp{
		Auth:                authService,
		OAuthRouter:         oAuthRouter,
		ProfileRouter:       profile,
		BlockRouter:         block,
		FollowRouter:        follow,
		FriendRouter:        friend,
		GroupRouter:         group,
		RateRouter:          rate,
		ReportHandler:       reportHandler,
		ReportReasonHandler: reportReasonHandler,
		ReportProofHandler:  reportProofHandler,
		ProfileService:      profileService,
		AuthService:         grpcAuthService,
		RateService:         rateService,
	}
}

// LoadConfig đọc cấu hình từ file `.env` hoặc `config.yaml`
func LoadProperties() {
	log.Printf("Đọc config")

	// var config UserProperties
	if err := viper.Unmarshal(&property.AppProperties); err != nil {
		log.Fatalf("Lỗi parse config: %v", err)
		// return nil, err
	}

	// return &config, nil
}
