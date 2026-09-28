package services

import (
	"context"
	"time"

	_utils "common/utils"
	sharepb "pb/types/shared"
	"userv4/infra/client"
	"userv4/infra/mapper"
	usecases "userv4/internal/usecase/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProfileUsecaseInfoService adapts legacy ProfileUsecase into modules/user ProfileService.
// Profile ownership lives here; ProfileUsecase / KYC / hub remain backends until
// those lookups move into this package.
type ProfileUsecaseInfoService struct {
	profile       *usecases.ProfileUsecase
	kyc           usecases.IKYCUsecase
	profileMapper *mapper.ProfileMapper
	kycMapper     *mapper.KYCMapper
	hub           *client.HubClient
	sync          *_utils.SyncUtil
}

func NewProfileUsecaseInfoService(
	profile *usecases.ProfileUsecase,
	kyc usecases.IKYCUsecase,
	profileMapper *mapper.ProfileMapper,
	kycMapper *mapper.KYCMapper,
	hub *client.HubClient,
	sync *_utils.SyncUtil,
) *ProfileUsecaseInfoService {
	return &ProfileUsecaseInfoService{
		profile:       profile,
		kyc:           kyc,
		profileMapper: profileMapper,
		kycMapper:     kycMapper,
		hub:           hub,
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
