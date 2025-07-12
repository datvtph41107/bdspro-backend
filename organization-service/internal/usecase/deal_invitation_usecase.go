package usecase

import (
	"context"
	"fmt"
	"time"

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
	SendInvitation(ctx context.Context, invitation *entity.DealInvitation) (*entity.DealInvitation, error)
	
	// 2. Xác nhận lời mời (dành cho user)
	AcceptInvitation(ctx context.Context, invitationID uint64) error
	
	// 3. Từ chối lời mời (dành cho user)
	RejectInvitation(ctx context.Context, invitationID uint64, reason string) error
	
	// 4. Gửi lại lời mời
	ResendInvitation(ctx context.Context, invitationID uint64) (*entity.DealInvitation, error)
	
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
	GetInvitationByID(ctx context.Context, invitationID uint64) (*entity.DealInvitation, error)
	GetPendingInvitations(ctx context.Context, inviteeID uint64, page, size int) ([]*dto.DealInvitationWithProfile, uint32, error)
	GetInvitationsByDealID(ctx context.Context, dealID uint64, page, size int) ([]*dto.DealInvitationWithProfile, uint32, error)
}

type dealInvitationUsecase struct {
	dealInvitationRepository repository.DealInvitationRepository
	dealRepository          repository.GroupDealRepository
	groupAuthUsecase        GroupAuthUsecase
	logWorker               *LogWorker
	transaction             iusecase.ITransaction
	notiClient              iusecase.INotiClient
	userClient              iusecase.IUserClient
}

func NewDealInvitationUsecase(
	dealInvitationRepository repository.DealInvitationRepository,
	dealRepository repository.GroupDealRepository,
	groupAuthUsecase GroupAuthUsecase,
	logWorker *LogWorker,
	transaction iusecase.ITransaction,
	notiClient iusecase.INotiClient,
	userClient iusecase.IUserClient,
) DealInvitationUsecase {
	return &dealInvitationUsecase{
		dealInvitationRepository: dealInvitationRepository,
		dealRepository:          dealRepository,
		groupAuthUsecase:        groupAuthUsecase,
		logWorker:               logWorker,
		transaction:             transaction,
		notiClient:              notiClient,
		userClient:              userClient,
	}
}

func (u *dealInvitationUsecase) SendInvitation(ctx context.Context, invitation *entity.DealInvitation) (*entity.DealInvitation, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	
	// Kiểm tra quyền - chỉ leader mới được gửi lời mời
	deal, err := u.dealRepository.GetByID(ctx, invitation.DealID)
	if err != nil {
		return nil, err
	}
	if deal == nil {
		return nil, custom_error.RecordNotFound("deal not found")
	}
	
	err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserId))
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to send invitation")
	}
	
	// Kiểm tra xem đã có lời mời cho user này chưa
	exists, err := u.dealInvitationRepository.ExistsByDealIDAndInviteeID(ctx, invitation.DealID, invitation.InviteeID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, custom_error.InvalidRequest(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "invitee_id",
					Description: "invitation already exists for this user",
				},
			},
		})
	}
	
	// Kiểm tra user không tự mời chính mình
	if invitation.InviteeID == currentUserId {
		return nil, custom_error.InvalidRequest(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "invitee_id",
					Description: "cannot invite yourself",
				},
			},
		})
	}
	
	// Set thông tin cho invitation
	invitation.InviterID = currentUserId
	invitation.Status = entity.DealInvitationStatusInvited
	invitation.InvitedAt = time.Now()
	invitation.CreatedBy = currentUserId
	invitation.UpdatedBy = currentUserId
	
	// Tạo invitation
	createdInvitation, err := u.dealInvitationRepository.Create(ctx, invitation)
	if err != nil {
		return nil, err
	}
	
	// Log activity
	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		ActorId: uint32(currentUserId),
		LogType: "SEND_INVITATION",
		LogData: fmt.Sprintf("Sent invitation to user %d for deal %d", invitation.InviteeID, invitation.DealID),
	})
	
	// Gửi thông báo
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		u.sendInvitationNotification(cloneCtx, createdInvitation, deal)
	}()
	
	return createdInvitation, nil
}

func (u *dealInvitationUsecase) AcceptInvitation(ctx context.Context, invitationID uint64) error {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	
	// Lấy thông tin invitation
	invitation, err := u.dealInvitationRepository.GetByID(ctx, invitationID)
	if err != nil {
		return err
	}
	if invitation == nil {
		return custom_error.RecordNotFound("invitation not found")
	}
	
	// Kiểm tra quyền - chỉ người được mời mới được accept
	if invitation.InviteeID != currentUserId {
		return custom_error.Forbidden("you can only accept your own invitation")
	}
	
	// Kiểm tra trạng thái
	if !invitation.IsPending() {
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
	
	// Thực hiện transaction
	err = u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		// Cập nhật trạng thái invitation
		now := time.Now()
		invitation.Status = entity.DealInvitationStatusAccepted
		invitation.RespondedAt = &now
		invitation.UpdatedBy = currentUserId
		
		_, err = u.dealInvitationRepository.Update(ctx, invitation)
		if err != nil {
			return err
		}
		
		// TODO: Thêm vào deal_members khi có repository
		// dealMember := &entity.DealMember{
		//     DealID:       invitation.DealID,
		//     MemberID:     invitation.InviteeID,
		//     MemberType:   invitation.MemberType,
		//     AmountCommit: 0,
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
	
	return nil
}

func (u *dealInvitationUsecase) RejectInvitation(ctx context.Context, invitationID uint64, reason string) error {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	
	// Lấy thông tin invitation
	invitation, err := u.dealInvitationRepository.GetByID(ctx, invitationID)
	if err != nil {
		return err
	}
	if invitation == nil {
		return custom_error.RecordNotFound("invitation not found")
	}
	
	// Kiểm tra quyền
	if invitation.InviteeID != currentUserId {
		return custom_error.Forbidden("you can only reject your own invitation")
	}
	
	// Kiểm tra trạng thái
	if !invitation.IsPending() {
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
	invitation.Status = entity.DealInvitationStatusRejected
	invitation.RespondedAt = &now
	invitation.UpdatedBy = currentUserId
	invitation.Message = reason // Lưu lý do từ chối
	
	_, err = u.dealInvitationRepository.Update(ctx, invitation)
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
	
	return nil
}

func (u *dealInvitationUsecase) ResendInvitation(ctx context.Context, invitationID uint64) (*entity.DealInvitation, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	
	// Lấy thông tin invitation
	invitation, err := u.dealInvitationRepository.GetByID(ctx, invitationID)
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
	invitation.Status = entity.DealInvitationStatusInvited
	invitation.InvitedAt = time.Now()
	invitation.UpdatedBy = currentUserId
	
	updatedInvitation, err := u.dealInvitationRepository.Update(ctx, invitation)
	if err != nil {
		return nil, err
	}
	
	// Log activity
	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		ActorId: uint32(currentUserId),
		LogType: "RESEND_INVITATION",
		LogData: fmt.Sprintf("Resent invitation to user %d for deal %d", invitation.InviteeID, invitation.DealID),
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
	invitations, err := u.dealInvitationRepository.GetAcceptedByDealID(ctx, dealID)
	if err != nil {
		return nil, err
	}
	
	// Populate user profiles
	return u.populateUserProfiles(ctx, invitations)
}

func (u *dealInvitationUsecase) WithdrawFromDeal(ctx context.Context, invitationID uint64, reason string) error {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	
	// Lấy thông tin invitation
	invitation, err := u.dealInvitationRepository.GetByID(ctx, invitationID)
	if err != nil {
		return err
	}
	if invitation == nil {
		return custom_error.RecordNotFound("invitation not found")
	}
	
	// Kiểm tra quyền - chỉ thành viên đã accept mới được rút
	if invitation.InviteeID != currentUserId {
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
	
	// Cập nhật trạng thái
	now := time.Now()
	invitation.Status = entity.DealInvitationStatusWithdrawn
	invitation.WithdrawnAt = &now
	invitation.UpdatedBy = currentUserId
	invitation.Message = reason
	
	_, err = u.dealInvitationRepository.Update(ctx, invitation)
	if err != nil {
		return err
	}
	
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
	
	return nil
}

func (u *dealInvitationUsecase) RemoveFromDeal(ctx context.Context, invitationID uint64, reason string) error {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	
	// Lấy thông tin invitation
	invitation, err := u.dealInvitationRepository.GetByID(ctx, invitationID)
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
	
	// Cập nhật trạng thái
	invitation.Status = entity.DealInvitationStatusWithdrawn
	invitation.UpdatedBy = currentUserId
	invitation.Message = reason
	
	_, err = u.dealInvitationRepository.Update(ctx, invitation)
	if err != nil {
		return err
	}
	
	// Log activity
	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		ActorId: uint32(currentUserId),
		LogType: "REMOVE_FROM_DEAL",
		LogData: fmt.Sprintf("Removed user %d from deal %d", invitation.InviteeID, invitation.DealID),
	})
	
	// Gửi thông báo
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		u.sendRemovalNotification(cloneCtx, invitation, deal)
	}()
	
	return nil
}

func (u *dealInvitationUsecase) ConfirmWithdrawal(ctx context.Context, invitationID uint64) error {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	
	// Lấy thông tin invitation
	invitation, err := u.dealInvitationRepository.GetByID(ctx, invitationID)
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
	if invitation.Status != entity.DealInvitationStatusWithdrawn {
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
	err = u.dealInvitationRepository.Delete(ctx, invitationID)
	if err != nil {
		return err
	}
	
	// Log activity
	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		ActorId: uint32(currentUserId),
		LogType: "CONFIRM_WITHDRAWAL",
		LogData: fmt.Sprintf("Confirmed withdrawal of user %d from deal %d", invitation.InviteeID, invitation.DealID),
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
	invitations, total, err := u.dealInvitationRepository.GetByDealID(ctx, dealID, page, size)
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
func (u *dealInvitationUsecase) GetInvitationByID(ctx context.Context, invitationID uint64) (*entity.DealInvitation, error) {
	return u.dealInvitationRepository.GetByID(ctx, invitationID)
}

func (u *dealInvitationUsecase) GetPendingInvitations(ctx context.Context, inviteeID uint64, page, size int) ([]*dto.DealInvitationWithProfile, uint32, error) {
	invitations, total, err := u.dealInvitationRepository.GetPendingByInviteeID(ctx, inviteeID, page, size)
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
	invitations, total, err := u.dealInvitationRepository.GetByDealID(ctx, dealID, page, size)
	if err != nil {
		return nil, 0, err
	}
	
	invitationsWithProfiles, err := u.populateUserProfiles(ctx, invitations)
	if err != nil {
		return nil, 0, err
	}
	
	return invitationsWithProfiles, total, nil
}

func (u *dealInvitationUsecase) populateUserProfiles(ctx context.Context, invitations []*entity.DealInvitation) ([]*dto.DealInvitationWithProfile, error) {
	if len(invitations) == 0 {
		return []*dto.DealInvitationWithProfile{}, nil
	}
	
	// Collect user IDs
	userIDs := make([]uint64, 0, len(invitations)*2)
	for _, invitation := range invitations {
		userIDs = append(userIDs, invitation.InviterID, invitation.InviteeID)
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
			DealInvitation: invitation,
			Inviter:        profileMap[invitation.InviterID],
			Invitee:        profileMap[invitation.InviteeID],
		}
	}
	
	return result, nil
}

// Notification methods
func (u *dealInvitationUsecase) sendInvitationNotification(ctx context.Context, invitation *entity.DealInvitation, deal *entity.Deal) error {
	notification := &dto.NotificationDTO{
		Type:       enums.NotiDealInvitation,
		UserId:     &invitation.InviteeID,
		AttachData: []string{deal.Name, invitation.Message},
		TargetID:   deal.ID,
		IsMerge:    false,
	}
	
	return u.notiClient.SendNotification(ctx, notification)
}

func (u *dealInvitationUsecase) sendAcceptanceNotification(ctx context.Context, invitation *entity.DealInvitation, deal *entity.Deal) error {
	notification := &dto.NotificationDTO{
		Type:       enums.NotiDealInvitationAccepted,
		UserId:     &invitation.InviterID,
		AttachData: []string{deal.Name},
		TargetID:   deal.ID,
		IsMerge:    false,
	}
	
	return u.notiClient.SendNotification(ctx, notification)
}

func (u *dealInvitationUsecase) sendRejectionNotification(ctx context.Context, invitation *entity.DealInvitation, deal *entity.Deal) error {
	notification := &dto.NotificationDTO{
		Type:       enums.NotiDealInvitationRejected,
		UserId:     &invitation.InviterID,
		AttachData: []string{deal.Name, invitation.Message},
		TargetID:   deal.ID,
		IsMerge:    false,
	}
	
	return u.notiClient.SendNotification(ctx, notification)
}

func (u *dealInvitationUsecase) sendWithdrawalNotification(ctx context.Context, invitation *entity.DealInvitation, deal *entity.Deal) error {
	notification := &dto.NotificationDTO{
		Type:       enums.NotiDealMemberWithdrawn,
		UserId:     &invitation.InviterID,
		AttachData: []string{deal.Name, invitation.Message},
		TargetID:   deal.ID,
		IsMerge:    false,
	}
	
	return u.notiClient.SendNotification(ctx, notification)
}

func (u *dealInvitationUsecase) sendRemovalNotification(ctx context.Context, invitation *entity.DealInvitation, deal *entity.Deal) error {
	notification := &dto.NotificationDTO{
		Type:       enums.NotiDealMemberRemoved,
		UserId:     &invitation.InviteeID,
		AttachData: []string{deal.Name, invitation.Message},
		TargetID:   deal.ID,
		IsMerge:    false,
	}
	
	return u.notiClient.SendNotification(ctx, notification)
} 