package providers

import (
	"core/kernel.v2/src/contracts/logger"
	"core/kernel.v2/src/pkg/zaplogger"
)

type LoggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func ProvideLogger(cfg LoggerConfig) (logger.Logger, error) {
	return zaplogger.New(logger.Config{
		Level:  logger.Level(cfg.Level),
		Format: cfg.Format,
	})
}
