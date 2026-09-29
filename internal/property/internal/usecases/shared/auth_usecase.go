package shared_usecase

import (
	"bdspro/internal/enums"
	_utils "common/utils"
	"context"
)

type RequestUsecase struct {
}

func NewRequestUsecase() *RequestUsecase {
	return &RequestUsecase{}
}

func (uc *RequestUsecase) GetOwnerIDAndType(c context.Context) (*uint64, enums.EOwnerOf, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	return &profileId, enums.EOwnerOfMember, nil
}
