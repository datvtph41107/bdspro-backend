package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	Service struct {
		Name    string `yaml:"name" validate:"required"`
		Version string `yaml:"version" validate:"required"`
		Env     string `yaml:"env"`
	} `yaml:"service"`

	Server struct {
		HTTP struct {
			Port         int           `yaml:"port" validate:"min=1,max=65535"`
			ReadTimeout  time.Duration `yaml:"read_timeout"`
			WriteTimeout time.Duration `yaml:"write_timeout"`
		} `yaml:"http"`
	} `yaml:"server"`

	Redis struct {
		Host     string `yaml:"host" validate:"required"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db" validate:"min=0,max=15"`
	} `yaml:"redis"`

	JWT struct {
		SecretKey           string   `yaml:"key-generate" validate:"required"`
		TokenPrefix         string   `yaml:"tokenPrefix"`
		TokenExpirationDays int      `yaml:"tokenExpirationAfterDays" validate:"min=1"`
		AccessExpMinutes    int      `yaml:"accessExpAfterMinutes" validate:"min=1"`
		RefreshExpMinutes   uint64   `yaml:"refreshExpAfterMinutes" validate:"min=1"`
		ListPermit          []string `yaml:"listPermit"`
		AuthorizationHeader string   `yaml:"authorizationHeader"`
	} `yaml:"jwt"`
}

func (c *Config) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

func (c *Config) String() string {
	return fmt.Sprintf("Config{Service: %s, Env: %s}",
		c.Service.Name, c.Service.Env)
}
