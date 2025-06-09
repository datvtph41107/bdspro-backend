package services

import (
	"context"
	"fmt"
	pb_user "pb/types/user"
	"time"
	"user/repo"

	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProfileService struct {
	pb_user.UnimplementedProfileServiceServer
	ProfileRepo *repo.ProfileRepo
	cache       *cache.Cache
}

func NewGrpcProfileService(profileRepo *repo.ProfileRepo) *ProfileService {
	return &ProfileService{
		ProfileRepo: profileRepo,
		cache:       cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (s *ProfileService) GetProfileByIds(ctx context.Context, req *pb_user.GetProfileByIdsRequest) (*pb_user.GetProfileByIdsResponse, error) {
	var profilesPb []*pb_user.ProfileItem
	var missingIDs []uint64

	// 1. Check cache trước
	for _, id := range req.Ids {
		if cached, found := s.cache.Get(fmt.Sprintf("%d", id)); found {
			profilesPb = append(profilesPb, cached.(*pb_user.ProfileItem))
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
			pb := &pb_user.ProfileItem{
				Id:       profile.ProfileID,
				FullName: profile.FullName,
				Avatar:   profile.Avatar,
			}
			// Cache lại
			s.cache.Set(fmt.Sprintf("%d", profile.ProfileID), pb, cache.DefaultExpiration)
			profilesPb = append(profilesPb, pb)
		}
	}

	return &pb_user.GetProfileByIdsResponse{
		Profiles: profilesPb,
	}, nil
}
