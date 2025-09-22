package handler

import (
	_jwt "common/jwt"
	_middleware "common/middleware"
	_routes "common/routes"
	"context"
	orgpb "pb/types/organization"
	userpb "pb/types/user"
	user_repo "user/infra/repo"
	auth_service "user/internal/services/auth"
)

type GrpcAuthService struct {
	userpb.UnimplementedAuthServiceServer
	ProfileService *auth_service.AuthService
	OrgClient      orgpb.OrganizationMemberServiceClient
	ProfileRepo    *user_repo.ProfileRepo
}

func NewGrpcAuthService(profileService *auth_service.AuthService,
	profileRepo *user_repo.ProfileRepo,
	orgClient orgpb.OrganizationMemberServiceClient,
) *GrpcAuthService {
	return &GrpcAuthService{
		ProfileService: profileService,
		ProfileRepo:    profileRepo,
		OrgClient:      orgClient,
	}
}

// @Summary Chuyển đổi tài khoản
// @Description Chuyển đổi tài khoản
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body userpb.SwitchAccountRequest true "Thông tin tài khoản"
// @Security BearerAuth
// @Router /v2/user/auth/switch-account [post]
func (s *GrpcAuthService) SwitchAccount(ctx context.Context, req *userpb.SwitchAccountRequest) (*userpb.SwitchAccountResponse, error) {
	principal, err := _middleware.PrincipalFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if principal.ProfileId == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không thể xác thực phiên đăng nhập",
		}
	}

	profile, err := s.ProfileRepo.GetByProfileID(principal.ProfileId)
	if err != nil {
		return nil, err
	}

	ownerType := uint32(10)
	results := struct {
		OrganizationId *uint64
		GroupId        *uint64
		FullName       string
		Avatar         string
		Email          string
		Phone          string
		OwnerType      *uint32
	}{
		OrganizationId: nil,
		FullName:       profile.FullName,
		Avatar:         profile.Avatar,
		Email:          profile.Email,
		Phone:          profile.Phone,
		OwnerType:      &ownerType,
	}

	payload := _jwt.JwtTokenProperties{
		AuthID:    principal.AuthID,
		ProfileID: principal.ProfileId,
		SessionID: uint64(principal.Session),
		Role:      "ROLE_USER",
		Type:      "ACCESS",
		PlanID:    profile.PlanID,
		PlanAt:    profile.PlanAt,
	}

	if req.OrganizationId != nil {
		clientResponse, err := s.OrgClient.SwitchOrganization(ctx, &orgpb.SwitchOrganizationRequest{
			OrganizationId: uint32(*req.OrganizationId),
			UserId:         uint32(principal.ProfileId),
		})
		if err != nil {
			return nil, err
		}
		if !clientResponse.IsSuccess {
			return nil, &_routes.Except{
				Code:    400,
				Message: "Không thể chuyển đổi tổ chức",
			}
		}
		results.FullName = profile.FullName
		results.Avatar = profile.Avatar
		results.Email = profile.Email
		results.Phone = profile.Phone
		ownerType = 30
		results.OwnerType = &ownerType
		payload.AdditionalInfo = map[string]interface{}{
			"organization": clientResponse.Id,
		}
		results.OrganizationId = req.OrganizationId
	}

	// authEntity, err := auth_model.AuthGetWithId(principal.UserID)
	// if err != nil {
	// 	return nil, err
	// }

	newAccessToken, err := _jwt.GenerateToken(payload)
	if err != nil {
		return nil, err
	}

	return &userpb.SwitchAccountResponse{
		AuthId:         principal.AuthID,
		ProfileId:      profile.ProfileID,
		OrganizationId: results.OrganizationId,
		GroupId:        results.GroupId,
		Role:           "ROLE_USER",
		FullName:       results.FullName,
		Avatar:         results.Avatar,
		Email:          results.Email,
		Phone:          results.Phone,
		AccessToken:    newAccessToken,
		// RefreshToken:   s.ProfileService.GetRefreshToken(authEntity, uint64(principal.Session)),
	}, nil
}
