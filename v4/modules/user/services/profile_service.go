package services

import (
	"context"
	"fmt"
	"time"

	_dto "common/domain/dto"
	_provider "common/domain/provider"
	_errors "common/errors"
	_utils "common/utils"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"userv4/infra/client"
	"userv4/infra/mapper"
	"userv4/internal/interface/repo"
	organizationusecase "userv4/internal/usecase/organization"
	usecases "userv4/internal/usecase/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProfileService owns profile HTTP/gRPC surface migrated into modules/user.
type ProfileService interface {
	GetInfoMe(ctx context.Context) (*sharepb.UserV3Proto, error)
	GetProfileInfoV3(ctx context.Context, id uint64, access string, timestamp int64) (*sharepb.TotalResponseProto, error)
	GetProfileInfo(ctx context.Context, id uint64) (*userpb.ProfilePublicInfo, error)
	ListAccountMe(ctx context.Context) (*userpb.AccountInfoResponse, error)
	SetInfo(ctx context.Context, req *userpb.SetInfoRequest) error
	SetPerson(ctx context.Context, req *sharepb.UserV3Proto) error
	UpdatePrivacySetting(ctx context.Context, req *userpb.UpdatePrivacySettingRequest) error
	QuickSetup(ctx context.Context, req *sharepb.UserV3Proto) error
	SearchProfilePublic(ctx context.Context, text string, page, size uint32) (*userpb.SearchProfilePublicResponse, error)
	FindByPhone(ctx context.Context, phone string) (*userpb.FindByPhoneResponse, error)
	GetProfileCompletion(ctx context.Context) (*userpb.ProfileCompletionResponse, error)
	GetStatusOnline(ctx context.Context, id uint64) (*userpb.StatusOnlineResponse, error)
}

// ProfileUsecaseInfoService adapts legacy ProfileUsecase into modules/user ProfileService.
type ProfileUsecaseInfoService struct {
	profile       *usecases.ProfileUsecase
	profileRepo   repo.IProfileRepo
	kyc           usecases.IKYCUsecase
	profileMapper *mapper.ProfileMapper
	kycMapper     *mapper.KYCMapper
	hub           *client.HubClient
	crm           _provider.CrmProvider
	bdspro        *client.BdsproClient
	sync          *_utils.SyncUtil
	org           *organizationusecase.Service
}

func NewProfileUsecaseInfoService(
	profile *usecases.ProfileUsecase,
	kyc usecases.IKYCUsecase,
	profileMapper *mapper.ProfileMapper,
	kycMapper *mapper.KYCMapper,
	hub *client.HubClient,
	crm _provider.CrmProvider,
	bdspro *client.BdsproClient,
	sync *_utils.SyncUtil,
) *ProfileUsecaseInfoService {
	return &ProfileUsecaseInfoService{
		profile:       profile,
		kyc:           kyc,
		profileMapper: profileMapper,
		kycMapper:     kycMapper,
		hub:           hub,
		crm:           crm,
		bdspro:        bdspro,
		sync:          sync,
	}
}

// WithProfileRepo attaches optional direct profile store access.
func (s *ProfileUsecaseInfoService) WithProfileRepo(profileRepo repo.IProfileRepo) *ProfileUsecaseInfoService {
	if s != nil {
		s.profileRepo = profileRepo
	}
	return s
}

// WithOrganization attaches optional organization listing for profile views.
func (s *ProfileUsecaseInfoService) WithOrganization(org *organizationusecase.Service) *ProfileUsecaseInfoService {
	if s != nil {
		s.org = org
	}
	return s
}

func (s *ProfileUsecaseInfoService) GetInfoMe(ctx context.Context) (*sharepb.UserV3Proto, error) {
	result, err := s.profile.GetInfoMe(ctx)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		return nil, status.Errorf(codes.Internal, "failed to get profile info: %v", err)
	}

	var kyc *sharepb.KYCV3Proto
	if s.kyc != nil {
		kycEntity, errKYC := s.kyc.GetMyKYC(ctx)
		if errKYC == nil && kycEntity != nil && s.kycMapper != nil {
			kyc = s.kycMapper.ToProtoWithOwnerCheck(kycEntity, true)
		}
	}

	profilePb := s.profileMapper.ModelToUserV3Proto(result, kyc)

	if s.hub != nil && (result.ProvinceID != nil || result.WardID != nil) {
		provinceName, wardName := s.hub.GetProvinceAndWardNames(ctx, result.ProvinceID, result.WardID)
		if profilePb.Address != nil {
			if provinceName != "" {
				profilePb.Address.ProvinceName = provinceName
			}
			if wardName != "" {
				profilePb.Address.WardName = wardName
			}
		}
	}

	now := time.Now().UnixMilli()
	if s.sync != nil {
		_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyProfileInfoMe, 0), now)
	}
	profilePb.Timestamp = now

	return profilePb, nil
}

func (s *ProfileUsecaseInfoService) GetProfileInfoV3(ctx context.Context, id uint64, access string, timestamp int64) (*sharepb.TotalResponseProto, error) {
	syncKey := fmt.Sprintf("time:profile:%d", id)
	if s.sync != nil {
		updated := s.sync.HasUpdated(ctx, syncKey, timestamp)
		if !updated {
			return &sharepb.TotalResponseProto{}, nil
		}
	}

	result, err := s.profile.UserInfoV3(ctx, id, access)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		return nil, status.Errorf(codes.Internal, "failed to get profile info: %v", err)
	}

	if s.sync != nil && result.User != nil {
		_ = s.sync.PutTimestamp(ctx, syncKey, result.User.UpdatedAt.UnixMilli())
	}

	if s.hub != nil && result.User != nil && (result.User.ProvinceID != nil || result.User.WardID != nil) {
		provinceName, wardName := s.hub.GetProvinceAndWardNames(ctx, result.User.ProvinceID, result.User.WardID)
		if result.User.Address != nil {
			if provinceName != "" {
				result.User.Address.ProvinceName = provinceName
			}
			if wardName != "" {
				result.User.Address.WardName = wardName
			}
		}
	}

	var contactInfo *sharepb.ContactInfoV3Proto
	if s.crm != nil {
		friendInfoResp, errFriendInfo := s.crm.GetFriendInfo(ctx, id)
		if errFriendInfo == nil && friendInfoResp != nil {
			contactInfo = &sharepb.ContactInfoV3Proto{
				NumFriend:    friendInfoResp.NumFriend,
				NumFollower:  friendInfoResp.NumFollower,
				NumFollowing: friendInfoResp.NumFollowing,
				ContactId:    friendInfoResp.ContactId,
				BlockId:      friendInfoResp.BlockId,
				ReportId:     friendInfoResp.ReportId,
				FollowId:     friendInfoResp.FollowId,
			}
		}
	}

	var bdsproInfo *sharepb.BDSProDashboardProto
	if s.bdspro != nil {
		bdsproInfoResp, errBdsPro := s.bdspro.GetDashboardByProfileId(ctx, id)
		if errBdsPro == nil && bdsproInfoResp != nil {
			bdsproInfo = bdsproInfoResp
		}
	}

	var kycProto *sharepb.KYCV3Proto
	requestProfileID := _utils.GetProfileIdWithContext(ctx)
	isOwner := requestProfileID > 0 && requestProfileID == id
	if result.KYC != nil {
		kycProto = s.profileMapper.DTOToKYCProtoWithOwnerCheck(result.KYC, isOwner)
	} else if s.kyc != nil && s.kycMapper != nil {
		kycEntity, errKYC := s.kyc.GetKYCByProfileID(ctx, id)
		if errKYC == nil && kycEntity != nil {
			kycProto = s.kycMapper.ToProtoWithOwnerCheck(kycEntity, isOwner)
		}
	}

	return &sharepb.TotalResponseProto{
		User:         s.profileMapper.DTOToUserV3Proto(result.User, nil),
		FriendStatus: s.profileMapper.DTOToFriendV3Proto(result.FriendStatus),
		RateStats:    s.profileMapper.DTOToRateStatsV3Proto(result.RateStats),
		Kyc:          kycProto,
		ContactInfo:  contactInfo,
		BdsproInfo:   bdsproInfo,
	}, nil
}

func (s *ProfileUsecaseInfoService) GetProfileInfo(ctx context.Context, id uint64) (*userpb.ProfilePublicInfo, error) {
	result, err := s.profile.PublicUserInfo(ctx, id, false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get profile info: %v", err)
	}

	var kyc *sharepb.KYCV3Proto
	if s.kyc != nil {
		kycEntity, errKYC := s.kyc.GetKYCByProfileID(ctx, id)
		if errKYC == nil && kycEntity != nil && s.kycMapper != nil {
			requestProfileID := _utils.GetProfileIdWithContext(ctx)
			isOwner := requestProfileID > 0 && requestProfileID == id
			kyc = s.kycMapper.ToProtoWithOwnerCheck(kycEntity, isOwner)
		}
	}

	profilePb := s.profileMapper.PublicUserInfoToPb(result, kyc)
	if s.hub != nil && (result.ProvinceID != nil || result.WardID != nil) {
		provinceName, wardName := s.hub.GetProvinceAndWardNames(ctx, result.ProvinceID, result.WardID)
		if provinceName != "" {
			profilePb.ProvinceName = &provinceName
		}
		if wardName != "" {
			profilePb.WardName = &wardName
		}
	}
	return profilePb, nil
}

func (s *ProfileUsecaseInfoService) ListAccountMe(ctx context.Context) (*userpb.AccountInfoResponse, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	var orgList []struct {
		ID    uint64
		Name  string
		Phone string
		Email string
	}
	if s.org != nil {
		orgs, err := s.org.ListForProfile(ctx, profileId)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "list profile organizations: %v", err)
		}
		for _, org := range orgs {
			orgList = append(orgList, struct {
				ID    uint64
				Name  string
				Phone string
				Email string
			}{ID: org.ID, Name: org.Name, Phone: org.Phone, Email: org.Email})
		}
	}

	if s.profileRepo == nil {
		return nil, status.Errorf(codes.Internal, "profile repository not configured")
	}
	profile, _ := s.profileRepo.GetByProfileID(profileId)
	if profile == nil {
		return nil, status.Errorf(codes.NotFound, "profile not found")
	}

	accountInfos := []*userpb.AccountInfo{}
	ownerType := uint32(10)
	accountInfos = append(accountInfos, &userpb.AccountInfo{
		FullName:  &profile.FullName,
		Avatar:    &profile.Avatar,
		Phone:     &profile.Phone,
		Email:     &profile.Email,
		ProfileId: &profileId,
		OwnerType: &ownerType,
	})

	for _, org := range orgList {
		ot := uint32(30)
		name, phone, email, id := org.Name, org.Phone, org.Email, org.ID
		accountInfos = append(accountInfos, &userpb.AccountInfo{
			FullName:       &name,
			Phone:          &phone,
			Email:          &email,
			OwnerType:      &ot,
			OrganizationId: &id,
		})
	}

	return &userpb.AccountInfoResponse{
		Data:  accountInfos,
		Total: int32(len(accountInfos)),
	}, nil
}

func (s *ProfileUsecaseInfoService) SetInfo(ctx context.Context, req *userpb.SetInfoRequest) error {
	if req == nil {
		return _errors.ReturnError(400, "request is required")
	}
	dto := s.profileMapper.SetProfileToDTO(req)
	if _, err := s.profile.UpdateUserInfo(ctx, dto); err != nil {
		return err
	}
	if s.sync != nil {
		_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyProfileInfoMe, 0), time.Now().UnixMilli())
	}
	return nil
}

func (s *ProfileUsecaseInfoService) SetPerson(ctx context.Context, req *sharepb.UserV3Proto) error {
	dto := s.profileMapper.UserV3ToDTO(req)
	return s.profile.SetPerson(ctx, dto)
}

func (s *ProfileUsecaseInfoService) UpdatePrivacySetting(ctx context.Context, req *userpb.UpdatePrivacySettingRequest) error {
	dto := s.profileMapper.PrivacyRequestToDTO(req)
	return s.profile.UpdatePrivacySetting(ctx, dto)
}

func (s *ProfileUsecaseInfoService) QuickSetup(ctx context.Context, req *sharepb.UserV3Proto) error {
	dto := s.profileMapper.UserV3ToDTO(req)
	return s.profile.QuickSetup(ctx, dto)
}

func (s *ProfileUsecaseInfoService) SearchProfilePublic(ctx context.Context, text string, page, size uint32) (*userpb.SearchProfilePublicResponse, error) {
	pageable := _dto.Pagable{Page: page, Size: size}
	results, total, err := s.profile.SearchProfilePublic(ctx, &userpb.SearchProfilePublicRequest{
		Text: text,
		Page: page,
		Size: size,
	}, pageable)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search profile: %v", err)
	}
	return &userpb.SearchProfilePublicResponse{
		Data:  s.profileMapper.UserSearchItemToPbs(results),
		Total: int32(total),
	}, nil
}

func (s *ProfileUsecaseInfoService) FindByPhone(ctx context.Context, phone string) (*userpb.FindByPhoneResponse, error) {
	if phone == "" {
		return nil, _errors.ReturnError(400, "Số điện thoại không đúng")
	}
	profile, err := s.profile.FindByPhone(ctx, phone)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find user by phone: %v", err)
	}
	if profile == nil {
		return nil, _errors.ReturnError(404, "Số điện thoại không tồn tại")
	}
	return &userpb.FindByPhoneResponse{
		FullName: profile.FullName,
		Phone:    profile.Phone,
		Avatar:   profile.Avatar,
	}, nil
}

func (s *ProfileUsecaseInfoService) GetProfileCompletion(ctx context.Context) (*userpb.ProfileCompletionResponse, error) {
	result, err := s.profile.GetProfileCompletion(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get profile completion: %v", err)
	}
	return &userpb.ProfileCompletionResponse{
		Percentage:      uint32(result.Percentage),
		CompletedFields: result.CompletedFields,
		TotalFields:     result.TotalFields,
		MissingFields:   result.MissingFields,
	}, nil
}

func (s *ProfileUsecaseInfoService) GetStatusOnline(ctx context.Context, id uint64) (*userpb.StatusOnlineResponse, error) {
	result, err := s.profile.GetStatusOnline(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get status online: %v", err)
	}
	return result, nil
}
