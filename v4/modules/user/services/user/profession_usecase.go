package services

import (
	_crud "common/domain/crud"
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"userv4/modules/user/enums"
	"userv4/modules/user/domain/user"
	"userv4/internal/interface/repo"
)

type IAdminProfessionUsecase interface {
	_crud.IBaseUsecase[user.ProfessionEntity]
	ListVerified(ctx context.Context, pagable _dto.IPagable) ([]user.ProfessionEntity, uint32, error)
}

type ProfessionUsecase struct {
	_crud.BaseUsecase[user.ProfessionEntity, repo.IProfessionRepo]
}

func NewProfessionUsecase(rp repo.IProfessionRepo) IAdminProfessionUsecase {
	return &ProfessionUsecase{
		BaseUsecase: _crud.BaseUsecase[user.ProfessionEntity, repo.IProfessionRepo]{
			Repo: rp,
		},
	}
}

func (u *ProfessionUsecase) Create(ctx context.Context, entity *user.ProfessionEntity) (*user.ProfessionEntity, error) {
	u.applyAdminMetadata(ctx, entity)
	return u.BaseUsecase.Create(ctx, entity)
}

func (u *ProfessionUsecase) Update(ctx context.Context, id uint64, entity *user.ProfessionEntity) (*user.ProfessionEntity, error) {
	u.applyAdminMetadata(ctx, entity)
	return u.BaseUsecase.Update(ctx, id, entity)
}

func (u *ProfessionUsecase) ListVerified(ctx context.Context, pagable _dto.IPagable) ([]user.ProfessionEntity, uint32, error) {
	return u.Repo.ListVerified(ctx, pagable)
}

func (u *ProfessionUsecase) applyAdminMetadata(ctx context.Context, entity *user.ProfessionEntity) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID != 0 {
		entity.ApproveUserID = &profileID
	}

	if entity.VerifiedStatus == 0 {
		entity.VerifiedStatus = enums.VerifyVerified
	}
}
