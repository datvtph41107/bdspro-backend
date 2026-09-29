package usecases

import (
	_crud "common/domain/crud"
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"user/enums"
	"user/internal/interface/repo"
	"user/internal/models"
)

type IAdminProfessionUsecase interface {
	_crud.IBaseUsecase[models.ProfessionEntity]
	ListVerified(ctx context.Context, pagable _dto.IPagable) ([]models.ProfessionEntity, uint32, error)
}

type ProfessionUsecase struct {
	_crud.BaseUsecase[models.ProfessionEntity, repo.IProfessionRepo]
}

func NewProfessionUsecase(rp repo.IProfessionRepo) IAdminProfessionUsecase {
	return &ProfessionUsecase{
		BaseUsecase: _crud.BaseUsecase[models.ProfessionEntity, repo.IProfessionRepo]{
			Repo: rp,
		},
	}
}

func (u *ProfessionUsecase) Create(ctx context.Context, entity *models.ProfessionEntity) (*models.ProfessionEntity, error) {
	u.applyAdminMetadata(ctx, entity)
	return u.BaseUsecase.Create(ctx, entity)
}

func (u *ProfessionUsecase) Update(ctx context.Context, id uint64, entity *models.ProfessionEntity) (*models.ProfessionEntity, error) {
	u.applyAdminMetadata(ctx, entity)
	return u.BaseUsecase.Update(ctx, id, entity)
}

func (u *ProfessionUsecase) ListVerified(ctx context.Context, pagable _dto.IPagable) ([]models.ProfessionEntity, uint32, error) {
	return u.Repo.ListVerified(ctx, pagable)
}

func (u *ProfessionUsecase) applyAdminMetadata(ctx context.Context, entity *models.ProfessionEntity) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID != 0 {
		entity.ApproveUserID = &profileID
	}

	if entity.VerifiedStatus == 0 {
		entity.VerifiedStatus = enums.VerifyVerified
	}
}
