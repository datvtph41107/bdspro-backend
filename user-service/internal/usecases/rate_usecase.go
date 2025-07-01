package usecases

import (
	_dto "common/dto"
	_utils "common/utils"
	"context"
	"errors"
	"time"
	"user/internal/dto"
	iusecase "user/internal/interface"
	user_model "user/internal/models/user"
	"user/internal/repo"
)

type RateUsecase struct {
	rateRepo    repo.RateRepo
	transaction iusecase.ITransaction
}

func NewRateUsecase(rateRepo repo.RateRepo, transaction iusecase.ITransaction) *RateUsecase {
	return &RateUsecase{
		rateRepo:    rateRepo,
		transaction: transaction,
	}
}

func (u *RateUsecase) CheckRateCondition(ctx context.Context, rate *user_model.Rate) error {
	// todo: điều kiện đã giao dịch/ tạo lịch/ gọi điện > 30s mới đánh giá
	// todo: AI check comment vi phạm quy định cộng đồng
	message := "Bạn cần hoàn tất giao dịch hoặc lịch hẹn với người bán để có thể đánh giá."

	if false {
		return errors.New(message)
	}

	return nil
}

func (u *RateUsecase) CreateRate(ctx context.Context, rate *user_model.Rate) (*user_model.Rate, error) {
	err := u.CheckRateCondition(ctx, rate)
	if err != nil {
		return nil, err
	}

	// chỉ đánh giá 1 lần cho 1 đối tượng
	profileId := _utils.GetProfileIdWithContext(ctx)
	existRate, err := u.rateRepo.GetByProfileIdAndTargetIdAndTargetType(ctx,
		profileId,
		rate.TargetID,
		rate.TargetType,
	)
	if err != nil {
		return nil, err
	}
	if existRate != nil {
		return nil, errors.New("Bạn đã đánh giá cho người này rồi")
	}

	rate, err = u.rateRepo.CreateRate(ctx, rate)
	if err != nil {
		return nil, err
	}
	return rate, nil
}

func (u *RateUsecase) GetRateList(ctx context.Context,
	id uint64,
	targetType uint32,
	req *dto.RateSearchDTO,
) ([]user_model.RateItem, int64, error) {
	rates, total, err := u.rateRepo.GetRateList(ctx, id, uint8(targetType), req)
	if err != nil {
		return nil, 0, err
	}

	return rates, total, nil
}

func (u *RateUsecase) UpdateRate(ctx context.Context, req *user_model.Rate) (*user_model.Rate, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)

	rate, err := u.rateRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if time.Now().After(rate.CreatedAt.Add(time.Hour * 24)) {
		return nil, errors.New("Chỉ có thể cập nhật đánh giá trong vòng 24 giờ sau khi tạo")
	}

	if rate.CreatedBy != nil && *rate.CreatedBy != profileId {
		return nil, errors.New("Bạn không có quyền cập nhật đánh giá này")
	}

	err = u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		rate, err = u.rateRepo.UpdateRate(ctx, req, rate)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rate, nil
}

func (u *RateUsecase) DeleteRate(ctx context.Context, id uint64) error {
	profileId := _utils.GetProfileIdWithContext(ctx)

	rate, err := u.rateRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if rate.CreatedBy != nil && *rate.CreatedBy != profileId {
		return errors.New("you are not allowed to delete this rate")
	}

	return u.rateRepo.DeleteRate(ctx, rate.ID)
}

func (u *RateUsecase) GetHistoryRate(ctx context.Context, rateId uint64, page uint32, size uint32) ([]user_model.RateItem, int64, error) {
	dto := &_dto.Pagable{
		Page: int(page),
		Size: int(size),
	}

	rates, total, err := u.rateRepo.GetHistoryRate(ctx, rateId, dto)
	if err != nil {
		return nil, 0, err
	}
	return rates, total, nil
}

func (u *RateUsecase) UpdateHidden(ctx context.Context, id uint64, hidden bool) error {
	// todo: check admin ở đây

	return u.rateRepo.UpdateHidden(ctx, id, hidden)
}
