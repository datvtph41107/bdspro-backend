package services

import (
	"context"
	"errors"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"

	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/entity"
	"userv4/modules/bdspro/enums"
	"userv4/modules/bdspro/repository"

	"gorm.io/gorm"
)

// DealInvitationService owns deal invitation send/accept/reject and list.
type DealInvitationService interface {
	Send(ctx context.Context, req dto.SendDealInvitationRequest) (*bdspropb.SendDealInvitationResponse, error)
	Accept(ctx context.Context, req dto.AcceptDealInvitationRequest) (*bdspropb.AcceptDealInvitationResponse, error)
	Reject(ctx context.Context, req dto.RejectDealInvitationRequest) (*bdspropb.RejectDealInvitationResponse, error)
	List(ctx context.Context, req dto.DealInvitationListRequest) (*bdspropb.GetDealInvitationsByDealResponse, error)
	Remove(ctx context.Context, req dto.RemoveFromDealRequest) (*bdspropb.RemoveFromDealResponse, error)
	UpdateMemberRole(ctx context.Context, req dto.UpdateDealMemberRoleRequest) (*bdspropb.UpdateDealMemberRoleResponse, error)
}

type LocalDealInvitationService struct {
	sync    *_utils.SyncUtil
	deal    repository.DealRepository
	members repository.DealMemberRepository
	deals   DealService
}

func NewLocalDealInvitationService(
	sync *_utils.SyncUtil,
	deal repository.DealRepository,
	members repository.DealMemberRepository,
	deals DealService,
) *LocalDealInvitationService {
	return &LocalDealInvitationService{sync: sync, deal: deal, members: members, deals: deals}
}

func (s *LocalDealInvitationService) List(ctx context.Context, req dto.DealInvitationListRequest) (*bdspropb.GetDealInvitationsByDealResponse, error) {
	if s.deals == nil {
		return nil, _errors.InternalServerException("deal service unavailable")
	}
	return s.deals.GetInvitations(ctx, dto.DealListByIDRequest{
		DealID:    req.DealID,
		Page:      req.Page,
		Size:      req.Size,
		Timestamp: req.Timestamp,
	})
}

func (s *LocalDealInvitationService) Send(ctx context.Context, req dto.SendDealInvitationRequest) (*bdspropb.SendDealInvitationResponse, error) {
	if req.DealID == 0 {
		return nil, _errors.BadRequestException("deal id is required")
	}
	if req.MemberID == 0 {
		return nil, _errors.BadRequestException("member id is required")
	}
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, _errors.ForbiddenException("profile is required")
	}
	if req.MemberID == profileID {
		return nil, _errors.BadRequestException("cannot invite yourself")
	}

	deal, err := s.deal.GetByID(ctx, req.DealID)
	if err != nil {
		return nil, err
	}
	if deal == nil {
		return nil, _errors.NotFoundException("deal not found")
	}

	exists, err := s.members.ExistsByDealIDAndMemberID(ctx, req.DealID, req.MemberID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, _errors.ReturnError(400, "invitation already exists for this user")
	}

	roleKey := req.RoleKey
	if roleKey == 0 {
		roleKey = uint32(enums.RoleKeyDealMember)
	}
	now := time.Now()
	row := &entity.DealMember{
		DealID:     req.DealID,
		MemberID:   req.MemberID,
		InviterID:  profileID,
		Message:    req.Message,
		RoleID:     req.RoleID,
		RoleKey:    roleKey,
		ColorID:    req.ColorID,
		Status:     uint32(enums.DealMemberStatusInvited),
		MemberType: uint32(enums.DealMemberTypeMember),
		InvitedAt:  &now,
	}
	created, err := s.members.Create(ctx, row)
	if err != nil {
		return nil, err
	}
	s.bumpInvitationSync(ctx, req.DealID)
	s.clearDealMeSync(ctx, req.MemberID)

	return &bdspropb.SendDealInvitationResponse{
		InvitationId: created.ID,
		Success:      true,
		Message:      req.Message,
		RoleId:       req.RoleID,
		RoleKey:      roleKey,
	}, nil
}

func (s *LocalDealInvitationService) Accept(ctx context.Context, req dto.AcceptDealInvitationRequest) (*bdspropb.AcceptDealInvitationResponse, error) {
	row, dealID, err := s.loadInvitationForDeal(ctx, req.DealID, req.InvitationID)
	if err != nil {
		return nil, err
	}
	if row.Status != uint32(enums.DealMemberStatusInvited) {
		return nil, _errors.BadRequestException("invitation is not in pending status")
	}
	now := time.Now()
	row.Status = uint32(enums.DealMemberStatusAccepted)
	row.RespondedAt = &now
	if err := s.members.Update(ctx, row); err != nil {
		return nil, err
	}
	s.bumpInvitationSync(ctx, dealID)
	s.clearDealMeSync(ctx, row.MemberID)
	return &bdspropb.AcceptDealInvitationResponse{
		InvitationId: req.InvitationID,
		Success:      true,
	}, nil
}

func (s *LocalDealInvitationService) UpdateMemberRole(ctx context.Context, req dto.UpdateDealMemberRoleRequest) (*bdspropb.UpdateDealMemberRoleResponse, error) {
	if req.DealID == 0 || req.MemberID == 0 {
		return nil, _errors.BadRequestException("deal id and member id are required")
	}
	if req.RoleKey == 0 {
		return nil, _errors.BadRequestException("role key is required")
	}
	if _utils.GetProfileIdWithContext(ctx) == 0 {
		return nil, _errors.ForbiddenException("profile is required")
	}

	err := s.members.UpdateRole(ctx, req.DealID, req.MemberID, req.RoleKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, _errors.NotFoundException("deal member not found")
		}
		return nil, err
	}
	s.bumpInvitationSync(ctx, req.DealID)

	const msg = "Update deal member role successfully"
	return &bdspropb.UpdateDealMemberRoleResponse{
		Success:  true,
		Message:  msg,
		DealId:   req.DealID,
		MemberId: req.MemberID,
		RoleKey:  req.RoleKey,
	}, nil
}

func (s *LocalDealInvitationService) Remove(ctx context.Context, req dto.RemoveFromDealRequest) (*bdspropb.RemoveFromDealResponse, error) {
	if req.InvitationID == 0 {
		return nil, _errors.BadRequestException("invitation id is required")
	}
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, _errors.ForbiddenException("profile is required")
	}

	row, err := s.members.GetByID(ctx, req.InvitationID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, _errors.NotFoundException("invitation not found")
	}
	if row.RoleKey == uint32(enums.RoleKeyDealOwner) || row.IsOwner {
		return nil, _errors.BadRequestException("cannot remove deal owner")
	}
	if row.Status != uint32(enums.DealMemberStatusInvited) && row.Status != uint32(enums.DealMemberStatusAccepted) {
		return nil, _errors.BadRequestException("member is not active in this deal")
	}

	now := time.Now()
	row.Status = uint32(enums.DealMemberStatusWithdrawn)
	row.WithdrawnAt = &now
	row.Message = req.Reason
	if err := s.members.Update(ctx, row); err != nil {
		return nil, err
	}
	s.bumpInvitationSync(ctx, row.DealID)
	s.clearDealMeSync(ctx, row.MemberID)

	return &bdspropb.RemoveFromDealResponse{
		InvitationId: req.InvitationID,
		Success:      true,
	}, nil
}

func (s *LocalDealInvitationService) Reject(ctx context.Context, req dto.RejectDealInvitationRequest) (*bdspropb.RejectDealInvitationResponse, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	row, dealID, err := s.loadInvitationForDeal(ctx, req.DealID, req.InvitationID)
	if err != nil {
		return nil, err
	}
	if profileID > 0 && row.MemberID != profileID {
		return nil, _errors.ForbiddenException("you can only reject your own invitation")
	}
	if row.Status != uint32(enums.DealMemberStatusInvited) {
		return nil, _errors.BadRequestException("invitation is not in pending status")
	}
	now := time.Now()
	row.Status = uint32(enums.DealMemberStatusRejected)
	row.RespondedAt = &now
	row.Message = req.Reason
	if err := s.members.Update(ctx, row); err != nil {
		return nil, err
	}
	s.bumpInvitationSync(ctx, dealID)
	return &bdspropb.RejectDealInvitationResponse{
		InvitationId: req.InvitationID,
		Success:      true,
	}, nil
}

func (s *LocalDealInvitationService) loadInvitationForDeal(ctx context.Context, dealID, invitationID uint64) (*entity.DealMember, uint64, error) {
	if dealID == 0 || invitationID == 0 {
		return nil, 0, _errors.BadRequestException("deal id and invitation id are required")
	}
	row, err := s.members.GetByID(ctx, invitationID)
	if err != nil {
		return nil, 0, err
	}
	if row == nil {
		return nil, 0, _errors.NotFoundException("invitation not found")
	}
	if row.DealID != dealID {
		return nil, 0, _errors.BadRequestException("invitation does not belong to this deal")
	}
	return row, dealID, nil
}

func (s *LocalDealInvitationService) bumpInvitationSync(ctx context.Context, dealID uint64) {
	if s.sync == nil || dealID == 0 {
		return
	}
	now := time.Now().UnixMilli()
	_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyDealInvitation, dealID), now)
	_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyDealMembers, dealID), now)
	_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyDealPermission, dealID), now)
}

func (s *LocalDealInvitationService) clearDealMeSync(ctx context.Context, profileID uint64) {
	if s.sync == nil || profileID == 0 {
		return
	}
	key := s.sync.GetKeyForProfile(profileID, 0, _utils.SyncKeyDealMe)
	_ = s.sync.Del(ctx, key)
}
