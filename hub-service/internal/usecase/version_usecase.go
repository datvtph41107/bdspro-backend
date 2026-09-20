package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	_crud "common/domain/crud"
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	_err "common/domain/err"
	_errors "common/errors"
	"hub/internal/domain"
	iprovider "hub/internal/interface"
	_repo "hub/internal/repo"
	sharepb "pb/types/shared"
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
	NotificationClient iprovider.INotificationClient
	UserClient         iprovider.IUserClient
}

func NewVersionUsecase(
	repo _repo.IVersionRepo,
	notificationClient iprovider.INotificationClient,
	userClient iprovider.IUserClient,
) IVersionUsecase {
	return &VersionUsecase{
		BaseUsecase: _crud.BaseUsecase[domain.VersionEntity, _repo.IVersionRepo]{
			Repo: repo,
		},
		NotificationClient: notificationClient,
		UserClient:         userClient,
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

	latest, err := u.Repo.GetLatestByAppAndPlatform(ctx, appName, platform, targetVersion, true)
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

	latest, err := u.Repo.GetLatestByAppAndPlatform(ctx, bundleVersion.AppName, bundleVersion.Platform, bundleVersion.VersionName, false)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể lấy thông tin phiên bản gần nhất: " + err.Error(),
		}
	}

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

	// Gửi thông báo cho tất cả user về version mới
	go u.sendNotificationToAllUsers(ctx, bundleVersion)

	return bundleVersion, nil
}

// sendNotificationToAllUsers gửi thông báo cho tất cả user về version mới
func (u *VersionUsecase) sendNotificationToAllUsers(ctx context.Context, version *domain.VersionEntity) {
	if u.NotificationClient == nil || u.UserClient == nil {
		slog.WarnContext(ctx, strings.TrimSuffix(fmt.Sprintln("NotificationClient or UserClient is nil, skipping notification"), "\n"))
		return
	}

	// Lấy tất cả user theo batch
	page := uint32(1)
	size := uint32(100) // Lấy 100 user mỗi lần

	for {
		usersResp, err := u.UserClient.ListAllUsers(ctx, page, size)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("Failed to get users list: %v", err))
			break
		}

		if usersResp == nil || len(usersResp.Data) == 0 {
			break
		}

		// Tạo danh sách notification requests
		notificationRequests := make([]*sharepb.NotificationRequest, 0, len(usersResp.Data))
		targetId := version.ID

		for _, user := range usersResp.Data {
			title := "Có phiên bản mới"
			message := []string{
				"Ứng dụng " + version.AppName + " đã có phiên bản mới " + version.VersionName,
			}
			if version.ReleaseNotes != "" {
				message = append(message, version.ReleaseNotes)
			}

			notificationRequests = append(notificationRequests, &sharepb.NotificationRequest{
				Avatar:           "",
				Title:            title,
				Message:          message,
				NotificationType: uint32(_enum.NotificationSystem),
				OwnerId:          user.ProfileId,
				TargetId:         targetId,
				OwnerOf:          uint32(_enum.EOwnerOfMember),
				AttachData:       []string{},
				SendToDevice:     true,
				IsMerge:          true, // Merge thông báo
			})
		}

		// Gửi batch notification
		if len(notificationRequests) > 0 {
			if err := u.NotificationClient.SendBatch(ctx, notificationRequests); err != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("Failed to send batch notification: %v", err))
			} else {
				slog.InfoContext(ctx, fmt.Sprintf("Sent notification to %d users for version %s", len(notificationRequests), version.VersionName))
			}
		}

		// Kiểm tra xem còn user nào không
		if int32(len(usersResp.Data)) < int32(size) {
			break
		}

		page++
	}
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
