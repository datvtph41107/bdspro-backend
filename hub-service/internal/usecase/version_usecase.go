package usecase

import (
	"context"
	"strconv"
	"strings"

	_crud "common/domain/crud"
	_dto "common/domain/dto"
	_err "common/domain/err"
	_errors "common/errors"
	"hub/internal/domain"
	_repo "hub/internal/repo"
)

type IVersionUsecase interface {
	_crud.IBaseUsecase[domain.VersionEntity]
	GetListWithFilter(ctx context.Context, appName string, platform string, pagable _dto.IPagable) ([]*domain.VersionEntity, int64, *_err.ErrorDTO)
	CheckAppUpdate(ctx context.Context, appName string, platform string, targetVersion string) (*domain.VersionEntity, bool, bool, *_err.ErrorDTO)
	CreateBundleVersion(ctx context.Context, bundleVersion *domain.VersionEntity) (*domain.VersionEntity, *_err.ErrorDTO)
	UpdateBundleVersion(ctx context.Context, bundleVersion *domain.VersionEntity) (*domain.VersionEntity, error)
}

type VersionUsecase struct {
	_crud.BaseUsecase[domain.VersionEntity, _repo.IVersionRepo]
}

func NewVersionUsecase(repo _repo.IVersionRepo) IVersionUsecase {
	return &VersionUsecase{
		BaseUsecase: _crud.BaseUsecase[domain.VersionEntity, _repo.IVersionRepo]{
			Repo: repo,
		},
	}
}

func (u *VersionUsecase) GetListWithFilter(ctx context.Context, appName string, platform string, pagable _dto.IPagable) ([]*domain.VersionEntity, int64, *_err.ErrorDTO) {
	data, total, err := u.Repo.GetListWithFilter(ctx, appName, platform, pagable)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể lấy danh sách phiên bản: " + err.Error(),
		}
	}

	return data, total, nil
}

func (u *VersionUsecase) CheckAppUpdate(ctx context.Context, appName string, platform string, targetVersion string) (*domain.VersionEntity, bool, bool, *_err.ErrorDTO) {
	if appName == "" {
		return nil, false, false, &_err.ErrorDTO{
			Code:    400,
			Message: "Thiếu appName",
		}
	}

	if platform == "" {
		return nil, false, false, &_err.ErrorDTO{
			Code:    400,
			Message: "Thiếu platform",
		}
	}

	latest, err := u.Repo.GetLatestByAppAndPlatform(ctx, appName, platform)
	if err != nil {
		return nil, false, false, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể lấy thông tin phiên bản: " + err.Error(),
		}
	}

	if latest == nil {
		return nil, false, false, nil
	}

	needUpdate := false
	if targetVersion == "" {
		needUpdate = true
	} else if compareVersionStrings(targetVersion, latest.VersionName) < 0 {
		needUpdate = true
	}

	return latest, needUpdate, latest.ForceUpdate, nil
}

func compareVersionStrings(current string, latest string) int {
	currentParts := parseVersionParts(current)
	latestParts := parseVersionParts(latest)

	maxLen := len(currentParts)
	if len(latestParts) > maxLen {
		maxLen = len(latestParts)
	}

	for i := 0; i < maxLen; i++ {
		var curVal, latVal int
		if i < len(currentParts) {
			curVal = currentParts[i]
		}
		if i < len(latestParts) {
			latVal = latestParts[i]
		}
		if curVal < latVal {
			return -1
		}
		if curVal > latVal {
			return 1
		}
	}

	return 0
}

func parseVersionParts(version string) []int {
	segments := strings.Split(version, ".")
	parts := make([]int, 0, len(segments))
	for _, segment := range segments {
		if segment == "" {
			parts = append(parts, 0)
			continue
		}
		value, err := strconv.Atoi(segment)
		if err != nil {
			parts = append(parts, 0)
			continue
		}
		parts = append(parts, value)
	}
	return parts
}

func (u *VersionUsecase) CreateBundleVersion(ctx context.Context, bundleVersion *domain.VersionEntity) (*domain.VersionEntity, *_err.ErrorDTO) {
	if bundleVersion.AppName == "" {
		return nil, &_err.ErrorDTO{
			Code:    400,
			Message: "Thiếu appName",
		}
	}

	if bundleVersion.Platform == "" {
		return nil, &_err.ErrorDTO{
			Code:    400,
			Message: "Thiếu platform",
		}
	}

	if bundleVersion.VersionName == "" {
		return nil, &_err.ErrorDTO{
			Code:    400,
			Message: "Thiếu versionName",
		}
	}

	latest, _ := u.Repo.GetLatestByAppAndPlatform(ctx, bundleVersion.AppName, bundleVersion.Platform)
	// if err != nil {
	// 	return nil, &_err.ErrorDTO{
	// 		Code:    500,
	// 		Message: "Không thể lấy thông tin phiên bản gần nhất: " + err.Error(),
	// 	}
	// }

	if latest == nil {
		bundleVersion.BuildNumber = 1
	} else {
		bundleVersion.BuildNumber = latest.BuildNumber + 1
	}

	if err := u.Repo.Create(ctx, bundleVersion); err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể tạo phiên bản bundle: " + err.Error(),
		}
	}
	return bundleVersion, nil
}

func (u *VersionUsecase) UpdateBundleVersion(ctx context.Context, bundleVersion *domain.VersionEntity) (*domain.VersionEntity, error) {
	if bundleVersion.ID == 0 {
		return nil, _errors.ReturnError(400, "Thiếu ID")
	}

	oldBundleVersion, err := u.Repo.GetByID(ctx, bundleVersion.ID)
	if err != nil {
		return nil, _errors.ReturnError(500, "Không thể lấy thông tin phiên bản bundle: "+err.Error())
	}

	oldBundleVersion.ForceUpdate = bundleVersion.ForceUpdate
	oldBundleVersion.Active = bundleVersion.Active
	oldBundleVersion.ReleaseNotes = bundleVersion.ReleaseNotes
	print("oldBundleVersion", oldBundleVersion.Active)

	if err := u.Repo.Update(ctx, oldBundleVersion.ID, oldBundleVersion); err != nil {
		return nil, _errors.ReturnError(500, "Không thể cập nhật phiên bản bundle: "+err.Error())
	}
	return bundleVersion, nil
}
