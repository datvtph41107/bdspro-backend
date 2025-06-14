package services

import (
	"context"
	"fmt"
	sharepb "pb/types/shared"
	pb_user "pb/types/user"
	"time"
	"user/repo"
	user_service "user/services/user"

	"github.com/jinzhu/copier"
	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcProfileService struct {
	pb_user.UnimplementedProfileServiceServer
	ProfileRepo *repo.ProfileRepo
	Profile     *user_service.ProfileService
	cache       *cache.Cache
}

func NewGrpcProfileService(profileRepo *repo.ProfileRepo, profileService *user_service.ProfileService) *GrpcProfileService {
	return &GrpcProfileService{
		ProfileRepo: profileRepo,
		Profile:     profileService,
		cache:       cache.New(5*time.Minute, 10*time.Minute),
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

// @Summary Lấy thông tin công khai của người dùng
// @Description Lấy thông tin công khai của người dùng
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Profile"
// @Success 200 {object} pb_user.ProfilePublicInfo "Thông tin công khai của người dùng"
// @Router /v2/user/profile/info/{id} [get]
func (s *GrpcProfileService) GetProfileInfo(ctx context.Context, req *pb_user.IdRequest) (*pb_user.ProfilePublicInfo, error) {
	result, err := s.Profile.PublicUserInfo(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get profile info: %v", err)
	}

	dto := &pb_user.ProfilePublicInfo{}
	copier.Copy(dto, result)
	return dto, nil
}
