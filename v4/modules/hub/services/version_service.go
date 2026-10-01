package services

import (
	"context"
	"strconv"
	"strings"

	_errors "common/errors"
	hubpb "pb/types/hub"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/entity"
	"userv4/modules/hub/repository"
)

// VersionService owns app-bundle version flows for modules/hub.
type VersionService interface {
	CheckAppUpdate(ctx context.Context, appName string, req dto.CheckAppUpdateRequest) (*hubpb.CheckAppUpdateResponse, error)
	CreateVersionFromDev(ctx context.Context, req *hubpb.CreateVersionRequest) (*hubpb.CreateVersionResponse, error)
	AdminList(ctx context.Context, appName, platform string, page, size int) ([]*hubpb.VersionItem, int64, error)
	AdminGet(ctx context.Context, id uint64) (*hubpb.VersionItem, error)
	AdminCreate(ctx context.Context, req dto.CreateVersionAdminRequest) (*hubpb.VersionItem, error)
	AdminUpdate(ctx context.Context, id uint64, req dto.CreateVersionAdminRequest) (*hubpb.VersionItem, error)
	AdminDelete(ctx context.Context, id uint64) error
}

type LocalVersionService struct {
	repo repository.VersionRepository
}

func NewLocalVersionService(repo repository.VersionRepository) *LocalVersionService {
	return &LocalVersionService{repo: repo}
}

func (s *LocalVersionService) CheckAppUpdate(ctx context.Context, appName string, req dto.CheckAppUpdateRequest) (*hubpb.CheckAppUpdateResponse, error) {
	if appName == "" {
		return nil, _errors.BadRequestException("Thiếu appName")
	}
	if req.Platform == "" {
		return nil, _errors.BadRequestException("Thiếu platform")
	}
	latest, err := s.repo.GetLatest(ctx, appName, req.Platform, true)
	if err != nil {
		return nil, err
	}
	if latest == nil {
		return &hubpb.CheckAppUpdateResponse{}, nil
	}
	needUpdate := req.TargetVersion == "" || compareVersionStrings(req.TargetVersion, latest.VersionName) < 0
	return &hubpb.CheckAppUpdateResponse{
		NeedUpdate:    needUpdate,
		ForceUpdate:   needUpdate && latest.ForceUpdate,
		LatestVersion: mapVersionItem(latest),
	}, nil
}

func (s *LocalVersionService) CreateVersionFromDev(ctx context.Context, req *hubpb.CreateVersionRequest) (*hubpb.CreateVersionResponse, error) {
	if req == nil {
		return nil, _errors.BadRequestException("body is required")
	}
	apiKey, _ := ctx.Value(apiKeyContextKey{}).(string)
	if err := s.repo.ValidateAPIKey(ctx, apiKey); err != nil {
		return nil, _errors.UnauthorizedException(err.Error())
	}
	v := &entity.Version{
		AppName:      req.GetAppName(),
		Platform:     req.GetPlatform(),
		VersionName:  req.GetVersionName(),
		BuildNumber:  req.GetBuildNumber(),
		ForceUpdate:  req.GetForceUpdate(),
		Active:       req.GetActive(),
		DownloadURL:  req.GetDownloadUrl(),
		ReleaseNotes: req.GetReleaseNotes(),
		BundleName:   req.GetBundleName(),
		BundleSize:   req.GetBundleSize(),
	}
	if err := s.repo.Create(ctx, v); err != nil {
		return nil, err
	}
	return &hubpb.CreateVersionResponse{Data: mapVersionItem(v)}, nil
}

func (s *LocalVersionService) AdminList(ctx context.Context, appName, platform string, page, size int) ([]*hubpb.VersionItem, int64, error) {
	rows, total, err := s.repo.List(ctx, appName, platform, page, size)
	if err != nil {
		return nil, 0, err
	}
	items := make([]*hubpb.VersionItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, mapVersionItem(r))
	}
	return items, total, nil
}

func (s *LocalVersionService) AdminGet(ctx context.Context, id uint64) (*hubpb.VersionItem, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, _errors.NotFoundException("version not found")
	}
	return mapVersionItem(row), nil
}

func (s *LocalVersionService) AdminCreate(ctx context.Context, req dto.CreateVersionAdminRequest) (*hubpb.VersionItem, error) {
	if req.AppName == "" || req.Platform == "" || req.VersionName == "" {
		return nil, _errors.BadRequestException("appName, platform and versionName are required")
	}
	v := &entity.Version{
		AppName: req.AppName, Platform: req.Platform, VersionName: req.VersionName,
		BuildNumber: req.BuildNumber, ForceUpdate: req.ForceUpdate, Active: req.Active,
		DownloadURL: req.DownloadURL, ReleaseNotes: req.ReleaseNotes,
		BundleName: req.BundleName, BundleSize: req.BundleSize,
	}
	if err := s.repo.Create(ctx, v); err != nil {
		return nil, err
	}
	return mapVersionItem(v), nil
}

func (s *LocalVersionService) AdminUpdate(ctx context.Context, id uint64, req dto.CreateVersionAdminRequest) (*hubpb.VersionItem, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, _errors.NotFoundException("version not found")
	}
	existing.AppName = req.AppName
	existing.Platform = req.Platform
	existing.VersionName = req.VersionName
	existing.BuildNumber = req.BuildNumber
	existing.ForceUpdate = req.ForceUpdate
	existing.Active = req.Active
	existing.DownloadURL = req.DownloadURL
	existing.ReleaseNotes = req.ReleaseNotes
	existing.BundleName = req.BundleName
	existing.BundleSize = req.BundleSize
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return mapVersionItem(existing), nil
}

func (s *LocalVersionService) AdminDelete(ctx context.Context, id uint64) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return _errors.NotFoundException("version not found")
	}
	return s.repo.Delete(ctx, id)
}

type apiKeyContextKey struct{}

// WithAPIKey attaches an HTTP API key into context for CreateVersionFromDev.
func WithAPIKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, apiKeyContextKey{}, strings.TrimSpace(key))
}

func mapVersionItem(v *entity.Version) *hubpb.VersionItem {
	if v == nil {
		return nil
	}
	out := &hubpb.VersionItem{
		Id: v.ID, AppName: v.AppName, Platform: v.Platform, VersionName: v.VersionName,
		BuildNumber: v.BuildNumber, ForceUpdate: v.ForceUpdate, Active: v.Active,
		DownloadUrl: v.DownloadURL, ReleaseNotes: v.ReleaseNotes,
		BundleName: v.BundleName, BundleSize: v.BundleSize,
	}
	if v.CreatedAt != nil {
		out.CreatedAt = strconv.FormatInt(v.CreatedAt.UnixMilli(), 10)
	}
	if v.UpdatedAt != nil {
		out.UpdatedAt = strconv.FormatInt(v.UpdatedAt.UnixMilli(), 10)
	}
	return out
}

func compareVersionStrings(current, latest string) int {
	currentParts := parseVersionParts(current)
	latestParts := parseVersionParts(latest)
	maxLen := len(currentParts)
	if len(latestParts) > maxLen {
		maxLen = len(latestParts)
	}
	for i := 0; i < maxLen; i++ {
		var c, l int
		if i < len(currentParts) {
			c = currentParts[i]
		}
		if i < len(latestParts) {
			l = latestParts[i]
		}
		if c < l {
			return -1
		}
		if c > l {
			return 1
		}
	}
	return 0
}

func parseVersionParts(v string) []int {
	parts := strings.Split(strings.TrimSpace(v), ".")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		n, _ := strconv.Atoi(p)
		out = append(out, n)
	}
	return out
}
