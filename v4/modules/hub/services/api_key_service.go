package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	_errors "common/errors"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/entity"
	"userv4/modules/hub/repository"
)

// APIKeyService manages api_keys CRUD.
type APIKeyService interface {
	CreateAPIKey(ctx context.Context, req dto.CreateAPIKeyRequest) (*dto.APIKeyItem, error)
	GetAPIKeys(ctx context.Context) ([]*dto.APIKeyItem, error)
	DeleteAPIKey(ctx context.Context, id uint64) error
	VerifyAPIKey(ctx context.Context, apiKey string) (*dto.APIKeyItem, error)
}

type LocalAPIKeyService struct {
	repo repository.APIKeyRepository
}

func NewLocalAPIKeyService(repo repository.APIKeyRepository) *LocalAPIKeyService {
	return &LocalAPIKeyService{repo: repo}
}

func (s *LocalAPIKeyService) CreateAPIKey(ctx context.Context, req dto.CreateAPIKeyRequest) (*dto.APIKeyItem, error) {
	if req.Name == "" {
		return nil, _errors.BadRequestException("name is required")
	}
	if req.AppName == "" {
		return nil, _errors.BadRequestException("appName is required")
	}
	if req.ExpiredAt != nil && req.ExpiredAt.Before(time.Now()) {
		return nil, _errors.BadRequestException("expiredAt must be in the future")
	}
	key, err := generateKey()
	if err != nil {
		return nil, err
	}
	apiKey := &entity.APIKey{
		Name: req.Name, AppName: req.AppName, Description: req.Description,
		APIKey: key, ExpiredAt: req.ExpiredAt,
	}
	if err := s.repo.Create(ctx, apiKey); err != nil {
		return nil, err
	}
	return mapAPIKeyItem(apiKey), nil
}

func (s *LocalAPIKeyService) GetAPIKeys(ctx context.Context) ([]*dto.APIKeyItem, error) {
	rows, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*dto.APIKeyItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapAPIKeyItem(r))
	}
	return out, nil
}

func (s *LocalAPIKeyService) DeleteAPIKey(ctx context.Context, id uint64) error {
	if id == 0 {
		return _errors.BadRequestException("id is required")
	}
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, repository.ErrAPIKeyNotFound) {
		return _errors.NotFoundException("api key not found")
	}
	return err
}

func (s *LocalAPIKeyService) VerifyAPIKey(ctx context.Context, apiKey string) (*dto.APIKeyItem, error) {
	if apiKey == "" {
		return nil, _errors.BadRequestException("apiKey is required")
	}
	row, err := s.repo.GetByValue(ctx, apiKey)
	if errors.Is(err, repository.ErrAPIKeyNotFound) {
		return nil, _errors.NotFoundException("api key invalid or not found")
	}
	if err != nil {
		return nil, err
	}
	if row.ExpiredAt != nil && row.ExpiredAt.Before(time.Now()) {
		return nil, _errors.UnauthorizedException("api key expired")
	}
	return mapAPIKeyItem(row), nil
}

func generateKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func mapAPIKeyItem(a *entity.APIKey) *dto.APIKeyItem {
	if a == nil {
		return nil
	}
	item := &dto.APIKeyItem{
		ID: a.ID, Name: a.Name, AppName: a.AppName,
		Description: a.Description, APIKey: a.APIKey, ExpiredAt: a.ExpiredAt,
	}
	if a.CreatedAt != nil {
		s := strconv.FormatInt(a.CreatedAt.UnixMilli(), 10)
		item.CreatedAt = &s
	}
	return item
}
