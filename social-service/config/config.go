package config

import (
	"bytes"
	_ "embed"
	"log"
	"social/infra/service"

	"github.com/spf13/viper"
)

const C_SESSION_ID = "C_SESSION_ID"

//go:embed app.yml
var ConfigFile []byte

// UserProperties chứa cấu hình của ứng dụng
type UserProperties struct {
	Jwt struct {
		SecretKey           string   `mapstructure:"key-generate"`
		TokenPrefix         string   `mapstructure:"tokenPrefix"`
		TokenExpirationDays int      `mapstructure:"tokenExpirationAfterDays"`
		AccessExpMinutes    int      `mapstructure:"accessExpAfterMinutes"`
		RefreshExpMinutes   uint64   `mapstructure:"refreshExpAfterMinutes"`
		ListPermit          []string `mapstructure:"listPermit"`
		AuthorizationHeader string   `mapstructure:"authorizationHeader"`
	} `mapstructure:"jwt"`
}

var AppProperties UserProperties

type ConfigApp struct {
	NewsFeedService *service.NewsFeedService
	ReportService   *service.ReportService
	CommentService  *service.CommentService
	LikeService     *service.LikeService
}

func NewConfigApp(
	newsFeedService *service.NewsFeedService,
	reportService *service.ReportService,
	commentService *service.CommentService,
	likeService *service.LikeService,
) *ConfigApp {
	return &ConfigApp{
		NewsFeedService: newsFeedService,
		ReportService:   reportService,
		CommentService:  commentService,
		LikeService:     likeService,
	}
}
func LoadConfig() {
	// viper.SetConfigName("config")   // Tên file không có phần mở rộng
	viper.SetConfigType("yaml") // Loại file config
	// viper.AddConfigPath("./config") // Đọc từ thư mục gốc

	// Đọc biến môi trường (nếu có)
	viper.AutomaticEnv()

	// Đọc file config
	err := viper.ReadConfig(bytes.NewReader(ConfigFile))
	if err != nil {
		log.Fatalf("Lỗi khi đọc file config: %v", err)
	}

	if err := viper.Unmarshal(&AppProperties); err != nil {
		log.Fatalf("Lỗi parse config: %v", err)
		// return nil, err
	}
}
