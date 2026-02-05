package provider

import (
	"core/kernel/contract"

	"github.com/spf13/viper"
)

type viperConfig struct {
	v *viper.Viper
}

func NewViperConfig(v *viper.Viper) contract.Config {
	return &viperConfig{v: v}
}

func (c *viperConfig) Get(key string) any          { return c.v.Get(key) }
func (c *viperConfig) GetString(key string) string { return c.v.GetString(key) }
func (c *viperConfig) GetInt(key string) int       { return c.v.GetInt(key) }
func (c *viperConfig) GetBool(key string) bool     { return c.v.GetBool(key) }
