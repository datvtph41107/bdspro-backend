package usecase

import (
	"context"
	"errors"
	"regexp"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"
)

type OneHouseUsecase interface {
	GetDetail(ctx context.Context, idOrUUID uint64) (*qh_domain.OneHouse, error)
}

type oneHouseUsecase struct {
	oneHouseRepo repo.OneHouseRepo
}

func NewOneHouseUsecase(repo repo.OneHouseRepo) OneHouseUsecase {
	return &oneHouseUsecase{oneHouseRepo: repo}
}

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func isUUID(s string) bool {
	return uuidRegex.MatchString(s)
}

func (u *oneHouseUsecase) GetDetail(ctx context.Context, idOrUUID uint64) (*qh_domain.OneHouse, error) {
	if idOrUUID == 0 {
		return nil, errors.New("id or uuid is required")
	}

	// Ưu tiên tìm theo property_uuid nếu đầu vào có dạng UUID
	// if isUUID(idOrUUID) {
	// 	oh, err := u.oneHouseRepo.GetByPropertyUUID(ctx, idOrUUID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if oh != nil {
	// 		return oh, nil
	// 	}
	// 	// Nếu không tìm thấy theo UUID, fallback tìm theo id (text)
	// 	oh, err = u.oneHouseRepo.GetByID(ctx, idOrUUID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if oh != nil {
	// 		return oh, nil
	// 	}
	// 	return nil, errors.New("one_house not found")
	// }

	// Nếu không phải UUID, chỉ tìm theo id (text)
	oh, err := u.oneHouseRepo.GetByID(ctx, idOrUUID)
	if err != nil {
		return nil, err
	}
	if oh == nil {
		return nil, errors.New("one_house not found")
	}
	return oh, nil
}
