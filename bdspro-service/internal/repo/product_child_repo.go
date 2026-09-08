package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

type UProductChildRepo interface {
	GetByID(ctx context.Context, id uint64) (*domain.Product, error)
	InfoArea(ctx context.Context, profileId, id *uint64) (*dto.InfoAreaResponse, error)
	TotalAreaDevide(parentID *uint64) float64
	ValidatorMerge(parentID uint64, childIds []uint64) bool
	MergeChilds(childIds []uint64) error
	ValidListChild(parentID uint64) bool
	AllChildId(parentId uint64) []uint64
	MergeAll(parentId uint64) error
	ValidatorDevide(parentID uint64) bool
	// CreateHistory(
	// 	c context.Context,
	// 	action enums.EProductHistory,
	// 	productID *uint64,
	// 	parentID uint64,
	// 	childIDs []uint64,
	// 	areas []float64,
	// ) error

	// History(
	// 	c context.Context,
	// 	profileId uint64,
	// 	id uint64,
	// 	dto *dto.ProductHistorySearch) (*[]domain.ProductHistory, int64, error)
}
