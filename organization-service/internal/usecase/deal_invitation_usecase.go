package usecase

import (
	"context"
	"fmt"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"
	"organization/internal/enums"
	iusecase "organization/internal/interface"
	sharepb "pb/types/shared"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type DealInvitationUsecase interface {
	// 1. Gửi lời mời
	SendInvitation(ctx context.Context, invitation *entity.DealMember) (*entity.DealMember, error)

	// 2. Xác nhận lời mời (dành cho user)
	AcceptInvitation(ctx context.Context, invitationID uint64) error

	// 3. Từ chối lời mời (dành cho user)
	RejectInvitation(ctx context.Context, invitationID uint64, reason string) error

	// 4. Gửi lại lời mời
	ResendInvitation(ctx context.Context, invitationID uint64) (*entity.DealMember, error)

	// 5. Lấy danh sách thành viên đã accept
	GetAcceptedMembers(ctx context.Context, dealID uint64) ([]*dto.DealInvitationWithProfile, error)

	// 6. Rút khỏi thương vụ (dành cho thành viên)
	WithdrawFromDeal(ctx context.Context, invitationID uint64, reason string) error

	// 7. Gỡ khỏi thương vụ (dành cho admin/phụ trách)
	RemoveFromDeal(ctx context.Context, invitationID uint64, reason string) error

	// 8. Xác nhận rút khỏi thương vụ (dành cho admin/phụ trách)
	ConfirmWithdrawal(ctx context.Context, invitationID uint64) error

	// 9. Tìm kiếm thành viên + sort theo thương vụ chung
	SearchMembers(ctx context.Context, dealID uint64, keyword string, page, size int) ([]*dto.DealInvitationWithProfile, uint32, error)

	// Additional methods
	GetInvitationByID(ctx context.Context, invitationID uint64) (*entity.DealMember, error)
	GetPendingInvitations(ctx context.Context, inviteeID uint64, page, size int) ([]*dto.DealInvitationWithProfile, uint32, error)
	GetInvitationsByDealID(ctx context.Context, dealID uint64, page, size int) ([]*dto.DealInvitationWithProfile, uint32, error)
}

type dealInvitationUsecase struct {
	// dealMemberRepository repository.DealInvitationRepository
	dealMemberRepository repository.DealMemberRepository
	dealRepository       repository.DealRepository
	historyUsecase       *EventHistoryUsecase
	groupAuthUsecase     GroupAuthUsecase
	logWorker            *LogWorker
	transaction          iusecase.ITransaction
	notiClient           iusecase.INotiClient
	userClient           iusecase.IUserClient
}

func NewDealInvitationUsecase(
	// dealMemberRepository repository.DealInvitationRepository,
	dealMemberRepository repository.DealMemberRepository,
	dealRepository repository.DealRepository,
	historyUsecase *EventHistoryUsecase,
	groupAuthUsecase GroupAuthUsecase,
	logWorker *LogWorker,
	transaction iusecase.ITransaction,
	notiClient iusecase.INotiClient,
	userClient iusecase.IUserClient,
) DealInvitationUsecase {
	return &dealInvitationUsecase{
		// dealMemberRepository: dealMemberRepository,
		dealMemberRepository: dealMemberRepository,
		dealRepository:       dealRepository,
		historyUsecase:       historyUsecase,
		groupAuthUsecase:     groupAuthUsecase,
		logWorker:            logWorker,
		transaction:          transaction,
		notiClient:           notiClient,
		userClient:           userClient,
	}
}

func (u *dealInvitationUsecase) SendInvitation(ctx context.Context, invitation *entity.DealMember) (*entity.DealMember, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	// Kiểm tra quyền - chỉ leader mới được gửi lời mời
	deal, err := u.dealRepository.GetByID(ctx, invitation.DealID)
	if err != nil {
		return nil, err
	}
	if deal == nil {
		return nil, custom_error.RecordNotFound("deal not found")
	}

	// Kiểm tra xem đã có lời mời cho user này chưa
	exists, err := u.dealMemberRepository.ExistsByDealIDAndInviteeID(ctx, invitation.DealID, invitation.MemberID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, _errors.ReturnError(400, "invitation already exists for this user")
		// return nil, custom_error.InvalidRequest(&errdetails.BadRequest{
		// 	FieldViolations: []*errdetails.BadRequest_FieldViolation{
		// 		{
		// 			Field:       "member_id",
		// 			Description: "invitation already exists for this user",
		// 		},
		// 	},
		// })
	}

	// Kiểm tra user không tự mời chính mình
	if invitation.MemberID == currentUserId {
		return nil, custom_error.InvalidRequest(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "member_id",
					Description: "cannot invite yourself",
				},
			},
		})
	}

	// Set thông tin cho invitation
	invitation.InviterID = currentUserId
	invitation.Status = entity.DealMemberStatusInvited
	invitation.InvitedAt = time.Now()

	// Tạo invitation
	createdInvitation, err := u.dealMemberRepository.Create(ctx, invitation)
	if err != nil {
		return nil, err
	}

	// Log activity
	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		ActorId: uint32(currentUserId),
		LogType: "SEND_INVITATION",
		LogData: fmt.Sprintf("Sent invitation to user %d for deal %d", invitation.MemberID, invitation.DealID),
	})

	// Gửi thông báo
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		u.sendInvitationNotification(cloneCtx, createdInvitation, deal)
	}()

	// Ghi lại lịch sử gửi lời mời
	u.historyUsecase.LogDealHistory(ctx, invitation.DealID, enums.DealHistoryEventSendInvitation,
		fmt.Sprintf("Gửi lời mời tham gia thương vụ '%s' cho %d", deal.Name, invitation.MemberID),
		map[string]interface{}{
			"deal_id":    invitation.DealID,
			"deal_name":  deal.Name,
			"member_id":  invitation.MemberID,
			"inviter_id": currentUserId,
		})

	return createdInvitation, nil
}

func (u *dealInvitationUsecase) AcceptInvitation(ctx context.Context, invitationID uint64) error {
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Lấy thông tin invitation
	invitation, err := u.dealMemberRepository.GetByID(ctx, invitationID)
	if err != nil {
		return err
	}
	if invitation == nil {
		return custom_error.RecordNotFound("invitation not found")
	}

	// Kiểm tra quyền - chỉ người được mời mới được accept
	// if invitation.MemberID != currentUserId {
	// 	return custom_error.Forbidden("you can only accept your own invitation")
	// }

	// Kiểm tra trạng thái
	if !invitation.IsInvited() {
		return _errors.ReturnError(400, "invitation is not in pending status")
	}

	// Lấy thông tin deal
	deal, err := u.dealRepository.GetByID(ctx, invitation.DealID)
	if err != nil {
		return err
	}

	// Thực hiện transaction
	err = u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		// Cập nhật trạng thái invitation
		now := time.Now()
		invitation.Status = entity.DealMemberStatusAccepted
		invitation.RespondedAt = &now

		_, err = u.dealMemberRepository.Update(ctx, invitation)
		if err != nil {
			return err
		}

		// Tạo deal member khi accept invitation
		// amountCommit := 0.0
		// if invitation.AmountCommit != 0 {
		// 	amountCommit = invitation.AmountCommit
		// }
		// dealMember := &entity.DealMember{
		// 	DealID:   invitation.DealID,
		// 	MemberID: invitation.MemberID,
		// 	// MemberType:      invitation.MemberType,
		// 	AmountCommit:    invitation.AmountCommit,
		// 	CommissionValue: 0, // Mặc định 0
		// 	CommissionType:  enums.CommissionTypePercent,
		// 	Note:            invitation.Message,
		// 	IsUnilateral:    false,
		// }
		// _, err = u.dealMemberRepository.CreateMember(ctx, dealMember)
		// if err != nil {
		// 	return err
		// }

		return nil
	})

	if err != nil {
		return err
	}

	// Log activity
	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		ActorId: uint32(currentUserId),
		LogType: "ACCEPT_INVITATION",
		LogData: fmt.Sprintf("User %d accepted invitation for deal %d", currentUserId, invitation.DealID),
	})

	// Gửi thông báo
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		u.sendAcceptanceNotification(cloneCtx, invitation, deal)
	}()

	// Ghi lại lịch sử chấp nhận lời mời
	u.historyUsecase.LogDealHistory(ctx, invitation.DealID, enums.DealHistoryEventAcceptInvitation,
		fmt.Sprintf("Chấp nhận lời mời tham gia thương vụ '%s'", deal.Name),
		map[string]interface{}{
			"deal_id":       invitation.DealID,
			"deal_name":     deal.Name,
			"member_id":     currentUserId,
			"invitation_id": invitationID,
		})

	return nil
}

func (u *dealInvitationUsecase) RejectInvitation(ctx context.Context, invitationID uint64, reason string) error {
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Lấy thông tin invitation
	invitation, err := u.dealMemberRepository.GetByID(ctx, invitationID)
	if err != nil {
		return err
	}
	if invitation == nil {
		return custom_error.RecordNotFound("invitation not found")
	}

	// Kiểm tra quyền
	if invitation.MemberID != currentUserId {
		return custom_error.Forbidden("you can only reject your own invitation")
	}

	// Kiểm tra trạng thái
	if !invitation.IsInvited() {
		return custom_error.InvalidRequest(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "status",
					Description: "invitation is not in pending status",
				},
			},
		})
	}

	// Lấy thông tin deal
	deal, err := u.dealRepository.GetByID(ctx, invitation.DealID)
	if err != nil {
		return err
	}

	// Cập nhật trạng thái
	now := time.Now()
	invitation.Status = entity.DealMemberStatusRejected
	invitation.RespondedAt = &now

	invitation.Message = reason // Lưu lý do từ chối

	_, err = u.dealMemberRepository.Update(ctx, invitation)
	if err != nil {
		return err
	}

	// Log activity
	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		ActorId: uint32(currentUserId),
		LogType: "REJECT_INVITATION",
		LogData: fmt.Sprintf("User %d rejected invitation for deal %d", currentUserId, invitation.DealID),
	})

	// Gửi thông báo
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		u.sendRejectionNotification(cloneCtx, invitation, deal)
	}()

	// Ghi lại lịch sử từ chối lời mời
	u.historyUsecase.LogDealHistory(ctx, invitation.DealID, enums.DealHistoryEventRejectInvitation,
		fmt.Sprintf("Từ chối lời mời tham gia thương vụ '%s' với lý do: %s", deal.Name, reason),
		map[string]interface{}{
			"deal_id":       invitation.DealID,
			"deal_name":     deal.Name,
			"member_id":     currentUserId,
			"invitation_id": invitationID,
			"reason":        reason,
		})

	return nil
}

func (u *dealInvitationUsecase) ResendInvitation(ctx context.Context, invitationID uint64) (*entity.DealMember, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Lấy thông tin invitation
	invitation, err := u.dealMemberRepository.GetByID(ctx, invitationID)
	if err != nil {
		return nil, err
	}
	if invitation == nil {
		return nil, custom_error.RecordNotFound("invitation not found")
	}

	// Kiểm tra quyền - chỉ leader mới được gửi lại
	deal, err := u.dealRepository.GetByID(ctx, invitation.DealID)
	if err != nil {
		return nil, err
	}

	err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to resend invitation")
	}

	// Kiểm tra có thể gửi lại không
	if !invitation.CanResend() {
		return nil, custom_error.InvalidRequest(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "status",
					Description: "invitation cannot be resent",
				},
			},
		})
	}

	// Cập nhật trạng thái
	invitation.Status = entity.DealMemberStatusInvited
	invitation.InvitedAt = time.Now()

	updatedInvitation, err := u.dealMemberRepository.Update(ctx, invitation)
	if err != nil {
		return nil, err
	}

	// Log activity
	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		ActorId: uint32(currentUserId),
		LogType: "RESEND_INVITATION",
		LogData: fmt.Sprintf("Resent invitation to user %d for deal %d", invitation.MemberID, invitation.DealID),
	})

	// Gửi thông báo
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		u.sendInvitationNotification(cloneCtx, updatedInvitation, deal)
	}()

	return updatedInvitation, nil
}

func (u *dealInvitationUsecase) GetAcceptedMembers(ctx context.Context, dealID uint64) ([]*dto.DealInvitationWithProfile, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra quyền - member của group mới được xem
	deal, err := u.dealRepository.GetByID(ctx, dealID)
	if err != nil {
		return nil, err
	}

	err = u.groupAuthUsecase.IsMemberGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to view members")
	}

	// Lấy danh sách invitation đã accept
	invitations, err := u.dealMemberRepository.GetAcceptedByDealID(ctx, dealID)
	if err != nil {
		return nil, err
	}

	// Populate user profiles
	return u.populateUserProfiles(ctx, invitations)
}

func (u *dealInvitationUsecase) WithdrawFromDeal(ctx context.Context, invitationID uint64, reason string) error {
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Lấy thông tin invitation
	invitation, err := u.dealMemberRepository.GetByID(ctx, invitationID)
	if err != nil {
		return err
	}
	if invitation == nil {
		return custom_error.RecordNotFound("invitation not found")
	}

	// Kiểm tra quyền - chỉ thành viên đã accept mới được rút
	if invitation.MemberID != currentUserId {
		return custom_error.Forbidden("you can only withdraw your own participation")
	}

	if !invitation.IsAccepted() {
		return custom_error.InvalidRequest(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "status",
					Description: "you can only withdraw from accepted invitation",
				},
			},
		})
	}

	// Lấy thông tin deal
	deal, err := u.dealRepository.GetByID(ctx, invitation.DealID)
	if err != nil {
		return err
	}

	// Thực hiện transaction đảm bảo ACID
	err = u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		// Cập nhật trạng thái invitation
		now := time.Now()
		invitation.Status = entity.DealMemberStatusWithdrawn
		invitation.WithdrawnAt = &now

		invitation.Message = reason

		_, err = u.dealMemberRepository.Update(ctx, invitation)
		if err != nil {
			return err
		}

		// Xóa deal member nếu đã tồn tại
		existingMember, err := u.dealMemberRepository.GetByDealIDAndMemberID(ctx, invitation.DealID, invitation.MemberID)
		if err != nil {
			return err
		}
		if existingMember != nil {
			// Xóa hẳn deal member
			err = u.dealMemberRepository.DeleteMember(ctx, invitation.DealID, invitation.MemberID)
			if err != nil {
				return err
			}
		}

		return nil
	})

	// Log activity
	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		ActorId: uint32(currentUserId),
		LogType: "WITHDRAW_FROM_DEAL",
		LogData: fmt.Sprintf("User %d withdrew from deal %d", currentUserId, invitation.DealID),
	})

	// Gửi thông báo
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		u.sendWithdrawalNotification(cloneCtx, invitation, deal)
	}()

	// Ghi lại lịch sử rút khỏi thương vụ
	u.historyUsecase.LogDealHistory(ctx, invitation.DealID, enums.DealHistoryEventWithdrawFromDeal,
		fmt.Sprintf("Rút khỏi thương vụ '%s' với lý do: %s", deal.Name, reason),
		map[string]interface{}{
			"deal_id":       invitation.DealID,
			"deal_name":     deal.Name,
			"member_id":     currentUserId,
			"invitation_id": invitationID,
			"reason":        reason,
		})

	return nil
}

func (u *dealInvitationUsecase) RemoveFromDeal(ctx context.Context, invitationID uint64, reason string) error {
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Lấy thông tin invitation
	invitation, err := u.dealMemberRepository.GetByID(ctx, invitationID)
	if err != nil {
		return err
	}
	if invitation == nil {
		return custom_error.RecordNotFound("invitation not found")
	}

	// Kiểm tra quyền - chỉ leader mới được gỡ thành viên
	deal, err := u.dealRepository.GetByID(ctx, invitation.DealID)
	if err != nil {
		return err
	}

	err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	if err != nil {
		return custom_error.Forbidden("you are not allowed to remove members")
	}

	// TODO: Kiểm tra xem có góp vốn nào chưa được phê duyệt không
	// if hasUnapprovedInvestment {
	//     return custom_error.BadRequest("cannot remove member with unapproved investments")
	// }

	// Thực hiện transaction đảm bảo ACID
	err = u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		// Cập nhật trạng thái invitation
		now := time.Now()
		invitation.Status = entity.DealMemberStatusWithdrawn
		invitation.WithdrawnAt = &now

		invitation.Message = reason

		_, err = u.dealMemberRepository.Update(ctx, invitation)
		if err != nil {
			return err
		}

		// Xóa deal member nếu đã tồn tại
		existingMember, err := u.dealMemberRepository.GetByDealIDAndMemberID(ctx, invitation.DealID, invitation.MemberID)
		if err != nil {
			return err
		}
		if existingMember != nil {
			// Xóa hẳn deal member
			err = u.dealMemberRepository.DeleteMember(ctx, invitation.DealID, invitation.MemberID)
			if err != nil {
				return err
			}
		}

		return nil
	})

	// Log activity
	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		ActorId: uint32(currentUserId),
		LogType: "REMOVE_FROM_DEAL",
		LogData: fmt.Sprintf("Removed user %d from deal %d", invitation.MemberID, invitation.DealID),
	})

	// Gửi thông báo
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		u.sendRemovalNotification(cloneCtx, invitation, deal)
	}()

	// Ghi lại lịch sử gỡ khỏi thương vụ
	u.historyUsecase.LogDealHistory(ctx, invitation.DealID, enums.DealHistoryEventRemoveFromDeal,
		fmt.Sprintf("Gỡ user %d khỏi thương vụ '%s' với lý do: %s", invitation.MemberID, deal.Name, reason),
		map[string]interface{}{
			"deal_id":       invitation.DealID,
			"deal_name":     deal.Name,
			"member_id":     invitation.MemberID,
			"removed_by":    currentUserId,
			"invitation_id": invitationID,
			"reason":        reason,
		})

	return nil
}

func (u *dealInvitationUsecase) ConfirmWithdrawal(ctx context.Context, invitationID uint64) error {
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Lấy thông tin invitation
	invitation, err := u.dealMemberRepository.GetByID(ctx, invitationID)
	if err != nil {
		return err
	}
	if invitation == nil {
		return custom_error.RecordNotFound("invitation not found")
	}

	// Kiểm tra quyền - chỉ leader mới được xác nhận
	deal, err := u.dealRepository.GetByID(ctx, invitation.DealID)
	if err != nil {
		return err
	}

	err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	if err != nil {
		return custom_error.Forbidden("you are not allowed to confirm withdrawal")
	}

	// Kiểm tra trạng thái
	if invitation.Status != entity.DealMemberStatusWithdrawn {
		return custom_error.InvalidRequest(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "status",
					Description: "invitation is not in withdrawn status",
				},
			},
		})
	}

	// Xóa invitation
	err = u.dealMemberRepository.Delete(ctx, invitationID)
	if err != nil {
		return err
	}

	// Log activity
	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		ActorId: uint32(currentUserId),
		LogType: "CONFIRM_WITHDRAWAL",
		LogData: fmt.Sprintf("Confirmed withdrawal of user %d from deal %d", invitation.MemberID, invitation.DealID),
	})

	// Ghi lại lịch sử xác nhận rút khỏi thương vụ
	u.historyUsecase.LogDealHistory(ctx, invitation.DealID, enums.DealHistoryEventConfirmWithdrawal,
		fmt.Sprintf("Xác nhận rút khỏi thương vụ '%s'", deal.Name),
		map[string]interface{}{
			"deal_id":       invitation.DealID,
			"deal_name":     deal.Name,
			"member_id":     invitation.MemberID,
			"confirm_by":    currentUserId,
			"invitation_id": invitationID,
		})

	return nil
}

func (u *dealInvitationUsecase) SearchMembers(ctx context.Context, dealID uint64, keyword string, page, size int) ([]*dto.DealInvitationWithProfile, uint32, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra quyền
	deal, err := u.dealRepository.GetByID(ctx, dealID)
	if err != nil {
		return nil, 0, err
	}

	err = u.groupAuthUsecase.IsMemberGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	if err != nil {
		return nil, 0, custom_error.Forbidden("you are not allowed to search members")
	}

	// TODO: Implement search logic with keyword
	// For now, just get all invitations
	invitations, total, err := u.dealMemberRepository.GetByDealID(ctx, dealID, page, size)
	if err != nil {
		return nil, 0, err
	}

	// Populate user profiles
	invitationsWithProfiles, err := u.populateUserProfiles(ctx, invitations)
	if err != nil {
		return nil, 0, err
	}

	return invitationsWithProfiles, total, nil
}

// Helper methods
func (u *dealInvitationUsecase) GetInvitationByID(ctx context.Context, invitationID uint64) (*entity.DealMember, error) {
	return u.dealMemberRepository.GetByID(ctx, invitationID)
}

func (u *dealInvitationUsecase) GetPendingInvitations(ctx context.Context, inviteeID uint64, page, size int) ([]*dto.DealInvitationWithProfile, uint32, error) {
	invitations, total, err := u.dealMemberRepository.GetPendingByInviteeID(ctx, inviteeID, page, size)
	if err != nil {
		return nil, 0, err
	}

	invitationsWithProfiles, err := u.populateUserProfiles(ctx, invitations)
	if err != nil {
		return nil, 0, err
	}

	return invitationsWithProfiles, total, nil
}

func (u *dealInvitationUsecase) GetInvitationsByDealID(ctx context.Context, dealID uint64, page, size int) ([]*dto.DealInvitationWithProfile, uint32, error) {
	invitations, total, err := u.dealMemberRepository.GetByDealID(ctx, dealID, page, size)
	if err != nil {
		return nil, 0, err
	}

	invitationsWithProfiles, err := u.populateUserProfiles(ctx, invitations)
	if err != nil {
		return nil, 0, err
	}

	return invitationsWithProfiles, total, nil
}

func (u *dealInvitationUsecase) populateUserProfiles(ctx context.Context, invitations []*entity.DealMember) ([]*dto.DealInvitationWithProfile, error) {
	if len(invitations) == 0 {
		return []*dto.DealInvitationWithProfile{}, nil
	}

	// Collect user IDs
	userIDs := make([]uint64, 0, len(invitations)*2)
	for _, invitation := range invitations {
		userIDs = append(userIDs, invitation.InviterID, invitation.MemberID)
	}

	// Get user profiles
	profiles, err := u.userClient.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: userIDs})
	if err != nil {
		return nil, err
	}

	// Create profile map
	profileMap := make(map[uint64]*sharepb.ProfileItem)
	for _, profile := range profiles.Profiles {
		profileMap[profile.Id] = profile
	}

	// Build result
	result := make([]*dto.DealInvitationWithProfile, len(invitations))
	for i, invitation := range invitations {
		result[i] = &dto.DealInvitationWithProfile{
			DealMember: invitation,
			Inviter:    profileMap[invitation.InviterID],
			Invitee:    profileMap[invitation.MemberID],
		}
	}

	return result, nil
}

// Notification methods
func (u *dealInvitationUsecase) sendInvitationNotification(ctx context.Context, invitation *entity.DealMember, deal *entity.Deal) error {
	notification := &dto.NotificationDTO{
		Type:       enums.NotiDealInvitation,
		UserId:     &invitation.MemberID,
		AttachData: []string{deal.Name, invitation.Message},
		TargetID:   deal.ID,
		IsMerge:    false,
	}

	return u.notiClient.SendNotification(ctx, notification)
}

func (u *dealInvitationUsecase) sendAcceptanceNotification(ctx context.Context, invitation *entity.DealMember, deal *entity.Deal) error {
	notification := &dto.NotificationDTO{
		Type:       enums.NotiDealInvitationAccepted,
		UserId:     &invitation.InviterID,
		AttachData: []string{deal.Name},
		TargetID:   deal.ID,
		IsMerge:    false,
	}

	return u.notiClient.SendNotification(ctx, notification)
}

func (u *dealInvitationUsecase) sendRejectionNotification(ctx context.Context, invitation *entity.DealMember, deal *entity.Deal) error {
	notification := &dto.NotificationDTO{
		Type:       enums.NotiDealInvitationRejected,
		UserId:     &invitation.InviterID,
		AttachData: []string{deal.Name, invitation.Message},
		TargetID:   deal.ID,
		IsMerge:    false,
	}

	return u.notiClient.SendNotification(ctx, notification)
}

func (u *dealInvitationUsecase) sendWithdrawalNotification(ctx context.Context, invitation *entity.DealMember, deal *entity.Deal) error {
	notification := &dto.NotificationDTO{
		Type:       enums.NotiDealMemberWithdrawn,
		UserId:     &invitation.InviterID,
		AttachData: []string{deal.Name, invitation.Message},
		TargetID:   deal.ID,
		IsMerge:    false,
	}

	return u.notiClient.SendNotification(ctx, notification)
}

func (u *dealInvitationUsecase) sendRemovalNotification(ctx context.Context, invitation *entity.DealMember, deal *entity.Deal) error {
	notification := &dto.NotificationDTO{
		Type:       enums.NotiDealMemberRemoved,
		UserId:     &invitation.MemberID,
		AttachData: []string{deal.Name, invitation.Message},
		TargetID:   deal.ID,
		IsMerge:    false,
	}

	return u.notiClient.SendNotification(ctx, notification)
}
