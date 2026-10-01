package handler

import (
	_provider "common/domain/provider"
	_redis "common/redis"
	_utils "common/utils"
	"context"
	"fmt"
	"log"
	"time"

	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"userv4/infra/client"
	"userv4/infra/mapper"
	"userv4/internal/interface/repo"
	"userv4/internal/job"
	organizationusecase "userv4/internal/usecase/organization"
	usecases "userv4/internal/usecase/user"
	usercontroller "userv4/modules/user/controller"

	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcProfileService struct {
	userpb.UnimplementedProfileServiceServer
	ProfileRepo           repo.IProfileRepo
	AdminUsecase          *usecases.AdminUsecase
	ProfileUsecase        *usecases.ProfileUsecase
	KYCUsecase            usecases.IKYCUsecase
	cache                 *cache.Cache
	OrganizationService   *organizationusecase.Service
	ProfileMapper         *mapper.ProfileMapper
	KYCMapper             *mapper.KYCMapper
	CrmClient             _provider.CrmProvider
	HubClient             *client.HubClient
	BdsproClient          *client.BdsproClient
	UserDashboardStatsJob *job.UserDashboardStatsJob
	SyncProvider          *_utils.SyncUtil
	Redis                 *_redis.RedisService
	ProfileController     *usercontroller.ProfileController
}

func NewGrpcProfileService(
	profileRepo repo.IProfileRepo,
	adminUsecase *usecases.AdminUsecase,
	profileUsecase *usecases.ProfileUsecase,
	kycUsecase usecases.IKYCUsecase,
	organizationService *organizationusecase.Service,
	profileMapper *mapper.ProfileMapper,
	kycMapper *mapper.KYCMapper,
	crmClient _provider.CrmProvider,
	hubClient *client.HubClient,
	bdsproClient *client.BdsproClient,
	userDashboardStatsJob *job.UserDashboardStatsJob,
	SyncProvider *_utils.SyncUtil,
	redis *_redis.RedisService,
	profileController *usercontroller.ProfileController,
) *GrpcProfileService {
	return &GrpcProfileService{
		ProfileRepo:           profileRepo,
		ProfileUsecase:        profileUsecase,
		AdminUsecase:          adminUsecase,
		KYCUsecase:            kycUsecase,
		cache:                 cache.New(5*time.Minute, 10*time.Minute),
		OrganizationService:   organizationService,
		ProfileMapper:         profileMapper,
		KYCMapper:             kycMapper,
		CrmClient:             crmClient,
		HubClient:             hubClient,
		BdsproClient:          bdsproClient,
		UserDashboardStatsJob: userDashboardStatsJob,
		SyncProvider:          SyncProvider,
		Redis:                 redis,
		ProfileController:     profileController,
	}
}

func (s *GrpcProfileService) GetProfileByIds(ctx context.Context, req *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error) {
	var profilesPb []*sharepb.ProfileItem
	var missingIDs []uint64

	for _, id := range req.Ids {
		if cached, found := s.cache.Get(fmt.Sprintf("%d", id)); found {
			profilesPb = append(profilesPb, cached.(*sharepb.ProfileItem))
		} else {
			missingIDs = append(missingIDs, id)
		}
	}

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
			s.cache.Set(fmt.Sprintf("%d", profile.ProfileID), pb, cache.DefaultExpiration)
			profilesPb = append(profilesPb, pb)
		}
	}

	return &sharepb.GetProfileByIdsResponse{Profiles: profilesPb}, nil
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

	return &sharepb.GetProfileByIdsResponse{Profiles: profilesPb}, nil
}

func (s *GrpcProfileService) GetProfileInfo(ctx context.Context, req *userpb.IdRequest) (*userpb.ProfilePublicInfo, error) {
	return s.ProfileController.GetProfileInfo(ctx, req.Id)
}

func (s *GrpcProfileService) GetProfileInfoV3(ctx context.Context, req *sharepb.SyncRequest) (*sharepb.TotalResponseProto, error) {
	return s.ProfileController.GetProfileInfoV3(ctx, req.Id, req.Access, req.Timestamp)
}

func (s *GrpcProfileService) GetProfileInfoMe(ctx context.Context, req *sharepb.Empty) (*sharepb.UserV3Proto, error) {
	return s.ProfileController.GetInfoMe(ctx)
}

func (s *GrpcProfileService) GenTileSessionToken(ctx context.Context, _ *sharepb.Empty) (*userpb.GenTileSessionTokenResponse, error) {
	result, err := s.ProfileController.GenTileSessionToken(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("gen tile session: %s, %s, %v", result.SessionK, result.SessionEncryptKey, result.ExpiresIn)
	return &userpb.GenTileSessionTokenResponse{
		SessionK:          result.SessionK,
		SessionEncryptKey: result.SessionEncryptKey,
		ExpiresIn:         result.ExpiresIn,
	}, nil
}

func (s *GrpcProfileService) ListAccountMe(ctx context.Context, req *sharepb.Empty) (*userpb.AccountInfoResponse, error) {
	return s.ProfileController.ListAccountMe(ctx)
}

func (s *GrpcProfileService) SetInfo(ctx context.Context, req *userpb.SetInfoRequest) (*sharepb.Empty, error) {
	if err := s.ProfileController.SetInfo(ctx, req); err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}

func (s *GrpcProfileService) SetPerson(ctx context.Context, req *sharepb.UserV3Proto) (*sharepb.Empty, error) {
	return &sharepb.Empty{}, s.ProfileController.SetPerson(ctx, req)
}

func (s *GrpcProfileService) UpdatePrivacySetting(ctx context.Context, req *userpb.UpdatePrivacySettingRequest) (*sharepb.Empty, error) {
	if err := s.ProfileController.UpdatePrivacySetting(ctx, req); err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}

func (s *GrpcProfileService) QuickSetup(ctx context.Context, req *sharepb.UserV3Proto) (*sharepb.Empty, error) {
	if err := s.ProfileController.QuickSetup(ctx, req); err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}

func (s *GrpcProfileService) SearchProfilePublic(ctx context.Context, req *userpb.SearchProfilePublicRequest) (*userpb.SearchProfilePublicResponse, error) {
	return s.ProfileController.SearchProfilePublic(ctx, req.Text, req.Page, req.Size)
}

func (s *GrpcProfileService) CheckAccountStatus(ctx context.Context, req *userpb.CheckAccountStatusRequest) (*userpb.CheckAccountStatusResponse, error) {
	result, err := s.ProfileUsecase.CheckAccountStatus(ctx, req.ProfileId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check account status: %v", err)
	}
	return result, nil
}

func (s *GrpcProfileService) FindByPhone(ctx context.Context, req *userpb.FindByPhoneRequest) (*userpb.FindByPhoneResponse, error) {
	return s.ProfileController.FindByPhone(ctx, req.Phone)
}

func (s *GrpcProfileService) GetProfileCompletion(ctx context.Context, req *sharepb.Empty) (*userpb.ProfileCompletionResponse, error) {
	return s.ProfileController.GetProfileCompletion(ctx)
}

func (s *GrpcProfileService) GetStatusOnline(ctx context.Context, req *userpb.IdRequest) (*userpb.StatusOnlineResponse, error) {
	return s.ProfileController.GetStatusOnline(ctx, req.Id)
}

func (s *GrpcProfileService) SetDemoMode(ctx context.Context, req *userpb.DemoModeRequest) (*sharepb.Empty, error) {
	if err := s.ProfileController.SetDemoMode(ctx, req.GetEnable()); err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}
