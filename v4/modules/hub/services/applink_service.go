package services

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"strconv"

	_errors "common/errors"
	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/entity"
	"mdn4/modules/hub/repository"
)

const (
	applinkCodeMin uint64 = 1_000_000
	applinkCodeMax uint64 = 10_000_000_000
	maxCodeRetry          = 5
)

var errApplinkNotFound = errors.New("applink: not found")
var errApplinkInvalidCode = errors.New("applink: invalid base64 code")
var errApplinkCodeExhausted = errors.New("applink: failed to generate unique code after retries")

// ApplinkService manages applink (deep-link code) creation and lookup.
type ApplinkService interface {
	CreateApplink(ctx context.Context, req dto.CreateApplinkRequest) (*dto.ApplinkResponse, error)
	GetApplinkByCode(ctx context.Context, b64Code string) (*dto.ApplinkResponse, error)
}

// @bind: mdn4/modules/hub/services.ApplinkService
type LocalApplinkService struct {
	repo repository.ApplinkRepository
}

func NewLocalApplinkService(repo repository.ApplinkRepository) *LocalApplinkService {
	return &LocalApplinkService{repo: repo}
}

func (s *LocalApplinkService) CreateApplink(ctx context.Context, req dto.CreateApplinkRequest) (*dto.ApplinkResponse, error) {
	if req.RefID == 0 {
		return nil, _errors.BadRequestException("ref_id is required")
	}

	// Return existing if same refID + action already exists.
	existing, err := s.repo.GetByRefIDAndAction(ctx, req.RefID, req.Action)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return toApplinkResponse(existing), nil
	}

	code, err := s.generateUniqueCode(ctx)
	if err != nil {
		return nil, err
	}
	a := &entity.Applink{Code: code, RefID: req.RefID, Action: req.Action}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return toApplinkResponse(a), nil
}

func (s *LocalApplinkService) GetApplinkByCode(ctx context.Context, b64Code string) (*dto.ApplinkResponse, error) {
	code, err := base64DecodeCode(b64Code)
	if err != nil {
		return nil, _errors.BadRequestException(err.Error())
	}
	a, err := s.repo.GetByCode(ctx, code)
	if errors.Is(err, repository.ErrApplinkNotFound) {
		return nil, _errors.NotFoundException("applink not found")
	}
	if err != nil {
		return nil, err
	}
	return toApplinkResponse(a), nil
}

func (s *LocalApplinkService) generateUniqueCode(ctx context.Context) (uint64, error) {
	for i := 0; i < maxCodeRetry; i++ {
		rang := applinkCodeMax - applinkCodeMin + 1
		code := applinkCodeMin + uint64(rand.Int63n(int64(rang)))
		exists, err := s.repo.ExistsByCode(ctx, code)
		if err != nil {
			return 0, err
		}
		if !exists {
			return code, nil
		}
	}
	return 0, errApplinkCodeExhausted
}

func base64DecodeCode(b64 string) (uint64, error) {
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", errApplinkInvalidCode, err.Error())
	}
	code, err := strconv.ParseUint(string(raw), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: not a valid number", errApplinkInvalidCode)
	}
	return code, nil
}

func base64EncodeCode(code uint64) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.FormatUint(code, 10)))
}

func toApplinkResponse(a *entity.Applink) *dto.ApplinkResponse {
	return &dto.ApplinkResponse{
		ID:     a.ID,
		Code:   base64EncodeCode(a.Code),
		RefID:  a.RefID,
		Action: a.Action,
	}
}
