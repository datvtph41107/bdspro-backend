package services

import (
	_utils "common/utils"
	"context"
	"fmt"
	organizationpb "pb/types/organization"
	orgpb "pb/types/organization"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"time"
	user_repo "user/infra/repo"
	user_service "user/internal/services/user"

	"github.com/jinzhu/copier"
	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcProfileService struct {
	userpb.UnimplementedProfileServiceServer
	ProfileRepo *user_repo.ProfileRepo
	Profile     *user_service.ProfileService
	cache       *cache.Cache
	OrgClient   orgpb.OrganizationServiceClient
	GroupClient organizationpb.GroupServiceClient
}

func NewGrpcProfileService(
	profileRepo *user_repo.ProfileRepo,
	profileService *user_service.ProfileService,
	orgClient orgpb.OrganizationServiceClient,
	groupClient organizationpb.GroupServiceClient,
) *GrpcProfileService {
	return &GrpcProfileService{
		ProfileRepo: profileRepo,
		Profile:     profileService,
		cache:       cache.New(5*time.Minute, 10*time.Minute),
		OrgClient:   orgClient,
		GroupClient: groupClient,
	}
}

func (s *GrpcProfileService) GetProfileByIds(ctx context.Context, req *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error) {
	var profilesPb []*sharepb.ProfileItem
	var missingIDs []uint64

	// 1. Check cache trước
	for _, id := range req.Ids {
		if cached, found := s.cache.Get(fmt.Sprintf("%d", id)); found {
			profilesPb = append(profilesPb, cached.(*sharepb.ProfileItem))
		} else {
			missingIDs = append(missingIDs, id)
		}
	}

	// 2. Lấy từ DB nếu thiếu
	if len(missingIDs) > 0 {
		profiles, err := s.ProfileRepo.GetProfileByIds(missingIDs)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get profiles: %v", err)
		}

		for _, profile := range profiles {
			pb := &sharepb.ProfileItem{
				Id:           profile.ProfileID,
				FullName:     profile.FullName,
				Avatar:       profile.Avatar,
				TickVerified: profile.TickVerified,
				Phone:        profile.Phone,
			}
			// Cache lại
			s.cache.Set(fmt.Sprintf("%d", profile.ProfileID), pb, cache.DefaultExpiration)
			profilesPb = append(profilesPb, pb)
		}
	}

	return &sharepb.GetProfileByIdsResponse{
		Profiles: profilesPb,
	}, nil
}

func (s *GrpcProfileService) GetProfileByPhones(ctx context.Context, req *userpb.GetProfileByPhonesRequest) (*sharepb.GetProfileByIdsResponse, error) {
	profiles, err := s.ProfileRepo.GetProfileByPhones(req.Phones)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get profiles: %v", err)
	}

	var profilesPb []*sharepb.ProfileItem
	for _, profile := range profiles {
		profilesPb = append(profilesPb, &sharepb.ProfileItem{
			Id:           profile.ProfileID,
			FullName:     profile.FullName,
			Avatar:       profile.Avatar,
			TickVerified: profile.TickVerified,
			Phone:        profile.Phone,
		})
	}

	return &sharepb.GetProfileByIdsResponse{
		Profiles: profilesPb,
	}, nil
}

// @Summary Lấy thông tin công khai của người dùng
// @Description Lấy thông tin công khai của người dùng
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Profile"
// @Success 200 {object} userpb.ProfilePublicInfo "Thông tin công khai của người dùng"
// @Router /v2/user/profile/info/{id} [get]
func (s *GrpcProfileService) GetProfileInfo(ctx context.Context, req *userpb.IdRequest) (*userpb.ProfilePublicInfo, error) {
	result, err := s.Profile.PublicUserInfo(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get profile info: %v", err)
	}

	dto := &userpb.ProfilePublicInfo{}
	copier.Copy(dto, result)
	return dto, nil
}

// @Summary Lấy danh sách tài khoản của người dùng
// @Description Lấy danh sách tài khoản của người dùng
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} userpb.AccountInfoResponse "Danh sách tài khoản của người dùng"
// @Router /v2/user/profile/accounts/me [get]
func (s *GrpcProfileService) ListAccountMe(ctx context.Context, req *sharepb.Empty) (*userpb.AccountInfoResponse, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)

	orgs, _ := s.OrgClient.GetOrganizationByUserId(ctx, &orgpb.GetOrganizationByUserIdRequest{
		UserId: uint32(profileId),
	})

	groups, _ := s.GroupClient.GetGroupByUserId(ctx, &organizationpb.GetGroupByUserIdRequest{
		UserId: uint32(profileId),
	})

	profile, _ := s.ProfileRepo.GetByProfileID(profileId)
	if profile == nil {
		return nil, status.Errorf(codes.NotFound, "profile not found")
	}

	accountInfos := []*userpb.AccountInfo{}
	ownerType := uint32(10)
	accountInfos = append(accountInfos, &userpb.AccountInfo{
		// Id:        uint64(profile.ProfileID),
		FullName:  &profile.FullName,
		Avatar:    &profile.Avatar,
		Phone:     &profile.Phone,
		Email:     &profile.Email,
		ProfileId: &profileId,
		OwnerType: &ownerType,
	})

	if orgs != nil && len(orgs.Data) > 0 {
		for _, org := range orgs.Data {
			ownerType := uint32(10)
			accountInfos = append(accountInfos, &userpb.AccountInfo{
				// Id:              uint64(org.Id),
				FullName:        &org.Name,
				Avatar:          org.LogoUrl,
				Phone:           org.Phone,
				Email:           org.Email,
				BackgroundImage: org.LogoUrl,
				OrganizationId:  &org.Id,
				OwnerType:       &ownerType,
			})
		}
	}

	if groups != nil && len(groups.Data) > 0 {
		for _, group := range groups.Data {
			accountInfos = append(accountInfos, &userpb.AccountInfo{
				// Id:       uint64(group.Id),
				FullName: &group.Name,
				Avatar:   &group.AvatarUrl,
				GroupId:  &group.Id,
			})
		}
	}

	return &userpb.AccountInfoResponse{
		Data:  accountInfos,
		Total: uint64(len(accountInfos)),
	}, nil
}
