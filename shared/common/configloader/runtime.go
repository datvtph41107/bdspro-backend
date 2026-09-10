package configloader

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	EnvironmentKey   = "QHPRO_ENVIRONMENT"
	ExecutionModeKey = "QHPRO_EXECUTION_MODE"
	legacyProfileKey = "ENV_RUNTIME"
)

// RuntimeSelection là quyết định cấu hình duy nhất tại process boundary.
// Environment diễn đạt vòng đời deploy; ExecutionMode chỉ diễn đạt process
// đang chạy trực tiếp trên host hay trong container. Không trường nào trong
// selection được dùng để chọn một committed YAML topology khác.
type RuntimeSelection struct {
	Environment   string
	ExecutionMode string
}

// ResolveRuntimeSelection ánh xạ một environment business-independent sang
// đúng YAML topology. Compatibility ENV_RUNTIME chỉ được đọc khi canonical key
// chưa tồn tại, để deployment cũ có thể chuyển đổi dần.
func ResolveRuntimeSelection() (RuntimeSelection, error) {
	rawEnvironment := strings.ToLower(strings.TrimSpace(os.Getenv(EnvironmentKey)))
	rawMode := strings.ToLower(strings.TrimSpace(os.Getenv(ExecutionModeKey)))
	if rawEnvironment == "" {
		return resolveLegacyProfile(strings.ToLower(strings.TrimSpace(os.Getenv(legacyProfileKey))))
	}

	environment := normalizeEnvironment(rawEnvironment)
	if environment == "" {
		return RuntimeSelection{}, fmt.Errorf("invalid %s %q", EnvironmentKey, rawEnvironment)
	}
	if rawMode == "" {
		rawMode = "host"
	}
	if rawMode != "host" && rawMode != "container" {
		return RuntimeSelection{}, fmt.Errorf("invalid %s %q", ExecutionModeKey, rawMode)
	}

	return RuntimeSelection{
		Environment:   environment,
		ExecutionMode: rawMode,
	}, nil
}

// LoadRuntimeYML resolve environment đúng một lần rồi nạp YAML tương ứng.
// Canonical defaults sống trong config/runtime.yml. Nếu developer có
// config/local.yml thì file đó được merge đè lên (secrets/DSN local) mà không
// cần export .env khi chạy air/make run-dev.
func LoadRuntimeYML() (RuntimeSelection, error) {
	selection, err := ResolveRuntimeSelection()
	if err != nil {
		return RuntimeSelection{}, err
	}
	if err := LoadYMLFile("runtime"); err != nil {
		return RuntimeSelection{}, err
	}
	if err := mergeLocalYMLIfPresent(); err != nil {
		return RuntimeSelection{}, err
	}
	return selection, nil
}

const localConfigPath = "config/local.yml"

// mergeLocalYMLIfPresent overlays optional workstation overrides onto the
// already-loaded runtime.yml. Missing local.yml is normal for CI/production.
func mergeLocalYMLIfPresent() error {
	// Operator-mounted config is already the single source of truth.
	if strings.TrimSpace(os.Getenv("QHPRO_CONFIG_FILE")) != "" || strings.TrimSpace(os.Getenv("CONFIG_FILE")) != "" {
		return nil
	}
	data, err := os.ReadFile(localConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read config %s: %w", localConfigPath, err)
	}
	if err := viper.MergeConfig(bytes.NewReader(data)); err != nil {
		return fmt.Errorf("merge config %s: %w", localConfigPath, err)
	}
	return nil
}

func normalizeEnvironment(value string) string {
	switch value {
	case "dev", "develop", "development", "local":
		return "development"
	case "test", "testing":
		return "test"
	case "stage", "staging":
		return "staging"
	case "prod", "product", "production":
		return "production"
	default:
		return ""
	}
}

func resolveLegacyProfile(profile string) (RuntimeSelection, error) {
	switch profile {
	case "", "local":
		return RuntimeSelection{Environment: "development", ExecutionMode: "host"}, nil
	case "dev", "develop", "development", "docker":
		return RuntimeSelection{Environment: "development", ExecutionMode: "container"}, nil
	case "test", "testing":
		return RuntimeSelection{Environment: "test", ExecutionMode: "host"}, nil
	case "stage", "staging":
		return RuntimeSelection{Environment: "staging", ExecutionMode: "container"}, nil
	case "prod", "product", "production":
		return RuntimeSelection{Environment: "production", ExecutionMode: "container"}, nil
	default:
		return RuntimeSelection{}, fmt.Errorf("invalid %s %q", legacyProfileKey, profile)
	}
}
