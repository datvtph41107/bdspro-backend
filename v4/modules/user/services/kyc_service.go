package services

import (
	"context"

	_dto "common/domain/dto"
	_errors "common/errors"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	mapper "userv4/modules/user/mapper"
	"userv4/modules/user/dto"
	usecases "userv4/modules/user/services/user"

	"google.golang.org/grpc/status"
)

// KYCService owns KYC HTTP surface migrated into modules/user.
type KYCService interface {
	SubmitKYC(ctx context.Context, req *userpb.SubmitKYCRequest) (*sharepb.KYCV3Proto, error)
	ListKYCs(ctx context.Context, req *userpb.ListKYCsRequest) (*userpb.ListKYCsResponse, error)
	ApproveKYC(ctx context.Context, id uint64) error
	RejectKYC(ctx context.Context, id uint64, reason string) error
}

// LocalKYCService adapts KYC usecases into modules/user.
type LocalKYCService struct {
	kyc      usecases.IKYCUsecase
	adminKYC usecases.IAdminKYCUsecase
	mapper   *mapper.KYCMapper
}

func NewLocalKYCService(
	kyc usecases.IKYCUsecase,
	adminKYC usecases.IAdminKYCUsecase,
	mapper *mapper.KYCMapper,
) *LocalKYCService {
	return &LocalKYCService{kyc: kyc, adminKYC: adminKYC, mapper: mapper}
}

var _ KYCService = (*LocalKYCService)(nil)

func (s *LocalKYCService) SubmitKYC(ctx context.Context, req *userpb.SubmitKYCRequest) (*sharepb.KYCV3Proto, error) {
	submitDTO := dto.SubmitKYCRequest{
		FullName:     req.FullName,
		IdentityCard: req.IdentityCard,
		FrontImage:   req.FrontImage,
		BackImage:    req.BackImage,
		SelfieImage:  req.SelfieImage,
	}
	if req.IdNumber != nil {
		submitDTO.IDNumber = *req.IdNumber
	}
	if req.DateOfBirth != nil {
		submitDTO.DateOfBirth = *req.DateOfBirth
	}
	if req.ExpiryDate != nil {
		submitDTO.ExpiryDate = *req.ExpiryDate
	}

	kyc, err := s.kyc.SubmitKYC(ctx, submitDTO)
	if err != nil {
		return nil, status.Error(_errors.RPCCodeForHTTP(int32(err.Code)), err.Message)
	}
	return s.mapper.ToProtoWithOwnerCheck(kyc, true), nil
}

func (s *LocalKYCService) ListKYCs(ctx context.Context, req *userpb.ListKYCsRequest) (*userpb.ListKYCsResponse, error) {
	listDTO := dto.ListKYCsRequest{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		Search:    req.Search,
		Status:    req.Status,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}
	entities, total, err := s.adminKYC.ListKYCs(ctx, listDTO)
	if err != nil {
		return nil, status.Error(_errors.RPCCodeForHTTP(int32(err.Code)), err.Message)
	}
	return &userpb.ListKYCsResponse{
		Data:  s.mapper.ToProtoList(entities),
		Total: total,
	}, nil
}

func (s *LocalKYCService) ApproveKYC(ctx context.Context, id uint64) error {
	if err := s.adminKYC.ApproveKYC(ctx, id); err != nil {
		return status.Error(_errors.RPCCodeForHTTP(int32(err.Code)), err.Message)
	}
	return nil
}

func (s *LocalKYCService) RejectKYC(ctx context.Context, id uint64, reason string) error {
	if err := s.adminKYC.RejectKYC(ctx, id, reason); err != nil {
		return status.Error(_errors.RPCCodeForHTTP(int32(err.Code)), err.Message)
	}
	return nil
}
