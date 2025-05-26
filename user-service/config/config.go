package config

import (
	_ "embed"
	"log"
	"user/config/property"
	"user/routes"

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
	// GroupRouter   *routes.GroupRouter
}

func NewConfigApp(authService *routes.AuthRouter,
	oAuthRouter *routes.OAuthRouter,
	profile *routes.ProfileRouter,
) *ConfigApp {
	return &ConfigApp{
		Auth:          authService,
		OAuthRouter:   oAuthRouter,
		ProfileRouter: profile,
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
