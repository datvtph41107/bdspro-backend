package usecase

import (
	_errors "common/errors"
	"context"
	"fmt"
	"hub/internal"
	"hub/internal/domain"
	"hub/internal/dto"
	"hub/internal/enums"
	"hub/internal/repo"
)

type SystemConfigUsecase struct {
	persistData dto.SystemConfigPersist
	defaultData dto.SystemConfigPersist
	repo        repo.ISystemConfigRepo
}

func NewSystemConfigUsecase(persistData dto.SystemConfigPersist, repo repo.ISystemConfigRepo) *SystemConfigUsecase {
	return &SystemConfigUsecase{
		persistData: persistData,
		defaultData: persistData,
		repo:        repo,
	}
}

// GetSystemConfigsByGroup lấy danh sách system config theo groupKey
func (uc *SystemConfigUsecase) GetSystemConfigsByGroup(ctx context.Context, groupKey string) ([]*domain.SystemConfigEntity, error) {
	// Parse groupKey string to enum
	configGroup, err := enums.GetSystemConfigGroup(groupKey)
	if err != nil {
		return nil, _errors.ReturnError(service.SystemConfigGroupInvalid, _errors.WithPublicMessage(fmt.Sprintf("Group key không hợp lệ: %s", groupKey)))
	}

	configs, err := uc.repo.GetByGroup(ctx, configGroup)
	if err != nil {
		return nil, fmt.Errorf("get system configs by group: %w", err)
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
			return nil, 0, 0, _errors.ReturnError(service.SystemConfigItemGroupInvalid, _errors.WithPublicMessage(fmt.Sprintf("Group key không hợp lệ cho config %s: %s", item.Key, item.GroupKey)))
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
		return nil, 0, 0, fmt.Errorf("bulk upsert system config: %w", err)
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
		return fmt.Errorf("load system config from database: %w", err)
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
		return nil, 0, _errors.ReturnError(service.SystemConfigGroupInvalid, _errors.WithPublicMessage(fmt.Sprintf("Group key không hợp lệ: %s", groupKey)))
	}

	// Lấy tất cả config mặc định từ property.DefaultSystemConfig theo group
	var defaultConfigs []dto.UpsertConfigItem
	// TODO: Implement logic to get default configs by group from property.DefaultSystemConfig

	// Nếu không có config mặc định nào cho group này
	if len(defaultConfigs) == 0 {
		return nil, 0, _errors.ReturnError(service.SystemConfigDefaultNotFound, _errors.WithPublicMessage(fmt.Sprintf("Không tìm thấy config mặc định cho group: %s", enums.GetSystemConfigGroupName(configGroup))))
	}

	// Thực hiện bulk upsert để reset về giá trị mặc định vào DB
	req := &dto.BulkUpsertSystemConfigRequest{
		Configs: defaultConfigs,
	}

	results, _, _, err := uc.BulkUpsertSystemConfig(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("reset system config: %w", err)
	}

	// persistData đã được update tự động trong BulkUpsertSystemConfig
	return results, len(results), nil
}

// GetSystemConfigByKey lấy system config theo key
func (uc *SystemConfigUsecase) GetSystemConfigByKey(ctx context.Context, key string) (*domain.SystemConfigEntity, error) {
	// Validate key
	if key == "" {
		return nil, _errors.ReturnError(service.SystemConfigKeyRequired)
	}

	// Repository technical failures pass through. Only normalized absence owns
	// the legacy business not-found mapping at the application boundary.
	config, err := uc.repo.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, _errors.ReturnError(service.SystemConfigNotFound, _errors.WithPublicMessage(fmt.Sprintf("Không tìm thấy config với key: %s", key)))
	}

	return config, nil
}

// GetUserSettingsByGroup lấy settings theo group dạng key-value map cho user
func (uc *SystemConfigUsecase) GetUserSettingsByGroup(ctx context.Context, groupKey string) (map[string]string, error) {
	// Parse groupKey string to enum
	configGroup, err := enums.GetSystemConfigGroup(groupKey)
	if err != nil {
		return nil, _errors.ReturnError(service.SystemConfigGroupInvalid, _errors.WithPublicMessage(fmt.Sprintf("Group key không hợp lệ: %s", groupKey)))
	}

	// Lấy configs từ DB
	configs, err := uc.repo.GetByGroup(ctx, configGroup)
	if err != nil {
		return nil, fmt.Errorf("get system settings: %w", err)
	}

	// Convert sang map key-value
	settings := make(map[string]string)
	for _, config := range configs {
		settings[config.Key] = config.Value
	}

	return settings, nil
}
