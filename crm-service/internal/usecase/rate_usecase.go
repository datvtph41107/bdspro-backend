package usecase

import (
	_enum "common/domain/enum"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/interface/provider"
	"crm/internal/repo"
	"errors"
	"fmt"
)

type RateUsecase struct {
	rateRepo             repo.RateRepo
	transaction          provider.ITransaction
	registryUsecase      *RegistryUsecase
	notificationProvider provider.NotificationProvider
}

func NewRateUsecase(rateRepo repo.RateRepo,
	transaction provider.ITransaction,
	registryUsecase *RegistryUsecase,
	notificationProvider provider.NotificationProvider) *RateUsecase {
	return &RateUsecase{
		rateRepo:             rateRepo,
		transaction:          transaction,
		registryUsecase:      registryUsecase,
		notificationProvider: notificationProvider,
	}
}

func (u *RateUsecase) CheckRateCondition(ctx context.Context, rate *domain.Rate) error {
	// todo: điều kiện đã giao dịch/ tạo lịch/ gọi điện > 30s mới đánh giá
	// todo: AI check comment vi phạm quy định cộng đồng
	message := "Bạn cần hoàn tất giao dịch hoặc lịch hẹn với người bán để có thể đánh giá."

	if false {
		return errors.New(message)
	}

	return nil
}

func (u *RateUsecase) CreateRate(ctx context.Context, registryName string, rate *domain.Rate) (*domain.Rate, error) {
	err := u.CheckRateCondition(ctx, rate)
	if err != nil {
		return nil, err
	}

	// chỉ đánh giá 1 lần cho 1 đối tượng
	profileId := _utils.GetProfileIdWithContext(ctx)
	rate.CreatedBy = &profileId

	ownerOf := u.registryUsecase.GetOwnerOfByRegistry(registryName)

	if rate.OwnerID <= 0 {
		return nil, _errors.ReturnError(400, "ownerId là bắt buộc")
	}
	if ownerOf == 0 {
		return nil, _errors.ReturnError(400, "registry không đúng")
	} else {
		rate.OwnerOf = ownerOf
	}

	rate, err = u.rateRepo.Create(ctx, rate)
	if err != nil {
		return nil, err
	}

	// Gửi thông báo cho owner (người được đánh giá) khi có đánh giá mới
	if u.notificationProvider != nil {
		go func() {
			// Tạo context mới để tránh context bị cancel
			notiCtx := _utils.CloneContext(ctx)
			ownerOfEnum := _enum.EOwnerOf(rate.OwnerOf)
			targetId := rate.ID
			title := "Bạn có đánh giá mới"
			message := []string{fmt.Sprintf("Bạn vừa nhận được một đánh giá %d sao", rate.Score)}

			// Thêm comment vào message nếu có
			if rate.Comment != "" {
				message = append(message, rate.Comment)
			}

			_ = u.notificationProvider.CreateNotification(
				notiCtx,
				"", // avatar - có thể lấy từ user service nếu cần
				title,
				message,
				_enum.NotificationRateCreate,
				&targetId,
				rate.OwnerID,
				ownerOfEnum,
				[]string{},
			)
		}()
	}

	return rate, nil
}

func (u *RateUsecase) GetRateList(ctx context.Context,
	id uint64,
	req *dto.RateSearchDTO,
) ([]domain.Rate, int64, error) {
	// Lấy ownerOf từ registry
	ownerOf := u.registryUsecase.GetOwnerOfByRegistry(req.Registry)
	if ownerOf == 0 {
		return nil, 0, errors.New("registry không tồn tại")
	}

	rates, total, err := u.rateRepo.GetRateList(ctx, id, ownerOf, req)
	if err != nil {
		return nil, 0, err
	}

	return rates, total, nil
}

func (u *RateUsecase) GetRateStats(ctx context.Context, ownerID uint64, registry string) (domain.RateStats, error) {
	// Lấy ownerOf từ registry
	ownerOf := u.registryUsecase.GetOwnerOfByRegistry(registry)
	if ownerOf == 0 {
		return domain.RateStats{}, errors.New("registry không tồn tại")
	}
	return u.rateRepo.CountInfo(ctx, ownerID, ownerOf)
}

func (u *RateUsecase) UpdateRate(ctx context.Context, id uint64, rate *domain.Rate) error {
	profileId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra quyền sở hữu
	existingRate, err := u.rateRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingRate.CreatedBy == nil || *existingRate.CreatedBy != profileId {
		return errors.New("Bạn không có quyền cập nhật đánh giá này")
	}

	return u.rateRepo.Update(ctx, id, rate)
}

func (u *RateUsecase) DeleteRate(ctx context.Context, id uint64) error {
	profileId := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra quyền sở hữu
	existingRate, err := u.rateRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingRate.CreatedBy == nil || *existingRate.CreatedBy != profileId {
		return errors.New("Bạn không có quyền xóa đánh giá này")
	}

	return u.rateRepo.Delete(ctx, id)
}

func (u *RateUsecase) GetRate(ctx context.Context, id uint64, registry string) (*domain.Rate, error) {
	// Lấy ownerOf từ registry để validate
	ownerOf := u.registryUsecase.GetOwnerOfByRegistry(registry)
	if ownerOf == 0 {
		return nil, errors.New("registry không tồn tại")
	}

	// Lấy rate và kiểm tra ownerOf
	rate, err := u.rateRepo.DetailRate(ctx, id)
	if err != nil {
		return nil, err
	}

	// Kiểm tra rate có đúng ownerOf không
	if uint32(rate.OwnerOf) != ownerOf {
		return nil, errors.New("rate không thuộc registry này")
	}

	return rate, nil
}

func (u *RateUsecase) GetHistoryRate(ctx context.Context, rateId uint64, page, size uint32) ([]domain.Rate, int64, error) {
	return u.rateRepo.GetHistoryRate(ctx, rateId, page, size)
}

func (u *RateUsecase) UpdateHidden(ctx context.Context, id uint64, hidden bool, registry string) error {
	// Lấy ownerOf từ registry để validate
	ownerOf := u.registryUsecase.GetOwnerOfByRegistry(registry)
	if ownerOf == 0 {
		return errors.New("registry không tồn tại")
	}

	// Lấy rate và kiểm tra ownerOf
	rate, err := u.rateRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Kiểm tra rate có đúng ownerOf không
	if rate.OwnerOf != ownerOf {
		return errors.New("rate không thuộc registry này")
	}

	return u.rateRepo.UpdateHidden(ctx, id, hidden)
}
