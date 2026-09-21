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
// đang chạy trực tiếp trên host hay trong container.
type RuntimeSelection struct {
	Environment   string
	ExecutionMode string
	// ConfigProfile is the YAML stem loaded as config/<profile>.yml
	// (from ENV_RUNTIME, or derived when only QHPRO_ENVIRONMENT is set).
	ConfigProfile string
}

// ResolveRuntimeSelection ánh xạ environment sang selection. Compatibility
// ENV_RUNTIME chỉ được đọc khi QHPRO_ENVIRONMENT chưa tồn tại.
func ResolveRuntimeSelection() (RuntimeSelection, error) {
	rawEnvironment := strings.ToLower(strings.TrimSpace(os.Getenv(EnvironmentKey)))
	rawMode := strings.ToLower(strings.TrimSpace(os.Getenv(ExecutionModeKey)))
	if rawEnvironment == "" {
		selection, err := resolveLegacyProfile(strings.ToLower(strings.TrimSpace(os.Getenv(legacyProfileKey))))
		if err != nil {
			return RuntimeSelection{}, err
		}
		selection.ConfigProfile = configProfileName(selection)
		return selection, nil
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

	selection := RuntimeSelection{
		Environment:   environment,
		ExecutionMode: rawMode,
	}
	selection.ConfigProfile = configProfileName(selection)
	return selection, nil
}

// LoadRuntimeYML resolve environment rồi nạp đúng một file:
// config/<ENV_RUNTIME>.yml (aliases chuẩn hoá về stem file thực tế).
// CONFIG_FILE / QHPRO_CONFIG_FILE vẫn override path tuyệt đối khi operator mount.
func LoadRuntimeYML() (RuntimeSelection, error) {
	selection, err := ResolveRuntimeSelection()
	if err != nil {
		return RuntimeSelection{}, err
	}
	if err := LoadYMLFile(selection.ConfigProfile); err != nil {
		return RuntimeSelection{}, err
	}
	return selection, nil
}

// configProfileName chọn stem YAML. ENV_RUNTIME thắng khi được set (kể cả khi
// QHPRO_ENVIRONMENT cũng có), vì đó là tên file operator muốn load khi run.
func configProfileName(selection RuntimeSelection) string {
	if raw := strings.TrimSpace(os.Getenv(legacyProfileKey)); raw != "" {
		return normalizeConfigProfile(raw)
	}
	switch selection.Environment {
	case "development":
		if selection.ExecutionMode == "container" {
			return "develop"
		}
		return "local"
	case "test":
		return "test"
	case "staging":
		return "staging"
	case "production":
		return "production"
	default:
		return "local"
	}
}

func normalizeConfigProfile(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "local":
		return "local"
	case "dev", "develop", "development", "docker":
		return "develop"
	case "test", "testing":
		return "test"
	case "stage", "staging":
		return "staging"
	case "prod", "product", "production":
		return "production"
	case "runtime":
		return "runtime"
	default:
		return strings.ToLower(strings.TrimSpace(raw))
	}
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
