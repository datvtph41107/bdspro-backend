package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Loader struct {
	opts *Options
	env  Env
}

func New(opts ...Option) (*Loader, error) {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	rawEnv := os.Getenv(options.EnvVar)
	if rawEnv == "" {
		rawEnv = "local"
	}

	env, err := ParseEnv(rawEnv)
	if err != nil {
		return nil, err
	}

	return &Loader{
		opts: options,
		env:  env,
	}, nil
}

func (l *Loader) Load(ctx context.Context, out interface{}) error {
	v := viper.New()

	// ===== base =====
	v.SetConfigName(l.opts.ConfigName)
	v.SetConfigType(l.opts.ConfigType)

	for _, path := range l.opts.ConfigPaths {
		v.AddConfigPath(path)
	}

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("base config: %w", err)
	}

	// ===== env file =====
	envFile := fmt.Sprintf("%s.%s", l.opts.ConfigName, l.env)

	vEnv := viper.New()
	vEnv.SetConfigName(envFile)
	vEnv.SetConfigType(l.opts.ConfigType)

	for _, path := range l.opts.ConfigPaths {
		vEnv.AddConfigPath(path)
	}

	if err := vEnv.ReadInConfig(); err == nil {
		_ = v.MergeConfigMap(vEnv.AllSettings())
	}

	// ===== env override =====
	if l.opts.ReadEnv {
		v.SetEnvPrefix(l.opts.EnvPrefix)
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		v.AutomaticEnv()
	}

	// ===== unmarshal =====
	if err := v.Unmarshal(out); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	// ===== validate =====
	if l.opts.Validate {
		if validator, ok := out.(Validator); ok {
			if err := validator.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}
