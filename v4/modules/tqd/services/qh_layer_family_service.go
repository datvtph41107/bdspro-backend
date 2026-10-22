package services

import (
	"common/validation"
	"context"
	"errors"
	"fmt"
	"strings"

	_dto "common/domain/dto"

	qh "mdn4/modules/tqd/entity/qh"
	repo "mdn4/modules/tqd/repository"
)

// ErrLayerFamilyNotFound chỉ biểu thị lookup thành công nhưng không có bản ghi.
var ErrLayerFamilyNotFound = errors.New("layer family not found")

// QHLayerFamilyUpdateInput — trường nil = không đổi.
type QHLayerFamilyUpdateInput struct {
	Name       *string
	SortNumber *int32
}

type QHLayerFamilyUsecase interface {
	Create(ctx context.Context, row *qh.QHLayerFamily) (*qh.QHLayerFamily, error)
	Update(ctx context.Context, id uint64, in *QHLayerFamilyUpdateInput) (*qh.QHLayerFamily, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh.QHLayerFamily, error)
	List(ctx context.Context, pagable *_dto.Pagable, search string) ([]qh.QHLayerFamily, int64, error)
	ListClient(ctx context.Context, pagable *_dto.Pagable, search string) ([]qh.QHLayerFamily, int64, error)
	GetZoomRangeByFamilyID(ctx context.Context, familyID uint64) (uint32, uint32, error)
	BuildFamilyPMTiles(ctx context.Context, familyID uint64, minZoom, maxZoom uint32, outputDir string, progressFn func(BuildProgress)) error
}

type qhLayerFamilyUsecase struct {
	repo         repo.QHLayerFamilyRepository
	layerRepo    repo.LayerRepository
	regionRepo   repo.RegionRepository
	buildTracker *buildTracker
}

func NewQHLayerFamilyUsecase(
	r repo.QHLayerFamilyRepository,
	layerRepo repo.LayerRepository,
	regionRepo repo.RegionRepository,
) QHLayerFamilyUsecase {
	return &qhLayerFamilyUsecase{
		repo:         r,
		layerRepo:    layerRepo,
		regionRepo:   regionRepo,
		buildTracker: newBuildTracker(),
	}
}

func (u *qhLayerFamilyUsecase) Create(ctx context.Context, row *qh.QHLayerFamily) (*qh.QHLayerFamily, error) {
	if row == nil {
		return nil, validation.New("payload", "LAYER_FAMILY_PAYLOAD_INVALID", "payload is required")
	}
	row.Name = strings.TrimSpace(row.Name)
	if row.Name == "" {
		return nil, validation.New("name", "LAYER_FAMILY_NAME_INVALID", "name is required")
	}
	if err := u.repo.Create(ctx, row); err != nil {
		return nil, fmt.Errorf("create layer family: %w", err)
	}
	return row, nil
}

func (u *qhLayerFamilyUsecase) Update(ctx context.Context, id uint64, in *QHLayerFamilyUpdateInput) (*qh.QHLayerFamily, error) {
	if id == 0 {
		return nil, validation.New("id", "LAYER_FAMILY_ID_INVALID", "id is required")
	}
	if in == nil {
		return nil, validation.New("payload", "LAYER_FAMILY_PAYLOAD_INVALID", "payload is required")
	}
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("%w: id=%d", ErrLayerFamilyNotFound, id)
	}

	if in.Name != nil {
		s := strings.TrimSpace(*in.Name)
		if s == "" {
			return nil, validation.New("name", "LAYER_FAMILY_NAME_INVALID", "name cannot be empty")
		}
		existing.Name = s
	}

	if in.SortNumber != nil {
		existing.SortNumber = *in.SortNumber
	}

	if err := u.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("update layer family: %w", err)
	}
	return existing, nil
}

func (u *qhLayerFamilyUsecase) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return validation.New("id", "LAYER_FAMILY_ID_INVALID", "id is required")
	}
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get by id: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("%w: id=%d", ErrLayerFamilyNotFound, id)
	}
	return u.repo.Delete(ctx, id)
}

func (u *qhLayerFamilyUsecase) GetByID(ctx context.Context, id uint64) (*qh.QHLayerFamily, error) {
	if id == 0 {
		return nil, validation.New("id", "LAYER_FAMILY_ID_INVALID", "id is required")
	}
	row, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	if row == nil {
		return nil, fmt.Errorf("%w: id=%d", ErrLayerFamilyNotFound, id)
	}
	return row, nil
}

func (u *qhLayerFamilyUsecase) List(ctx context.Context, pagable *_dto.Pagable, search string) ([]qh.QHLayerFamily, int64, error) {
	offset := 0
	limit := 20
	sort := ""
	if pagable != nil {
		offset = pagable.GetOffset()
		limit = pagable.GetLimit()
		sort = pagable.GetSort()
	}
	return u.repo.List(ctx, offset, limit, search, sort)
}

func (u *qhLayerFamilyUsecase) ListClient(ctx context.Context, pagable *_dto.Pagable, search string) ([]qh.QHLayerFamily, int64, error) {
	offset := 0
	limit := 20
	if pagable != nil {
		offset = pagable.GetOffset()
		limit = pagable.GetLimit()
	}
	return u.repo.ListClient(ctx, offset, limit, search)
}

func (u *qhLayerFamilyUsecase) GetZoomRangeByFamilyID(ctx context.Context, familyID uint64) (uint32, uint32, error) {
	return u.layerRepo.GetZoomRangeByFamilyID(ctx, familyID)
}

func (u *qhLayerFamilyUsecase) BuildFamilyPMTiles(ctx context.Context, familyID uint64, minZoom, maxZoom uint32, outputDir string, progressFn func(BuildProgress)) error {
	buildCtx, gen, err := u.buildTracker.tryStart(familyID)
	if err != nil {
		if buildErr, ok := err.(*ErrBuildInProgress); ok {
			return &ErrFamilyBuildInProgress{FamilyID: familyID, Progress: buildErr.Progress}
		}
		return err
	}
	defer u.buildTracker.finish(familyID, gen)

	family, err := u.repo.GetByID(ctx, familyID)
	if err != nil {
		return fmt.Errorf("get family for build: %w", err)
	}
	if family == nil {
		return fmt.Errorf("%w: id=%d", ErrLayerFamilyNotFound, familyID)
	}

	familyMinZoom, familyMaxZoom, err := u.layerRepo.GetZoomRangeByFamilyID(ctx, familyID)
	if err != nil {
		return fmt.Errorf("get zoom range for family %d: %w", familyID, err)
	}
	if minZoom == 0 {
		minZoom = familyMinZoom
	}
	if maxZoom == 0 {
		maxZoom = familyMaxZoom
	}
	if outputDir == "" {
		outputDir = "./files/tiles"
	}

	wrappedFn := func(p BuildProgress) {
		u.buildTracker.updateProgress(familyID, gen, p)
		progressFn(p)
	}

	return BuildFamilyPMTilesFile(buildCtx, u.regionRepo, familyID, minZoom, maxZoom, outputDir, wrappedFn)
}
