package config

import (
	"common/configloader"
	"fmt"

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
	Email struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		From     string `mapstructure:"from"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"email"`

	Seo struct {
		PublicBaseURL         string `mapstructure:"public_base_url"`
		InternalSecretKey     string `mapstructure:"internal_secret_key"`
		RuntimeEventIngestKey string `mapstructure:"runtime_event_ingest_key"`
		Timezone              string `mapstructure:"timezone"`

		Sitemap struct {
			DefaultChangeFreq string  `mapstructure:"default_change_freq"`
			DefaultPriority   float32 `mapstructure:"default_priority"`
		} `mapstructure:"sitemap"`

		GenerationJob struct {
			Enabled   bool   `mapstructure:"enabled"`
			Cron      string `mapstructure:"cron"`
			BatchSize int    `mapstructure:"batch_size"`
		} `mapstructure:"generation_job"`
	} `mapstructure:"seo"`
}

var AppProperties UserProperties

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig() error {
	if _, err := configloader.LoadRuntimeYML(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&AppProperties); err != nil {
		return fmt.Errorf("parse CRM config: %w", err)
	}
	return nil
}
