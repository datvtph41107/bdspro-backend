package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const providerModeEnv = "QHPRO_AI_PROVIDER_MODE"

type ProviderConfig struct {
	APIKey    string
	BaseURL   string
	Model     string
	Timeout   time.Duration
	MaxTokens int
}

type Runtime struct {
	GRPCPort     string
	ProviderMode string

	Deepseek ProviderConfig
	OpenAI   ProviderConfig
	Gemini   ProviderConfig
}

// NewRuntime materializes Assistant process configuration exactly once for the
// Wire graph. configloader.LoadRuntimeYML must already have loaded runtime.yml
// and established YAML/environment precedence at the process boundary.
func NewRuntime() (Runtime, error) {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv(providerModeEnv)))
	if mode == "" {
		mode = "live"
	}
	switch mode {
	case "live", "stub", "disabled", "off":
	default:
		return Runtime{}, fmt.Errorf("unsupported %s=%q (want stub or live)", providerModeEnv, mode)
	}

	grpcPort := strings.TrimSpace(viper.GetString("app.port.grpc"))
	if grpcPort == "" {
		grpcPort = "50061"
	}

	return Runtime{
		GRPCPort:     grpcPort,
		ProviderMode: mode,
		Deepseek: providerConfig(
			"DEEPSEEK_API_KEY", "deepseek",
			"https://api.deepseek.com/v1", "deepseek-chat", 30*time.Second, 4000,
		),
		OpenAI: providerConfig(
			"OPENAI_API_KEY", "openai",
			"https://api.openai.com/v1", "gpt-4o-mini", 60*time.Second, 4000,
		),
		Gemini: providerConfig(
			"GEMINI_API_KEY", "gemini",
			"https://generativelanguage.googleapis.com/v1beta", "gemini-1.5-flash", 60*time.Second, 4000,
		),
	}, nil
}

func providerConfig(
	apiKeyEnv string,
	prefix string,
	defaultBaseURL string,
	defaultModel string,
	defaultTimeout time.Duration,
	defaultMaxTokens int,
) ProviderConfig {
	apiKey := strings.TrimSpace(os.Getenv(apiKeyEnv))
	if apiKey == "" {
		apiKey = strings.TrimSpace(viper.GetString(prefix + ".api_key"))
	}

	baseURL := strings.TrimSpace(viper.GetString(prefix + ".base_url"))
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	model := strings.TrimSpace(viper.GetString(prefix + ".model"))
	if model == "" {
		model = defaultModel
	}

	timeout := time.Duration(viper.GetInt(prefix+".timeout_seconds")) * time.Second
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	maxTokens := viper.GetInt(prefix + ".max_tokens")
	if maxTokens <= 0 {
		maxTokens = defaultMaxTokens
	}

	return ProviderConfig{
		APIKey:    apiKey,
		BaseURL:   baseURL,
		Model:     model,
		Timeout:   timeout,
		MaxTokens: maxTokens,
	}
}
