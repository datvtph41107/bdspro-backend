package config

import (
	server "bdspro/external/handler"
	public_router "bdspro/external/router/public"
	user_router "bdspro/external/router/user"
	service "bdspro/external/services"
	"bdspro/v1/mod/admin"
	"bdspro/v1/mod/user"
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
	// GrpcServer        *server.BdsproGrpcServer
	GrpcGenericServer *server.GenericGrpcServer
	UserRoute         *user.UserRoute
	AdminRouter       *admin.AdminRoute
	RouterManager     *user_router.RouterManager
	EnumRoute         *public_router.PublicRouter

	// @bind: register_server_config
	CostDocumentServer      *service.CostDocumentServer
	IncomeDocumentServer    *service.IncomeDocumentServer
	AssetIncomeTypeServer   *service.AssetIncomeTypeServer
	AssetExploitationServer *service.AssetExploitationServer
	AssetCostTypeServer     *service.AssetCostTypeServer
	AssetCostServer         *service.AssetCostServer
	AssetLegalServer        *service.AssetLegalServer
	AssetSplitHistoryServer *service.AssetSplitHistoryServer
}

func NewConfigApp(
	GrpcGenericServer *server.GenericGrpcServer,
	userRoute *user.UserRoute,
	adminRoute *admin.AdminRoute,
	// @bind: register_server_param
	CostDocumentServer *service.CostDocumentServer,
	IncomeDocumentServer *service.IncomeDocumentServer,
	AssetIncomeTypeServer *service.AssetIncomeTypeServer,
	AssetExploitationServer *service.AssetExploitationServer,
	AssetCostTypeServer *service.AssetCostTypeServer,
	AssetCostServer *service.AssetCostServer,
	AssetLegalServer *service.AssetLegalServer,
	AssetSplitHistoryServer *service.AssetSplitHistoryServer,
	RouterManager *user_router.RouterManager,
	EnumRoute *public_router.PublicRouter,
) *ConfigApp {
	return &ConfigApp{
		// GrpcServer:        GrpcServer,
		GrpcGenericServer: GrpcGenericServer,
		UserRoute:         userRoute,
		AdminRouter:       adminRoute,
		// @bind: register_server_new
		CostDocumentServer:      CostDocumentServer,
		IncomeDocumentServer:    IncomeDocumentServer,
		AssetIncomeTypeServer:   AssetIncomeTypeServer,
		AssetExploitationServer: AssetExploitationServer,
		AssetCostTypeServer:     AssetCostTypeServer,
		AssetCostServer:         AssetCostServer,
		AssetLegalServer:        AssetLegalServer,
		AssetSplitHistoryServer: AssetSplitHistoryServer,
		RouterManager:           RouterManager,
		EnumRoute:               EnumRoute,
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
