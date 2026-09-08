package configloader

import (
	"fmt"
	"os"
	"strings"
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
func LoadRuntimeYML() (RuntimeSelection, error) {
	selection, err := ResolveRuntimeSelection()
	if err != nil {
		return RuntimeSelection{}, err
	}
	if err := LoadYMLFile("runtime"); err != nil {
		return RuntimeSelection{}, err
	}
	return selection, nil
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
