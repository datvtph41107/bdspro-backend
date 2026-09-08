package usecase

import (
	"context"
	"log"
	"user/internal/dto"
	"user/internal/interface/repo"

	_utils "common/utils"
)

type OAuthUsecase struct {
	AuthMethodRepo repo.AuthMethodRepository
}

func NewOAuthUsecase(authMethodRepo repo.AuthMethodRepository) *OAuthUsecase {
	return &OAuthUsecase{
		AuthMethodRepo: authMethodRepo,
	}
}

// GetConnectedAccounts lấy danh sách các tài khoản MXH đã liên kết
func (u *OAuthUsecase) GetConnectedAccounts(ctx context.Context) ([]*dto.ConnectedAccountDTO, error) {
	// Lấy profileId từ context
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		log.Printf("ProfileID not found in context")
		return []*dto.ConnectedAccountDTO{}, nil
	}

	// Lấy danh sách OAuth accounts đã liên kết
	accounts, err := u.AuthMethodRepo.GetConnectedOAuthAccounts(ctx, profileID)
	if err != nil {
		log.Printf("Failed to get connected OAuth accounts for user %d: %v", profileID, err)
		return nil, err
	}

	// Convert sang DTO
	var result []*dto.ConnectedAccountDTO
	for _, account := range accounts {
		result = append(result, &dto.ConnectedAccountDTO{
			AuthID:    account.ID,
			Provider:  string(account.Provider),
			AuthName:  account.AuthName,
			FullName:  account.FullName,
			Email:     account.Email,
			Avatar:    account.Avatar,
			CreatedAt: _utils.FormatTimeToString(&account.CreatedAt),
		})
	}

	log.Printf("Found %d connected OAuth accounts for user %d", len(result), profileID)
	return result, nil
}
