package usecase

import (
	_errors "common/errors"
	"common/pkg/crypto"
	"context"
	"fmt"
	"hub/internal/domain"
	"hub/internal/dto"
	"hub/internal/enums"
	"hub/internal/repo"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type SystemConfigUsecase struct {
	persistData dto.SystemConfigPersist
	defaultData dto.SystemConfigPersist
	repo        repo.ISystemConfigRepo
	encryptor   crypto.Encryptor
}

func NewSystemConfigUsecase(persistData dto.SystemConfigPersist, repo repo.ISystemConfigRepo) *SystemConfigUsecase {
	return &SystemConfigUsecase{
		persistData: persistData,
		defaultData: persistData,
		repo:        repo,
		encryptor:   getEncryptor(),
	}
}

var (
	globalEncryptor     crypto.Encryptor
	globalEncryptorOnce sync.Once
)

func getEncryptor() crypto.Encryptor {
	globalEncryptorOnce.Do(func() {
		// Initialize key manager
		keyManager := crypto.GetGlobalKeyManager()

		// Try to initialize from config
		ctx := context.Background()
		if err := keyManager.InitializeFromConfig(ctx); err != nil {
			// Fallback to env if config fails
			passphrase := os.Getenv("CONFIG_ENCRYPTION_PASSPHRASE")
			if passphrase == "" {
				// Log warning but don't crash - encryption will be disabled
				// In production, you might want to panic here
				globalEncryptor = nil
				return
			}
			keyManager.InitializeWithPassphrase(ctx, passphrase)
		}

		// Create encryptor
		encryptor, err := keyManager.GetActiveEncryptor(ctx)
		if err != nil {
			globalEncryptor = nil
			return
		}
		globalEncryptor = encryptor
	})
	return globalEncryptor
}

// GetEncryptedSystemConfigByKey lấy config theo key và trả về giá trị đã mã hóa
func (uc *SystemConfigUsecase) GetEncryptedSystemConfigByKey(ctx context.Context, key string) (*dto.SystemConfigResponse, error) {
	// Validate key
	if key == "" {
		return nil, _errors.BadRequestException("key is required")
	}

	// Lấy config từ DB
	config, err := uc.repo.GetByKey(ctx, key)
	if err != nil {
		return nil, _errors.NotFoundException(fmt.Sprintf("config not found: %s", key))
	}

	// Kiểm tra encryption có enabled không
	if !viper.GetBool("encryption.enabled") {
		return &dto.SystemConfigResponse{
			Key:   config.Key,
			Value: config.Value,
			// Encrypted: false,
			// KeyID:     "",
			// Version:   "none",
		}, nil
	}

	// Kiểm tra encryptor có sẵn sàng không
	if uc.encryptor == nil {
		// Encryption enabled but encryptor not available - return error
		return nil, _errors.InternalServerException("encryption service not available")
	}

	// Mã hóa value
	encryptedValue, err := uc.encryptor.Encrypt(ctx, config.Value)
	if err != nil {
		return nil, _errors.InternalServerException(fmt.Sprintf("encryption failed: %v", err))
	}

	return &dto.SystemConfigResponse{
		Key:   config.Key,
		Value: encryptedValue,
		// Encrypted: true,
		// KeyID:     uc.encryptor.GetKeyID(),
		// Version:   "v1",
	}, nil
}

// GetEncryptedUserSettingsByGroup lấy settings đã mã hóa theo group
func (uc *SystemConfigUsecase) GetEncryptedUserSettingsByGroup(ctx context.Context, groupKey string) (map[string]string, string, bool, error) {
	// Parse groupKey
	configGroup, err := enums.GetSystemConfigGroup(groupKey)
	if err != nil {
		return nil, "", false, _errors.BadRequestException(fmt.Sprintf("invalid group key: %s", groupKey))
	}

	// Lấy configs từ DB
	configs, err := uc.repo.GetByGroup(ctx, configGroup)
	if err != nil {
		return nil, "", false, _errors.InternalServerException("failed to get settings")
	}

	// Kiểm tra encryption
	encryptionEnabled := viper.GetBool("encryption.enabled")

	// Convert sang map key-value
	settings := make(map[string]string)
	for _, cfg := range configs {
		if encryptionEnabled && uc.encryptor != nil {
			// Mã hóa value
			encrypted, err := uc.encryptor.Encrypt(ctx, cfg.Value)
			if err != nil {
				// Log error, use empty string
				settings[cfg.Key] = ""
				continue
			}
			settings[cfg.Key] = encrypted
		} else {
			settings[cfg.Key] = cfg.Value
		}
	}

	keyID := ""
	if uc.encryptor != nil {
		keyID = uc.encryptor.GetKeyID()
	}

	return settings, keyID, encryptionEnabled && uc.encryptor != nil, nil
}

// GetSystemConfigsByGroup lấy danh sách system config theo groupKey
func (uc *SystemConfigUsecase) GetSystemConfigsByGroup(ctx context.Context, groupKey string) ([]*domain.SystemConfigEntity, error) {
	// Parse groupKey string to enum
	configGroup, err := enums.GetSystemConfigGroup(groupKey)
	if err != nil {
		return nil, _errors.ReturnError(400, fmt.Sprintf("Group key không hợp lệ: %s", groupKey))
	}

	configs, err := uc.repo.GetByGroup(ctx, configGroup)
	if err != nil {
		return nil, _errors.ReturnError(500, "Lỗi khi lấy danh sách system config theo group")
	}
	return configs, nil
}

// BulkUpsertSystemConfig bulk upsert system config
func (uc *SystemConfigUsecase) BulkUpsertSystemConfig(ctx context.Context, req *dto.BulkUpsertSystemConfigRequest) ([]*domain.SystemConfigEntity, int, int, error) {
	// Validate and convert DTO to entities
	configs := make([]*domain.SystemConfigEntity, len(req.Configs))
	for i, item := range req.Configs {
		// Parse groupKey string to enum
		configGroup, err := enums.GetSystemConfigGroup(item.GroupKey)
		if err != nil {
			return nil, 0, 0, _errors.ReturnError(400, fmt.Sprintf("Group key không hợp lệ cho config %s: %s", item.Key, item.GroupKey))
		}

		configs[i] = &domain.SystemConfigEntity{
			Name:        item.Name,
			Key:         item.Key,
			Value:       item.Value,
			GroupConfig: configGroup,
		}
	}

	// Execute bulk upsert
	results, createdCount, updatedCount, err := uc.repo.BulkUpsert(ctx, configs)
	if err != nil {
		return nil, 0, 0, _errors.ReturnError(500, fmt.Sprintf("Lỗi khi bulk upsert system config: %v", err))
	}

	// Cập nhật lại persistData từ DB
	uc.LoadDataFromDB(ctx)

	return results, createdCount, updatedCount, nil
}

// InitDataFromDB load dữ liệu từ DB lên và cập nhật vào persistData
func (uc *SystemConfigUsecase) LoadDataFromDB(ctx context.Context) error {
	// Lấy tất cả config từ DB
	configs, err := uc.repo.GetAll(ctx)
	if err != nil {
		return _errors.ReturnError(500, fmt.Sprintf("Lỗi khi lấy danh sách config từ DB: %v", err))
	}

	// Clear map hiện tại và cập nhật lại
	uc.persistData.Configs = make(map[string]string)

	// Duyệt qua các config và ghép key:value vào persistData
	for _, config := range configs {
		uc.persistData.Configs[config.Key] = config.Value
	}

	return nil
}

// InitializeSystemConfig khởi tạo system config mặc định từ persistData vào DB nếu chưa có
func (uc *SystemConfigUsecase) InitializeSystemConfig(ctx context.Context) error {
	// Load lại toàn bộ config từ DB và cập nhật vào persistData
	if err := uc.LoadDataFromDB(ctx); err != nil {
		return err
	}

	// Load config vào memory provider
	return nil
}

// ResetDefaultByGroup reset config về giá trị mặc định theo groupKey
func (uc *SystemConfigUsecase) ResetDefaultByGroup(ctx context.Context, groupKey string) ([]*domain.SystemConfigEntity, int, error) {
	// Parse groupKey string to enum
	configGroup, err := enums.GetSystemConfigGroup(groupKey)
	if err != nil {
		return nil, 0, _errors.ReturnError(400, fmt.Sprintf("Group key không hợp lệ: %s", groupKey))
	}

	// Lấy tất cả config mặc định từ property.DefaultSystemConfig theo group
	var defaultConfigs []dto.UpsertConfigItem
	// TODO: Implement logic to get default configs by group from property.DefaultSystemConfig

	// Nếu không có config mặc định nào cho group này
	if len(defaultConfigs) == 0 {
		return nil, 0, _errors.ReturnError(404, fmt.Sprintf("Không tìm thấy config mặc định cho group: %s", enums.GetSystemConfigGroupName(configGroup)))
	}

	// Thực hiện bulk upsert để reset về giá trị mặc định vào DB
	req := &dto.BulkUpsertSystemConfigRequest{
		Configs: defaultConfigs,
	}

	results, _, _, err := uc.BulkUpsertSystemConfig(ctx, req)
	if err != nil {
		return nil, 0, _errors.ReturnError(500, fmt.Sprintf("Lỗi khi reset config: %v", err))
	}

	// persistData đã được update tự động trong BulkUpsertSystemConfig
	return results, len(results), nil
}

// GetSystemConfigByKey lấy system config theo key
func (uc *SystemConfigUsecase) GetSystemConfigByKey(ctx context.Context, key string) (*domain.SystemConfigEntity, error) {
	// Validate key
	if key == "" {
		return nil, _errors.ReturnError(400, "Key không được để trống")
	}

	// Lấy config từ DB
	config, err := uc.repo.GetByKey(ctx, key)
	if err != nil {
		return nil, _errors.ReturnError(404, fmt.Sprintf("Không tìm thấy config với key: %s", key))
	}

	return config, nil
}

// GetUserSettingsByGroup lấy settings theo group dạng key-value map cho user
func (uc *SystemConfigUsecase) GetUserSettingsByGroup(ctx context.Context, groupKey string) (map[string]string, error) {
	// Parse groupKey string to enum
	configGroup, err := enums.GetSystemConfigGroup(groupKey)
	if err != nil {
		return nil, _errors.ReturnError(400, fmt.Sprintf("Group key không hợp lệ: %s", groupKey))
	}

	// Lấy configs từ DB
	configs, err := uc.repo.GetByGroup(ctx, configGroup)
	if err != nil {
		return nil, _errors.ReturnError(500, "Lỗi khi lấy settings")
	}

	// Convert sang map key-value
	settings := make(map[string]string)
	for _, config := range configs {
		settings[config.Key] = config.Value
	}

	return settings, nil
}
