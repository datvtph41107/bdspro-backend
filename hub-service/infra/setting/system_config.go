package setting

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
)

// SystemConfigRepo interface để tránh circular dependency
type SystemConfigRepo interface {
	GetAll(ctx context.Context) ([]*SystemConfigEntity, error)
	Create(ctx context.Context, config *SystemConfigEntity) (*SystemConfigEntity, error)
}

// SystemConfigEntity simplified struct
type SystemConfigEntity struct {
	ID    uint64
	Name  string
	Key   string
	Value string
	Group int32
}

// SystemConfig map lưu config trong memory
var (
	SystemConfig     = make(map[string]string)
	systemConfigLock sync.RWMutex
)

// DefaultSystemConfig chứa các giá trị mặc định
var DefaultSystemConfig = map[string]struct {
	Name  string
	Value string
	Group int32
}{
	"date_format": {
		Name:  "Định dạng ngày tháng",
		Value: "DD/MM/YYYY",
		Group: 20, // Format
	},
	"price_format": {
		Name:  "Định dạng giá tiền",
		Value: "10",
		Group: 20, // Format
	},
	"currency": {
		Name:  "Đơn vị tiền tệ",
		Value: "VND",
		Group: 20, // Format
	},
	"timezone": {
		Name:  "Múi giờ",
		Value: "Asia/Ho_Chi_Minh",
		Group: 30, // Localization
	},
	"language": {
		Name:  "Ngôn ngữ mặc định",
		Value: "vi",
		Group: 30, // Localization
	},
	"items_per_page": {
		Name:  "Số item trên mỗi trang",
		Value: "20",
		Group: 10, // General
	},
	"max_upload_size": {
		Name:  "Kích thước upload tối đa (MB)",
		Value: "10",
		Group: 40, // Upload
	},
}

// GetSystemConfig lấy giá trị config từ memory
func GetSystemConfig(key string) string {
	systemConfigLock.RLock()
	defer systemConfigLock.RUnlock()
	return SystemConfig[key]
}

// GetSystemConfigWithDefault lấy giá trị config từ memory, nếu không có thì trả về default
func GetSystemConfigWithDefault(key string, defaultValue string) string {
	systemConfigLock.RLock()
	defer systemConfigLock.RUnlock()
	if value, exists := SystemConfig[key]; exists {
		return value
	}
	return defaultValue
}

// SetSystemConfig set giá trị config vào memory
func SetSystemConfig(key, value string) {
	systemConfigLock.Lock()
	defer systemConfigLock.Unlock()
	SystemConfig[key] = value
}

// InitSystemConfig khởi tạo config khi startup
func InitSystemConfig(ctx context.Context, repo SystemConfigRepo) error {
	slog.InfoContext(ctx, strings.TrimSuffix(fmt.Sprintln("Initializing system config..."), "\n"))

	// Lấy tất cả config từ DB
	configs, err := repo.GetAll(ctx)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Error loading system config from DB: %v", err))
		// Nếu lỗi, load default values
		loadDefaultConfig()
		return nil
	}

	// Nếu DB trống, seed default values
	if len(configs) == 0 {
		slog.InfoContext(ctx, strings.TrimSuffix(fmt.Sprintln("Database is empty, seeding default system config..."), "\n"))
		if err := seedDefaultConfig(ctx, repo); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("Error seeding default config: %v", err))
			loadDefaultConfig()
			return nil
		}
		// Load lại từ DB sau khi seed
		configs, _ = repo.GetAll(ctx)
	}

	// Load config vào memory
	systemConfigLock.Lock()
	defer systemConfigLock.Unlock()
	for _, config := range configs {
		SystemConfig[config.Key] = config.Value
	}
	slog.InfoContext(ctx, fmt.Sprintf("Loaded %d system config entries", len(SystemConfig)))
	return nil
}

// ReloadSystemConfig reload config từ DB vào memory
func ReloadSystemConfig(ctx context.Context, repo SystemConfigRepo) (int, error) {
	configs, err := repo.GetAll(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to reload system config: %w", err)
	}

	systemConfigLock.Lock()
	defer systemConfigLock.Unlock()

	// Clear old config
	SystemConfig = make(map[string]string)

	// Load new config
	for _, config := range configs {
		SystemConfig[config.Key] = config.Value
	}
	slog.InfoContext(ctx, fmt.Sprintf("Reloaded %d system config entries", len(SystemConfig)))
	return len(SystemConfig), nil
}

// loadDefaultConfig load default config vào memory
func loadDefaultConfig() {
	systemConfigLock.Lock()
	defer systemConfigLock.Unlock()
	for key, config := range DefaultSystemConfig {
		SystemConfig[key] = config.Value
	}
	slog.Info(fmt.Sprintf("Loaded %d default config entries into memory", len(DefaultSystemConfig)))
}

// seedDefaultConfig seed default config vào DB
func seedDefaultConfig(ctx context.Context, repo SystemConfigRepo) error {
	for key, config := range DefaultSystemConfig {
		entity := &SystemConfigEntity{
			Name:  config.Name,
			Key:   key,
			Value: config.Value,
			Group: config.Group,
		}
		if _, err := repo.Create(ctx, entity); err != nil {
			return fmt.Errorf("failed to seed config %s: %w", key, err)
		}
	}
	slog.InfoContext(ctx, fmt.Sprintf("Seeded %d default config entries into database", len(DefaultSystemConfig)))
	return nil
}
