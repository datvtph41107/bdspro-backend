package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	_crud "common/domain/crud"
	_err "common/domain/err"
	"hub/internal/domain"
	_repo "hub/internal/repo"
)

// IApiKeyUsecase định nghĩa các hành vi nghiệp vụ cho API key.
type IApiKeyUsecase interface {
	CreateApiKey(ctx context.Context, entity *domain.ApiKeyEntity) (*domain.ApiKeyEntity, *_err.ErrorDTO)
	GetApiKeys(ctx context.Context) ([]domain.ApiKeyEntity, *_err.ErrorDTO)
	DeleteApiKey(ctx context.Context, id uint64) *_err.ErrorDTO
	VerifyApiKey(ctx context.Context, apiKey string) (*domain.ApiKeyEntity, *_err.ErrorDTO)
}

// ApiKeyUsecase triển khai nghiệp vụ quản lý API key.
type ApiKeyUsecase struct {
	_crud.BaseUsecase[domain.ApiKeyEntity, _repo.IApiKeyRepo]
}

func NewApiKeyUsecase(repo _repo.IApiKeyRepo) IApiKeyUsecase {
	return &ApiKeyUsecase{
		BaseUsecase: _crud.BaseUsecase[domain.ApiKeyEntity, _repo.IApiKeyRepo]{
			Repo: repo,
		},
	}
}

func (u *ApiKeyUsecase) CreateApiKey(ctx context.Context, entity *domain.ApiKeyEntity) (*domain.ApiKeyEntity, *_err.ErrorDTO) {
	if entity == nil {
		return nil, &_err.ErrorDTO{Code: 400, Message: "Dữ liệu không hợp lệ"}
	}
	if entity.Name == "" {
		return nil, &_err.ErrorDTO{Code: 400, Message: "Tên key không được để trống"}
	}
	if entity.AppName == "" {
		return nil, &_err.ErrorDTO{Code: 400, Message: "Tên ứng dụng không được để trống"}
	}

	if entity.ExpiredAt != nil {
		if entity.ExpiredAt.Before(time.Now()) {
			return nil, &_err.ErrorDTO{Code: 400, Message: "Thời gian hết hạn phải lớn hơn hiện tại"}
		}
	}

	generatedKey, genErr := generateAPIKey()
	if genErr != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: genErr.Error()}
	}
	entity.ApiKey = generatedKey

	if err := u.Repo.Create(ctx, entity); err != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: err.Error()}
	}

	return entity, nil
}

func (u *ApiKeyUsecase) GetApiKeys(ctx context.Context) ([]domain.ApiKeyEntity, *_err.ErrorDTO) {
	entities, err := u.Repo.GetAll(ctx)
	if err != nil {
		return nil, &_err.ErrorDTO{Code: 500, Message: err.Error()}
	}
	return entities, nil
}

func (u *ApiKeyUsecase) DeleteApiKey(ctx context.Context, id uint64) *_err.ErrorDTO {
	if id == 0 {
		return &_err.ErrorDTO{Code: 400, Message: "Thiếu ID"}
	}

	if err := u.Repo.Delete(ctx, id); err != nil {
		if errors.Is(err, _repo.ErrApiKeyNotFound) {
			return &_err.ErrorDTO{Code: 404, Message: "Không tìm thấy API key"}
		}
		return &_err.ErrorDTO{Code: 500, Message: err.Error()}
	}

	return nil
}

func (u *ApiKeyUsecase) VerifyApiKey(ctx context.Context, apiKey string) (*domain.ApiKeyEntity, *_err.ErrorDTO) {
	if apiKey == "" {
		return nil, &_err.ErrorDTO{Code: 400, Message: "Thiếu API key"}
	}

	entity, err := u.Repo.GetByValue(ctx, apiKey)
	if err != nil {
		if errors.Is(err, _repo.ErrApiKeyNotFound) {
			return nil, &_err.ErrorDTO{Code: 404, Message: "API key không hợp lệ"}
		}
		return nil, &_err.ErrorDTO{Code: 500, Message: err.Error()}
	}

	if entity.ExpiredAt != nil && entity.ExpiredAt.Before(time.Now()) {
		return nil, &_err.ErrorDTO{Code: 401, Message: "API key đã hết hạn"}
	}

	return entity, nil
}

func generateAPIKey() (string, error) {
	buff := make([]byte, 32)
	if _, err := rand.Read(buff); err != nil {
		return "", err
	}
	return hex.EncodeToString(buff), nil
}
