package usecase

import (
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"hub/internal"
	"log/slog"
	"strconv"
	"strings"
	"time"

	providers "hub/internal/provider"
	_repo "hub/internal/repo"
)

// IUpdateDataUsecase interface cho UpdateDataService (sync tracking)
type IUpdateDataUsecase interface {
	CheckVersionSync(ctx context.Context, resource string, lastSync int64, limit int64) ([]uint64, error)
	CheckVersionSyncById(ctx context.Context, resource string, lastSync int64, id uint64) (int64, error)
	RefreshSyncIds(ctx context.Context, ownerID uint64, resource string, lastSync int64, limit int64) (needSync bool, currentVersion int64, changeCount int32, changedIds []uint64, next *int64, err error)
	FlushSyncIds(ctx context.Context, ownerID uint64, resource string, limit int64) error
	Put(ctx context.Context, ownerIDs []uint64, resource string, resourceId uint64) error
}

type UpdateDataUsecase struct {
	repo           _repo.IUpdateDataRepo
	cacheProvider  providers.CacheProvider
	updateProvider providers.UpdateDataProvider
}

const (
	// lưu danh sách id cần update
	keyUserUpdate = "u:%s:%d"

	// lưu update time của entity
	keyResourceUpdate = "time:%s:%d"
)

func NewUpdateDataUsecase(repo _repo.IUpdateDataRepo,
	cacheProvider providers.CacheProvider,
	updateProvider providers.UpdateDataProvider) IUpdateDataUsecase {
	return &UpdateDataUsecase{
		repo:           repo,
		cacheProvider:  cacheProvider,
		updateProvider: updateProvider,
	}
}

func (uc *UpdateDataUsecase) CheckVersionSync(ctx context.Context, resource string, lastSync int64, limit int64) ([]uint64, error) {
	ownerId := _utils.GetOriginIdFromContext(ctx)
	key := fmt.Sprintf(keyUserUpdate, resource, ownerId)
	// client, err := u.cacheProvider.Client(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("get user update data: %w", err)
	// }
	resourceIdStrs, err := uc.cacheProvider.LRange(ctx, key, 0, -1)
	if err != nil {
		return nil, fmt.Errorf("get user update data: %w", err)
	}
	// if len(productIds) == 0 {
	// 	return nil, _errors.NotFoundException("User products not found")
	// }
	resourceIds := make([]uint64, len(resourceIdStrs))
	for i, resourceIdStr := range resourceIdStrs {
		resourceIds[i], err = strconv.ParseUint(resourceIdStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse resource id %q: %w", resourceIdStr, err)
		}
	}
	return resourceIds, nil
}

func (uc *UpdateDataUsecase) CheckVersionSyncById(ctx context.Context, resource string, lastSync int64, id uint64) (int64, error) {
	key := fmt.Sprintf(keyResourceUpdate, resource, id)
	// client, err := u.cacheProvider.Client(ctx)
	// if err != nil {
	// 	return nil, _errors.InternalServerException("get user update data error: %w", err.Error())
	// }
	lastUpdateStr, err := uc.cacheProvider.Get(ctx, key)
	if err != nil {
		return -1, fmt.Errorf("get resource update data: %w", err)
	}
	lastUD := _utils.ParseInt64(lastUpdateStr)
	if lastUD < 1 {
		timestamp, err := uc.updateProvider.GetUpdatedAtOfId(ctx, resource, id)
		if err != nil {
			return -1, err
		}
		uc.cacheProvider.Set(ctx, key, strconv.FormatInt(timestamp, 10))
		return timestamp, nil
	}

	return lastUD, nil
}

func (uc *UpdateDataUsecase) RefreshSyncIds(ctx context.Context, ownerID uint64, resource string, lastSync int64, limit int64) (bool, int64, int32, []uint64, *int64, error) {
	changedIds, err := uc.updateProvider.GetRefreshIds(ctx, resource, ownerID)
	if err != nil {
		return false, 0, 0, nil, nil, err
	}

	needSync := len(changedIds) > 0
	currentVersion := time.Now().Unix()
	return needSync, currentVersion, int32(len(changedIds)), changedIds, nil, nil
}

func (uc *UpdateDataUsecase) FlushSyncIds(ctx context.Context, ownerID uint64, resource string, limit int64) error {
	ownerId := _utils.GetOriginIdFromContext(ctx)
	if ownerId == 0 {
		return _errors.ReturnError(service.OwnerIDRequired)
	}
	key := fmt.Sprintf(keyUserUpdate, resource, ownerId)
	slog.InfoContext(ctx, strings.TrimSuffix(fmt.Sprintln("key_trimmed", key, limit), "\n"))
	err := uc.cacheProvider.LTrim(ctx, key, limit, -1)
	// if err != nil {
	// 	return _errors.InternalServerException("get user products error: %w", err.Error())
	// }
	return err
}

func (uc *UpdateDataUsecase) Put(ctx context.Context, ownerIDs []uint64, resource string, resourceId uint64) error {
	for _, ownerId := range ownerIDs {
		key := fmt.Sprintf(keyUserUpdate, resource, ownerId)
		// Sử dụng Redis Set để lưu danh sách product IDs
		err := uc.cacheProvider.RPush(ctx, key, resourceId)
		if err != nil {
			continue
		}
		// client.RPush(cloneCtx, key, entity.ID).Err()
	}
	// if productUser.OriginProfileID != nil {
	// 	ids[idx] = *productUser.OriginProfileID

	// }
	return nil
}
