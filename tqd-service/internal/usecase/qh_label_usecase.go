package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"
)

type QHLabelUsecase interface {
	Create(ctx context.Context, label *qh_domain.QHLabel) (*qh_domain.QHLabel, error)
	Update(ctx context.Context, id uint64, label *qh_domain.QHLabel) (*qh_domain.QHLabel, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLabel, error)
	ListByLayerID(ctx context.Context, layerID uint64) ([]qh_domain.QHLabel, error)
}

type qhLabelUsecaseImpl struct {
	labelRepo repo.QHLabelRepository
}

func NewQHLabelUsecase(labelRepo repo.QHLabelRepository) QHLabelUsecase {
	return &qhLabelUsecaseImpl{labelRepo: labelRepo}
}

func (u *qhLabelUsecaseImpl) Create(ctx context.Context, label *qh_domain.QHLabel) (*qh_domain.QHLabel, error) {
	if strings.TrimSpace(label.Name) == "" {
		return nil, errors.New("tên nhãn không được để trống")
	}
	if label.LayerID == 0 {
		return nil, errors.New("layerId không được để trống")
	}
	if label.Color == "" {
		label.Color = "#1890ff"
	}

	if err := u.labelRepo.Create(ctx, label); err != nil {
		return nil, fmt.Errorf("tạo nhãn thất bại: %w", err)
	}
	return label, nil
}

func (u *qhLabelUsecaseImpl) Update(ctx context.Context, id uint64, label *qh_domain.QHLabel) (*qh_domain.QHLabel, error) {
	existing, err := u.labelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy nhãn thất bại: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("không tìm thấy nhãn id=%d", id)
	}

	if strings.TrimSpace(label.Name) != "" {
		existing.Name = label.Name
	}
	if label.Color != "" {
		existing.Color = label.Color
	}
	if label.Description != "" {
		existing.Description = label.Description
	}
	existing.SortOrder = label.SortOrder

	if err := u.labelRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("cập nhật nhãn thất bại: %w", err)
	}
	return existing, nil
}

func (u *qhLabelUsecaseImpl) Delete(ctx context.Context, id uint64) error {
	existing, err := u.labelRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("lấy nhãn thất bại: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("không tìm thấy nhãn id=%d", id)
	}
	return u.labelRepo.Delete(ctx, id)
}

func (u *qhLabelUsecaseImpl) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLabel, error) {
	label, err := u.labelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy nhãn thất bại: %w", err)
	}
	if label == nil {
		return nil, fmt.Errorf("không tìm thấy nhãn id=%d", id)
	}
	return label, nil
}

func (u *qhLabelUsecaseImpl) ListByLayerID(ctx context.Context, layerID uint64) ([]qh_domain.QHLabel, error) {
	if layerID == 0 {
		return nil, errors.New("layerId không hợp lệ")
	}
	return u.labelRepo.ListByLayerID(ctx, layerID)
}
