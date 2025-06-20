package config

import (
	_ "embed"
	"log"
	"user/config/property"
	"user/infra/routes"
	"user/infra/services"

	"github.com/spf13/viper"
)

//go:embed app.yml
var ConfigFile []byte

type ConfigApp struct {
	Auth           *routes.AuthRouter
	OAuthRouter    *routes.OAuthRouter
	ProfileRouter  *routes.ProfileRouter
	BlockRouter    *routes.BlockRouter
	FollowRouter   *routes.FollowRouter
	FriendRouter   *routes.FriendRouter
	GroupRouter    *routes.GroupRouter
	RateRouter     *routes.RateRouter
	ProfileService *services.GrpcProfileService
	AuthService    *services.GrpcAuthService
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
	profileService *services.GrpcProfileService,
	grpcAuthService *services.GrpcAuthService,
) *ConfigApp {
	return &ConfigApp{
		Auth:           authService,
		OAuthRouter:    oAuthRouter,
		ProfileRouter:  profile,
		BlockRouter:    block,
		FollowRouter:   follow,
		FriendRouter:   friend,
		GroupRouter:    group,
		RateRouter:     rate,
		ProfileService: profileService,
		AuthService:    grpcAuthService,
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
