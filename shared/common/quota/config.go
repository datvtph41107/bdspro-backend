package quota

import (
	"github.com/spf13/viper"
)

type Config struct {
	Enabled    bool   `mapstructure:"enabled"`
	PolicyFile string `mapstructure:"policy_file"`
}

func LoadQuotaConfig() (*Config, error) {
	var cfg Config
	if err := viper.UnmarshalKey("quota", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
