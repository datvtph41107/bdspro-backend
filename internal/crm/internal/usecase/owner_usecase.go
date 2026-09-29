package usecase

import (
	base_enum "base/enum"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"crm/internal"
)

type OwnerUsecase struct {
}

func NewOwnerUsecase() *OwnerUsecase {
	return &OwnerUsecase{}
}

func (u *OwnerUsecase) GetOwnerInfo(ctx context.Context, ownerOf base_enum.EOwnerOf, ownerId uint64) (uint64, error) {
	id := uint64(0)
	if !ownerOf.IsValid() {
		return 0, _errors.ReturnError(service.OwnerTypeInvalid)
	}

	if ownerOf == base_enum.EOwnerOfOrgnization {
		id = _utils.GetOrganizationIdFromContext(ctx)
	}

	if ownerOf == base_enum.EOwnerOfMember {
		id = _utils.GetProfileIdWithContext(ctx)
	}

	if ownerOf == base_enum.EOwnerOfGroup {
		id = ownerId
	}

	if ownerOf == base_enum.EOwnerOfAdmin {
		id = _utils.GetAuthIdFromContext(ctx)
		if id == 0 {
			id = _utils.GetProfileIdWithContext(ctx)
		}
	}

	if id == 0 {
		return 0, _errors.ReturnError(service.RequestCannotBePerformed)
	}

	return id, nil
}
