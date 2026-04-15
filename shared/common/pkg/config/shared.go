// shared/pkg/config/shared.go
// Config được chia sẻ giữa nhiều service — phải nhất quán
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// SharedAuthConfig: JWT config phải giống nhau ở TẤT CẢ service
// Nếu gateway dùng secret A để sign, chat-service dùng secret B để verify → fail
//
// Giải pháp: định nghĩa một struct chung, mỗi service embed nó
// Operator chỉ cần set 1 env var: SHARED_AUTH_JWT_SECRET
// Tất cả service đọc cùng giá trị
type SharedAuthConfig struct {
	JWTSecret       string        `mapstructure:"jwt_secret"        validate:"required,min=32" sensitive:"true"`
	JWTIssuer       string        `mapstructure:"jwt_issuer"        validate:"required"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"  validate:"min=5m"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl" validate:"min=1h"`
}

// ServiceDiscovery: địa chỉ của các service — thay vì hardcode trong từng service
// Operator chỉ cần maintain một bộ địa chỉ, tất cả service đều biết
//
// Trong môi trường không có K8s (như hệ thống của bạn), đây là "poor man's service registry"
// Đủ tốt cho small-medium scale. Khi cần K8s, chỉ cần thay đổi loader,
// không cần thay đổi code của service.
type ServiceDiscovery struct {
	ChatServiceAddr    string `mapstructure:"chat_service_addr"    validate:"required"`
	CRMServiceAddr     string `mapstructure:"crm_service_addr"     validate:"required"`
	PostServiceAddr    string `mapstructure:"post_service_addr"    validate:"required"`
	GatewayServiceAddr string `mapstructure:"gateway_service_addr" validate:"required"`
}

// LoadShared: load config chung từ prefix SHARED_
// Không phụ thuộc vào tên service nào
func LoadShared(target interface{}) error {
	v := viper.New()

	v.SetEnvPrefix("SHARED")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetConfigName("shared")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("read shared config: %w", err)
		}
	}

	if err := v.Unmarshal(target); err != nil {
		return fmt.Errorf("unmarshal shared config: %w", err)
	}

	return validateStruct(target)
}
