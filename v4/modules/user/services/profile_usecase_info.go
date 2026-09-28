package services

import (
	"context"
	"fmt"
	"time"

	_provider "common/domain/provider"
	_utils "common/utils"
	sharepb "pb/types/shared"
	"userv4/infra/client"
	"userv4/infra/mapper"
	usecases "userv4/internal/usecase/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProfileUsecaseInfoService adapts legacy ProfileUsecase into modules/user ProfileService.
// Profile ownership lives here; ProfileUsecase / KYC / hub / CRM / BDSPro remain backends until
// those lookups move into this package. Reads use the v4 process DB (database.dsn), not db_user.
type ProfileUsecaseInfoService struct {
	profile       *usecases.ProfileUsecase
	kyc           usecases.IKYCUsecase
	profileMapper *mapper.ProfileMapper
	kycMapper     *mapper.KYCMapper
	hub           *client.HubClient
	crm           _provider.CrmProvider
	bdspro        *client.BdsproClient
	sync          *_utils.SyncUtil
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
			// GetInfoMe is always an owner request.
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

	// Province / ward names from hub (in-process on v4 DSN).
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

	// Friend / contact counts from CRM.
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

	// BDSPro dashboard stats.
	var bdsproInfo *sharepb.BDSProDashboardProto
	if s.bdspro != nil {
		bdsproInfoResp, errBdsPro := s.bdspro.GetDashboardByProfileId(ctx, id)
		if errBdsPro == nil && bdsproInfoResp != nil {
			bdsproInfo = bdsproInfoResp
		}
	}

	// KYC with owner check.
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
