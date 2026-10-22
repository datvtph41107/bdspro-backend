package services

import (
	"context"
	"fmt"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"

	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/enums"
	"mdn4/modules/bdspro/repository"
	bdsutils "mdn4/modules/bdspro/utils"
)

// DealCommissionService owns deal member commission updates (migrated from DealCommissionService).
type DealCommissionService interface {
	UpdateMembersCommission(ctx context.Context, req dto.UpdateDealMembersCommissionRequest) (*bdspropb.UpdateDealMembersCommissionResponse, error)
}

// @bind: mdn4/modules/bdspro/services.DealCommissionService
type LocalDealCommissionService struct {
	sync    *_utils.SyncUtil
	deal    repository.DealRepository
	members repository.DealMemberRepository
	history *bdsutils.DealHistoryLogger
}

func NewLocalDealCommissionService(
	sync *_utils.SyncUtil,
	deal repository.DealRepository,
	members repository.DealMemberRepository,
	history *bdsutils.DealHistoryLogger,
) *LocalDealCommissionService {
	return &LocalDealCommissionService{sync: sync, deal: deal, members: members, history: history}
}

func (s *LocalDealCommissionService) UpdateMembersCommission(
	ctx context.Context,
	req dto.UpdateDealMembersCommissionRequest,
) (*bdspropb.UpdateDealMembersCommissionResponse, error) {
	dealID := req.DealID
	if dealID == 0 {
		return nil, _errors.BadRequestException("deal id is required")
	}
	if len(req.Members) == 0 {
		return nil, _errors.BadRequestException("members is required")
	}
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, _errors.ForbiddenException("profile is required")
	}

	deal, err := s.deal.GetByID(ctx, dealID)
	if err != nil {
		return nil, err
	}
	if deal == nil {
		return nil, _errors.NotFoundException("deal not found")
	}
	ok, err := s.deal.IsAcceptedMember(ctx, dealID, profileID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, _errors.ForbiddenException("you are not allowed to access this deal")
	}

	outMembers := make([]*bdspropb.DealMemberCommissionUpdate, 0, len(req.Members))
	for i, m := range req.Members {
		if m.MemberID == 0 {
			return nil, _errors.BadRequestException(fmt.Sprintf("members[%d].memberId is required", i))
		}
		// if m.CommissionValue <= 0 {
		// 	return nil, _errors.BadRequestException(fmt.Sprintf("members[%d].commissionValue must be greater than 0", i))
		// }
		commissionType := enums.CommissionType(m.CommissionType)
		if commissionType == 0 {
			commissionType = enums.CommissionTypePercent
		}
		if commissionType != enums.CommissionTypePercent && commissionType != enums.CommissionTypeVND {
			return nil, _errors.BadRequestException(fmt.Sprintf("members[%d].commissionType is invalid", i))
		}
		if commissionType == enums.CommissionTypePercent && m.CommissionValue > 100 {
			return nil, _errors.BadRequestException(fmt.Sprintf("members[%d].commissionValue cannot exceed 100%%", i))
		}

		row, err := s.members.GetByDealIDAndMemberID(ctx, dealID, m.MemberID)
		if err != nil {
			return nil, err
		}
		if row == nil {
			return nil, _errors.NotFoundException(fmt.Sprintf("deal member %d not found", m.MemberID))
		}
		if err := s.validateCommission(ctx, dealID, deal.TargetProfit, m.CommissionValue, commissionType); err != nil {
			return nil, err
		}
		if err := s.members.UpdateCommission(ctx, dealID, m.MemberID, m.CommissionValue, commissionType, m.Note); err != nil {
			return nil, err
		}
		outMembers = append(outMembers, &bdspropb.DealMemberCommissionUpdate{
			MemberId:        m.MemberID,
			CommissionValue: m.CommissionValue,
			CommissionType:  uint32(commissionType),
			Note:            m.Note,
		})
	}

	s.bumpCommissionSync(ctx, dealID)
	if s.history != nil {
		s.history.LogDealHistory(ctx, dealID, enums.DealHistoryEventUpdateMembersCommission,
			fmt.Sprintf("Cập nhật hoa hồng cho %d thành viên", len(outMembers)),
			map[string]interface{}{
				"deal_id":      dealID,
				"member_count": len(outMembers),
			})
	}

	return &bdspropb.UpdateDealMembersCommissionResponse{
		DealId:  dealID,
		Members: outMembers,
	}, nil
}

func (s *LocalDealCommissionService) validateCommission(
	ctx context.Context,
	dealID uint64,
	targetProfit float64,
	commissionValue float64,
	commissionType enums.CommissionType,
) error {
	currentTotal, err := s.members.CalculateTotalCommission(ctx, dealID)
	if err != nil {
		return err
	}
	newAmount, err := s.calculateCommissionAmount(targetProfit, commissionValue, commissionType)
	if err != nil {
		return err
	}
	if currentTotal+newAmount > targetProfit {
		return _errors.BadRequestException("total commission cannot exceed target profit")
	}
	return nil
}

func (s *LocalDealCommissionService) calculateCommissionAmount(
	targetProfit, commissionValue float64,
	commissionType enums.CommissionType,
) (float64, error) {
	switch commissionType {
	case enums.CommissionTypePercent:
		return (commissionValue / 100) * targetProfit, nil
	case enums.CommissionTypeVND:
		return commissionValue, nil
	default:
		return 0, _errors.BadRequestException("invalid commission type")
	}
}

func (s *LocalDealCommissionService) bumpCommissionSync(ctx context.Context, dealID uint64) {
	if s.sync == nil || dealID == 0 {
		return
	}
	now := time.Now().UnixMilli()
	_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyDealCommission, dealID), now)
	_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyDealMembers, dealID), now)
	// Commission lives on deal_members rows; drop invitation list cache so clients refetch.
	_ = s.sync.Del(ctx, s.sync.GetKey(ctx, _utils.SyncKeyDealInvitation, dealID))
}
