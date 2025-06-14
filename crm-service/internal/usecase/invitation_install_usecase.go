package usecase

import (
	_utils "common/utils"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	iusecase "crm/internal/interface"
	"crm/internal/repo"
)

type InvitationInstallUsecase struct {
	repo        repo.InvitationInstallRepo
	historyRepo repo.HistoryRepo
	userClient  iusecase.UserClient
}

func NewInvitationInstallUsecase(repo repo.InvitationInstallRepo, historyRepo repo.HistoryRepo, userClient iusecase.UserClient) *InvitationInstallUsecase {
	return &InvitationInstallUsecase{repo: repo, historyRepo: historyRepo, userClient: userClient}
}

func (u *InvitationInstallUsecase) SendInvitation(ctx context.Context, ids []uint64, req *domain.AppInvitedEntity) ([]*domain.AppInvitedEntity, error) {
	entities := make([]*domain.AppInvitedEntity, len(ids))
	for i, id := range ids {
		entities[i] = &domain.AppInvitedEntity{
			ContactID: id,
			Content:   req.Content,
			Link:      req.Link,
			Channel:   req.Channel,
		}
	}
	entities, err := u.repo.SendInvitation(ctx, entities)
	if err != nil {
		return nil, err
	}

	go func() {
		newCtx := _utils.CloneContext(ctx)
		profileId := _utils.GetProfileIdWithContext(newCtx)
		currentUser, _ := u.userClient.GetProfileById(newCtx, profileId)
		fullName := ""
		if currentUser != nil {
			fullName = currentUser.FullName
		}
		for _, entity := range entities {
			u.historyRepo.CreateHistory(newCtx, &domain.HistoryEntity{
				Title:      fullName + " đã mời bạn cài app",
				Note:       fullName + " đã gửi lời mời bạn cài app BDSPro \"" + entity.Content + "\"",
				TargetId:   entity.ContactID,
				TargetType: enums.TargetHistoryContact,
				ActionType: enums.HistoryContactInstall,
				OwnerID:    &profileId,
				OwnerType:  enums.EOwnerTypeMember,
				// AfterStage: ,
			})
		}
	}()

	return entities, nil
}

func (u *InvitationInstallUsecase) ListInvitation(ctx context.Context, req *dto.ListInvitationRequest) ([]domain.AppInvitedEntity, int64, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	entities, total, err := u.repo.ListInvitation(ctx, profileId, req)
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
