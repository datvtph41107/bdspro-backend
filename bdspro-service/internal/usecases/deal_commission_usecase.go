package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	"context"
	"errors"
	"fmt"
)

type DealCommissionUsecase interface {
	UpdateDealMemberCommission(ctx context.Context, dealID, memberID uint64, commissionValue float64, commissionType enums.CommissionType, note string) error
	UpdateDealMemberNote(ctx context.Context, dealID, memberID uint64, note string) error
	UpdateDealMemberUnilateralStatus(ctx context.Context, dealID, memberID uint64, isUnilateral bool) error
	UpdateDealMembersCommission(ctx context.Context, dealID uint64, members []*domain.DealMember) error
	GetDealCommissionStats(ctx context.Context, dealID uint64) (*domain.CommissionStats, error)
	ValidateCommission(ctx context.Context, dealID uint64, commissionValue float64, commissionType enums.CommissionType) error
	ValidateMembersCommission(ctx context.Context, dealID uint64, members []*domain.DealMember) error
	CalculateCommissionAmount(ctx context.Context, dealID uint64, commissionValue float64, commissionType enums.CommissionType) (float64, error)
}

type dealCommissionUsecase struct {
	dealMemberRepository repo.DealMemberRepository
	dealRepository       repo.DealRepository
	// groupAuthUsecase     GroupAuthUsecase
	// logWorker            *LogWorker
}

func NewDealCommissionUsecase(
	dealMemberRepository repo.DealMemberRepository,
	dealRepository repo.DealRepository,
	// groupAuthUsecase GroupAuthUsecase,
	// logWorker *LogWorker,
) DealCommissionUsecase {
	return &dealCommissionUsecase{
		dealMemberRepository: dealMemberRepository,
		dealRepository:       dealRepository,
		// groupAuthUsecase:     groupAuthUsecase,
		// logWorker:            logWorker,
	}
}

func (u *dealCommissionUsecase) UpdateDealMemberCommission(ctx context.Context, dealID, memberID uint64, commissionValue float64, commissionType enums.CommissionType, note string) error {
	// currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra deal có tồn tại không
	deal, err := u.dealRepository.GetByID(ctx, dealID)
	if err != nil {
		return err
	}
	if deal == nil {
		return errors.New("deal not found")
	}

	// Kiểm tra quyền - chỉ leader mới được cập nhật hoa hồng
	// err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	// if err != nil {
	// 	return errors.New("you are not allowed to update commission")
	// }

	// Kiểm tra member có tồn tại trong deal không
	_, err = u.dealMemberRepository.GetByDealIDAndMemberID(ctx, dealID, memberID)
	if err != nil {
		return errors.New("deal member not found")
	}

	// Validate hoa hồng
	err = u.ValidateCommission(ctx, dealID, commissionValue, commissionType)
	if err != nil {
		return err
	}

	// Cập nhật hoa hồng
	err = u.dealMemberRepository.UpdateCommission(ctx, dealID, memberID, commissionValue, commissionType, note)
	if err != nil {
		return err
	}

	// Ghi log
	// u.logWorker.Push(LogEvent{
	// 	GroupId: uint32(deal.OwnerId),
	// 	ActorId: uint32(currentUserId),
	// 	LogType: "UPDATE_DEAL_MEMBER_COMMISSION",
	// 	LogData: fmt.Sprintf("Updated commission for member %d in deal %d", memberID, dealID),
	// })

	return nil
}

func (u *dealCommissionUsecase) UpdateDealMemberNote(ctx context.Context, dealID, memberID uint64, note string) error {
	// currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra deal có tồn tại không
	deal, err := u.dealRepository.GetByID(ctx, dealID)
	if err != nil {
		return err
	}
	if deal == nil {
		return errors.New("deal not found")
	}

	// Kiểm tra quyền - chỉ leader mới được cập nhật hoa hồng
	// err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	// if err != nil {
	// 	return errors.New("you are not allowed to update commission")
	// }

	// Kiểm tra member có tồn tại trong deal không
	_, err = u.dealMemberRepository.GetByDealIDAndMemberID(ctx, dealID, memberID)
	if err != nil {
		return errors.New("deal member not found")
	}

	err = u.dealMemberRepository.UpdateNote(ctx, dealID, memberID, note)
	if err != nil {
		return err
	}

	// Ghi log
	// u.logWorker.Push(LogEvent{
	// 	GroupId: uint32(deal.OwnerId),
	// 	ActorId: uint32(currentUserId),
	// 	LogType: "UPDATE_DEAL_MEMBER_NOTE",
	// 	LogData: fmt.Sprintf("Updated note for member %d in deal %d", memberID, dealID),
	// })

	return nil
}

func (u *dealCommissionUsecase) UpdateDealMemberUnilateralStatus(ctx context.Context, dealID, memberID uint64, isUnilateral bool) error {
	// currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra deal có tồn tại không
	deal, err := u.dealRepository.GetByID(ctx, dealID)
	if err != nil {
		return err
	}
	if deal == nil {
		return errors.New("deal not found")
	}

	// Kiểm tra quyền - chỉ leader mới được cập nhật hoa hồng
	// err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	// if err != nil {
	// 	return errors.New("you are not allowed to update commission")
	// }

	// Kiểm tra member có tồn tại trong deal không
	_, err = u.dealMemberRepository.GetByDealIDAndMemberID(ctx, dealID, memberID)
	if err != nil {
		return errors.New("deal member not found")
	}

	err = u.dealMemberRepository.UpdateUnilateralStatus(ctx, dealID, memberID, isUnilateral)
	if err != nil {
		return err
	}

	// Ghi log
	// u.logWorker.Push(LogEvent{
	// 	GroupId: uint32(deal.OwnerId),
	// 	ActorId: uint32(currentUserId),
	// 	LogType: "UPDATE_DEAL_MEMBER_UNILATERAL_STATUS",
	// 	LogData: fmt.Sprintf("Updated unilateral status for member %d in deal %d to %t", memberID, dealID, isUnilateral),
	// })

	return nil
}

func (u *dealCommissionUsecase) UpdateDealMembersCommission(ctx context.Context, dealID uint64, members []*domain.DealMember) error {
	// currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra deal có tồn tại không
	deal, err := u.dealRepository.GetByID(ctx, dealID)
	if err != nil {
		return err
	}
	if deal == nil {
		return errors.New("deal not found")
	}

	// Kiểm tra quyền - chỉ leader mới được cập nhật hoa hồng
	// err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	// if err != nil {
	// 	return errors.New("you are not allowed to update commission")
	// }

	// Kiểm tra tất cả members có tồn tại trong deal không
	for _, member := range members {
		_, err = u.dealMemberRepository.GetByDealIDAndMemberID(ctx, dealID, member.MemberID)
		if err != nil {
			return errors.New("deal member not found")
		}
	}

	// Validate hoa hồng cho tất cả members
	for _, member := range members {
		err = u.ValidateCommission(ctx, dealID, member.CommissionValue, member.CommissionType)
		if err != nil {
			return err
		}
	}

	// Cập nhật hoa hồng cho tất cả members
	for _, member := range members {
		err = u.dealMemberRepository.UpdateCommission(ctx, dealID, member.MemberID, member.CommissionValue, member.CommissionType, member.Note)
		if err != nil {
			return err
		}
	}

	// Ghi log
	// u.logWorker.Push(LogEvent{
	// 	GroupId: uint32(deal.OwnerId),
	// 	ActorId: uint32(currentUserId),
	// 	LogType: "UPDATE_DEAL_MEMBERS_COMMISSION",
	// 	LogData: fmt.Sprintf("Updated commission for %d members in deal %d", len(members), dealID),
	// })

	return nil
}

func (u *dealCommissionUsecase) GetDealCommissionStats(ctx context.Context, dealID uint64) (*domain.CommissionStats, error) {
	deal, err := u.dealRepository.GetByID(ctx, dealID)
	if err != nil {
		return nil, err
	}
	if deal == nil {
		return nil, errors.New("deal not found")
	}

	stats, err := u.dealMemberRepository.GetCommissionStats(ctx, dealID)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (u *dealCommissionUsecase) ValidateCommission(ctx context.Context, dealID uint64, commissionValue float64, commissionType enums.CommissionType) error {
	// Kiểm tra giá trị hoa hồng phải > 0
	if commissionValue <= 0 {
		return errors.New("commission value must be greater than 0")
	}

	// Kiểm tra loại hoa hồng hợp lệ
	if commissionType != enums.CommissionTypePercent && commissionType != enums.CommissionTypeVND {
		return errors.New("invalid commission type")
	}

	// Nếu là hoa hồng theo %, kiểm tra không vượt quá 100%
	if commissionType == enums.CommissionTypePercent && commissionValue > 100 {
		return errors.New("commission percentage cannot exceed 100%")
	}

	// Tính tổng hoa hồng hiện tại
	currentTotalCommission, err := u.dealMemberRepository.CalculateTotalCommission(ctx, dealID)
	if err != nil {
		return err
	}

	// Tính hoa hồng mới sẽ thêm vào
	newCommissionAmount, err := u.CalculateCommissionAmount(ctx, dealID, commissionValue, commissionType)
	if err != nil {
		return err
	}

	// Lấy thông tin deal để kiểm tra lợi nhuận
	deal, err := u.dealRepository.GetByID(ctx, dealID)
	if err != nil {
		return err
	}

	// Tính tổng hoa hồng sau khi cập nhật
	totalCommissionAfterUpdate := currentTotalCommission + newCommissionAmount

	// Kiểm tra tổng hoa hồng không được vượt quá lợi nhuận mục tiêu
	if totalCommissionAfterUpdate > deal.TargetProfit {
		return errors.New("total commission cannot exceed target profit")
	}

	return nil
}

func (u *dealCommissionUsecase) ValidateMembersCommission(ctx context.Context, dealID uint64, members []*domain.DealMember) error {
	for _, member := range members {
		err := u.ValidateCommission(ctx, dealID, member.CommissionValue, member.CommissionType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *dealCommissionUsecase) CalculateCommissionAmount(ctx context.Context, dealID uint64, commissionValue float64, commissionType enums.CommissionType) (float64, error) {
	// Lấy thông tin deal
	deal, err := u.dealRepository.GetByID(ctx, dealID)
	if err != nil {
		return 0, err
	}

	// Tính số tiền hoa hồng dựa trên loại
	var commissionAmount float64
	switch commissionType {
	case enums.CommissionTypePercent:
		// Hoa hồng theo % của lợi nhuận mục tiêu
		commissionAmount = (commissionValue / 100) * deal.TargetProfit
	case enums.CommissionTypeVND:
		// Hoa hồng theo tiền VND trực tiếp
		commissionAmount = commissionValue
	default:
		return 0, fmt.Errorf("invalid commission type")
	}

	return commissionAmount, nil
}
