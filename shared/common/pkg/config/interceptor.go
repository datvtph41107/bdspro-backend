package config

import "time"

type InterceptorConfig struct {
	RecoverPanic   bool          `mapstructure:"recover_panic"`
	RequestTimeout time.Duration `mapstructure:"request_timeout"  validate:"min=1s,max=120s"`
	MaxRetry       int           `mapstructure:"max_retry"         validate:"min=0,max=5"`
	RetryDelay     time.Duration `mapstructure:"retry_delay"`
	RateLimit      int           `mapstructure:"rate_limit"        reload:"hot"`
	RateLimitBurst int           `mapstructure:"rate_limit_burst"  reload:"hot"`
}
