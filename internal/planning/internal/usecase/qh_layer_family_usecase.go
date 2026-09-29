package usecase

import (
	"context"
	"fmt"
	"strings"

	_dto "common/domain/dto"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"
)

// QHLayerFamilyUpdateInput — trường nil = không đổi.
type QHLayerFamilyUpdateInput struct {
	Name       *string
	SortNumber *int32
}

type QHLayerFamilyUsecase interface {
	Create(ctx context.Context, row *qh_domain.QHLayerFamily) (*qh_domain.QHLayerFamily, error)
	Update(ctx context.Context, id uint64, in *QHLayerFamilyUpdateInput) (*qh_domain.QHLayerFamily, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayerFamily, error)
	List(ctx context.Context, pagable *_dto.Pagable, search string) ([]qh_domain.QHLayerFamily, int64, error)
	ListClient(ctx context.Context, pagable *_dto.Pagable, search string) ([]qh_domain.QHLayerFamily, int64, error)
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

func (u *qhLayerFamilyUsecase) Create(ctx context.Context, row *qh_domain.QHLayerFamily) (*qh_domain.QHLayerFamily, error) {
	if row == nil {
		return nil, ErrQHLayerFamilyPayloadRequired
	}
	row.Name = strings.TrimSpace(row.Name)
	if row.Name == "" {
		return nil, ErrQHLayerFamilyNameRequired
	}
	if err := u.repo.Create(ctx, row); err != nil {
		return nil, fmt.Errorf("create layer family: %w", err)
	}
	return row, nil
}

func (u *qhLayerFamilyUsecase) Update(ctx context.Context, id uint64, in *QHLayerFamilyUpdateInput) (*qh_domain.QHLayerFamily, error) {
	if id == 0 {
		return nil, ErrQHLayerFamilyIDRequired
	}
	if in == nil {
		return nil, ErrQHLayerFamilyPayloadRequired
	}
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	if existing == nil {
		return nil, qhLayerFamilyNotFound(id)
	}

	if in.Name != nil {
		s := strings.TrimSpace(*in.Name)
		if s == "" {
			return nil, ErrQHLayerFamilyNameRequired
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
		return ErrQHLayerFamilyIDRequired
	}
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get by id: %w", err)
	}
	if existing == nil {
		return qhLayerFamilyNotFound(id)
	}
	return u.repo.Delete(ctx, id)
}

func (u *qhLayerFamilyUsecase) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayerFamily, error) {
	if id == 0 {
		return nil, ErrQHLayerFamilyIDRequired
	}
	row, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	if row == nil {
		return nil, qhLayerFamilyNotFound(id)
	}
	return row, nil
}

func (u *qhLayerFamilyUsecase) List(ctx context.Context, pagable *_dto.Pagable, search string) ([]qh_domain.QHLayerFamily, int64, error) {
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

func (u *qhLayerFamilyUsecase) ListClient(ctx context.Context, pagable *_dto.Pagable, search string) ([]qh_domain.QHLayerFamily, int64, error) {
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
		return fmt.Errorf("get family %d: %w", familyID, err)
	}
	if family == nil {
		return qhLayerFamilyNotFound(familyID)
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
