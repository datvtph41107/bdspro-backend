package config

import (
	"bdspro/v1/mod/admin"
	"bdspro/v1/mod/user"
	server "bdspro/v2/external/handler"
	"bytes"
	_ "embed"
	"log"

	"github.com/spf13/viper"
)

const C_SESSION_ID = "C_SESSION_ID"

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

// LoadConfig đọc cấu hình từ file `.env` hoặc `config.yaml`
// func LoadProperties() {
// 	log.Printf("Đọc config")

// 	// var config UserProperties
// 	if err := viper.Unmarshal(&AppProperties); err != nil {
// 		log.Fatalf("Lỗi parse config: %v", err)
// 		// return nil, err
// 	}

// 	// return &config, nil
// }

//go:embed app.yml
var ConfigFile []byte

type ConfigApp struct {
	GrpcServer  *server.BdsproGrpcServer
	UserRoute   *user.UserRoute
	AdminRouter *admin.AdminRoute
}

func NewConfigApp(GrpcServer *server.BdsproGrpcServer,
	userRoute *user.UserRoute,
	adminRoute *admin.AdminRoute) *ConfigApp {
	return &ConfigApp{
		GrpcServer:  GrpcServer,
		UserRoute:   userRoute,
		AdminRouter: adminRoute,
	}
}

// LoadConfig đọc cấu hình từ file config.yml
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
