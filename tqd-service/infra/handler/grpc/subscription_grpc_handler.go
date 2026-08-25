package handler_grpc

import (
	"context"
	"log"
	"strings"
	"time"

	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/interface/provider"
	"tqd/internal/usecase"

	"google.golang.org/protobuf/types/known/emptypb"
)

type SubscriptionGrpcHandler struct {
	tqdpb.UnimplementedSubscriptionServiceServer
	subUsecase    usecase.SubscriptionUsecase
	notifUsecase  usecase.NotificationUsecase
	regionUsecase usecase.RegionUsecase
	parcelUsecase usecase.ParcelUsecase
	userProvider  provider.UserProvider
	parcelMapper  *mapper.ParcelMapper
	SyncProvider  *_utils.SyncUtil
}

func NewSubscriptionGrpcHandler(
	subUsecase usecase.SubscriptionUsecase,
	notifUsecase usecase.NotificationUsecase,
	regionUsecase usecase.RegionUsecase,
	parcelUsecase usecase.ParcelUsecase,
	userProvider provider.UserProvider,
	syncProvider *_utils.SyncUtil,
) *SubscriptionGrpcHandler {
	return &SubscriptionGrpcHandler{
		subUsecase:    subUsecase,
		notifUsecase:  notifUsecase,
		regionUsecase: regionUsecase,
		parcelUsecase: parcelUsecase,
		userProvider:  userProvider,
		parcelMapper:  mapper.NewParcelMapper(),
		SyncProvider:  syncProvider,
	}
}

// =====================================================
// PARCEL SUBSCRIPTIONS
// =====================================================

// ListParcelSubscriptions - Danh sách subscription theo dõi parcel
func (h *SubscriptionGrpcHandler) ListParcelSubscriptions(
	ctx context.Context,
	req *tqdpb.ListParcelSubscriptionsRequest,
) (*tqdpb.ListParcelSubscriptionsResponse, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDSubscriptionList, userID)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListParcelSubscriptionsResponse{}, nil
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)
	page := int(pagable.GetPage())
	limit := pagable.GetLimit()

	log.Printf("[ListParcelSubscriptions] UserID: %d, Page: %d, Limit: %d", userID, page, limit)

	subs, total, err := h.subUsecase.ListParcelSubscriptions(ctx, userID, req.Status, page, limit)
	if err != nil {
		log.Printf("[ListParcelSubscriptions] Error: %v", err)
		return nil, _errors.ReturnError(500, err.Error())
	}

	if len(subs) == 0 {
		return &tqdpb.ListParcelSubscriptionsResponse{
			Data:  []*tqdpb.ParcelSubscriptionItem{},
			Total: 0,
			Page:  int32(pagable.GetPage()),
			Size:  int32(pagable.GetSize()),
		}, nil
	}

	// Lấy thông tin parcels
	parcelIDs := make([]uint64, 0, len(subs))
	for _, sub := range subs {
		parcelIDs = append(parcelIDs, sub.TargetID)
	}

	parcels, err := h.parcelUsecase.GetByIDs(ctx, parcelIDs)
	if err != nil {
		parcels = []qh_domain.Parcel{}
	}

	parcelMap := make(map[uint64]qh_domain.Parcel)
	for _, p := range parcels {
		parcelMap[p.ID] = p
	}

	// Build response
	items := make([]*tqdpb.ParcelSubscriptionItem, 0, len(subs))
	for _, sub := range subs {
		scope, _ := dto.SubscriptionScopeJSONToProto(sub.SubscriptionScope)

		item := &tqdpb.ParcelSubscriptionItem{
			SubscriptionId:    sub.ID,
			ParcelId:          sub.TargetID,
			SubscriptionScope: scope,
			Status:            uint32(sub.Status),
			CreatedAt:         sub.CreatedAt.Format(time.RFC3339),
		}

		if parcel, ok := parcelMap[sub.TargetID]; ok {
			parcelProto := h.parcelMapper.DomainToProto(&parcel)
			item.Geometry = parcelProto.Geometry
			lat := parcel.Lat
			lng := parcel.Lng
			item.Lat = &lat
			item.Lng = &lng
		}

		items = append(items, item)
	}

	return &tqdpb.ListParcelSubscriptionsResponse{
		Data:  items,
		Total: total,
		Page:  int32(pagable.GetPage()),
		Size:  int32(pagable.GetSize()),
	}, nil
}

// GetParcelSubscriptionStatus - Lấy trạng thái theo dõi parcel
func (h *SubscriptionGrpcHandler) GetParcelSubscriptionStatus(
	ctx context.Context,
	req *tqdpb.GetParcelSubscriptionStatusRequest,
) (*tqdpb.ParcelSubscriptionStatusResponse, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	if req.ParcelId == 0 {
		return nil, _errors.ReturnError(400, "parcel_id is required")
	}

	log.Printf("[GetParcelSubscriptionStatus] UserID: %d, ParcelID: %d", userID, req.ParcelId)

	isFollowing, sub, err := h.subUsecase.GetParcelStatus(ctx, userID, req.ParcelId)
	if err != nil {
		log.Printf("[GetParcelSubscriptionStatus] Error: %v", err)
		return nil, _errors.ReturnError(500, err.Error())
	}

	resp := &tqdpb.ParcelSubscriptionStatusResponse{IsFollowing: isFollowing}
	if sub != nil {
		subscriptionId := sub.ID
		resp.SubscriptionId = &subscriptionId

		scope, _ := dto.SubscriptionScopeJSONToProto(sub.SubscriptionScope)
		resp.SubscriptionScope = scope

		createdAt := sub.CreatedAt.Format(time.RFC3339)
		resp.CreatedAt = &createdAt
	}
	return resp, nil
}

// CreateParcelSubscription - Tạo subscription theo dõi parcel
func (h *SubscriptionGrpcHandler) CreateParcelSubscription(
	ctx context.Context,
	req *tqdpb.CreateParcelSubscriptionRequest,
) (*tqdpb.SubscriptionResponse, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	if req.ParcelId == 0 {
		return nil, _errors.ReturnError(400, "parcel_id is required")
	}

	log.Printf("[CreateParcelSubscription] UserID: %d, ParcelID: %d", userID, req.ParcelId)

	sub, err := h.subUsecase.CreateParcelSubscription(ctx, userID, req.ParcelId, req.SubscriptionScope, req.GetTriggerConfig())
	if err != nil {
		log.Printf("[CreateParcelSubscription] Error: %v", err)
		return nil, _errors.ReturnError(400, err.Error())
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDSubscriptionList, userID), t.UnixMilli())

	scope, _ := dto.SubscriptionScopeJSONToProto(sub.SubscriptionScope)
	return &tqdpb.SubscriptionResponse{
		Id:                sub.ID,
		TargetType:        sub.TargetType,
		TargetId:          sub.TargetID,
		SubscriptionScope: scope,
		Status:            uint32(sub.Status),
		CreatedAt:         sub.CreatedAt.Format(time.RFC3339),
	}, nil
}

// DeleteParcelSubscription - Xóa subscription theo dõi parcel
func (h *SubscriptionGrpcHandler) DeleteParcelSubscription(
	ctx context.Context,
	req *tqdpb.DeleteParcelSubscriptionRequest,
) (*emptypb.Empty, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	if req.ParcelId == 0 {
		return nil, _errors.ReturnError(400, "parcel_id is required")
	}

	log.Printf("[DeleteParcelSubscription] UserID: %d, ParcelID: %d", userID, req.ParcelId)

	if err := h.subUsecase.DeleteParcelSubscription(ctx, userID, req.ParcelId); err != nil {
		log.Printf("[DeleteParcelSubscription] Error: %v", err)
		return nil, _errors.ReturnError(400, err.Error())
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDSubscriptionList, userID), t.UnixMilli())

	return &emptypb.Empty{}, nil
}

func (h *SubscriptionGrpcHandler) ListRegionSubscriptions(
	ctx context.Context,
	req *tqdpb.ListRegionSubscriptionsRequest,
) (*tqdpb.ListRegionSubscriptionsResponse, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDSubscriptionList, userID)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListRegionSubscriptionsResponse{}, nil
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)
	page := int(pagable.GetPage())
	limit := pagable.GetLimit()

	log.Printf("[ListRegionSubscriptions] UserID: %d, Page: %d, Limit: %d", userID, page, limit)

	subs, regions, total, err := h.subUsecase.ListUserRegionSubscriptions(ctx, userID, page, limit)
	if err != nil {
		log.Printf("[ListRegionSubscriptions] Error: %v", err)
		return nil, _errors.ReturnError(500, err.Error())
	}

	items := make([]*tqdpb.RegionSubscriptionItem, 0, len(subs))
	for i, sub := range subs {
		scope, _ := dto.SubscriptionScopeJSONToProto(sub.SubscriptionScope)
		item := &tqdpb.RegionSubscriptionItem{
			SubscriptionId:    sub.ID,
			SubscriptionScope: scope,
			Status:            uint32(sub.Status),
			CreatedAt:         sub.CreatedAt.Format(time.RFC3339),
		}
		if i < len(regions) {
			item.RegionName = regions[i].Name
			// Geometry is bytes, need to convert
			// if len(regions[i].RawGeometry) > 0 {
			// 	geom := regions[i].RawGeometry.
			// 	item.Geometry = geom
			// }
		}
		items = append(items, item)
	}

	return &tqdpb.ListRegionSubscriptionsResponse{
		Data:  items,
		Total: total,
		Page:  int32(pagable.GetPage()),
		Size:  int32(pagable.GetSize()),
	}, nil
}

// CreateRegionSubscription - Tạo subscription theo dõi region
func (h *SubscriptionGrpcHandler) CreateRegionSubscription(
	ctx context.Context,
	req *tqdpb.CreateRegionSubscriptionRequest,
) (*tqdpb.SubscriptionResponse, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	if req.RegionId == 0 {
		return nil, _errors.ReturnError(400, "region_id is required")
	}

	log.Printf("[CreateRegionSubscription] UserID: %d, RegionID: %d", userID, req.RegionId)

	_, err := h.regionUsecase.GetByID(ctx, req.RegionId)
	if err != nil {
		return nil, _errors.ReturnError(404, "region not found")
	}

	sub, err := h.subUsecase.CreateRegionSubscription(ctx, userID, req.RegionId, req.SubscriptionScope, req.GetTriggerConfig())
	if err != nil {
		log.Printf("[CreateRegionSubscription] Error: %v", err)
		return nil, _errors.ReturnError(400, err.Error())
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDSubscriptionList, userID), t.UnixMilli())

	scope, _ := dto.SubscriptionScopeJSONToProto(sub.SubscriptionScope)
	return &tqdpb.SubscriptionResponse{
		Id:                sub.ID,
		TargetType:        sub.TargetType,
		TargetId:          sub.TargetID,
		SubscriptionScope: scope,
		Status:            uint32(sub.Status),
		CreatedAt:         sub.CreatedAt.Format(time.RFC3339),
	}, nil
}

// DeleteRegionSubscription - Xóa subscription theo dõi region
func (h *SubscriptionGrpcHandler) DeleteRegionSubscription(
	ctx context.Context,
	req *tqdpb.DeleteRegionSubscriptionRequest,
) (*emptypb.Empty, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	if req.SubscriptionId == 0 {
		return nil, _errors.ReturnError(400, "subscription_id is required")
	}

	log.Printf("[DeleteRegionSubscription] UserID: %d, SubscriptionID: %d", userID, req.SubscriptionId)

	if err := h.subUsecase.DeleteRegionSubscription(ctx, userID, req.SubscriptionId); err != nil {
		log.Printf("[DeleteRegionSubscription] Error: %v", err)
		return nil, _errors.ReturnError(400, err.Error())
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDSubscriptionList, userID), t.UnixMilli())

	return &emptypb.Empty{}, nil
}

// =====================================================
// NOTIFICATIONS
// =====================================================

// ListNotifications - Danh sách thông báo
func (h *SubscriptionGrpcHandler) ListNotifications(
	ctx context.Context,
	req *tqdpb.ListNotificationsRequest,
) (*tqdpb.ListNotificationsResponse, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDNotificationList, userID)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListNotificationsResponse{}, nil
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)
	page := int(pagable.GetPage())
	limit := pagable.GetLimit()

	log.Printf("[ListNotifications] UserID: %d, Page: %d, Limit: %d", userID, page, limit)

	notifs, total, unread, err := h.notifUsecase.List(ctx, userID, req.Type, req.Read, page, limit)
	if err != nil {
		log.Printf("[ListNotifications] Error: %v", err)
		return nil, _errors.ReturnError(500, err.Error())
	}

	items := make([]*tqdpb.NotificationItem, 0, len(notifs))
	for _, n := range notifs {
		item := &tqdpb.NotificationItem{
			Id:               n.ID,
			Title:            n.Title,
			Message:          n.Message,
			NotificationType: uint32(n.NotificationType),
			Severity:         uint32(n.Severity),
			IsRead:           n.IsRead,
			CreatedAt:        n.CreatedAt.Format(time.RFC3339),
		}
		if n.ActionLink != nil && *n.ActionLink != "" {
			actionLink := *n.ActionLink
			item.ActionLink = &actionLink
		}
		items = append(items, item)
	}

	return &tqdpb.ListNotificationsResponse{
		Data:        items,
		Total:       total,
		UnreadCount: unread,
		Page:        int32(pagable.GetPage()),
		Size:        int32(pagable.GetSize()),
	}, nil
}

// MarkNotificationRead - Đánh dấu đã đọc 1 thông báo
func (h *SubscriptionGrpcHandler) MarkNotificationRead(
	ctx context.Context,
	req *tqdpb.MarkNotificationReadRequest,
) (*emptypb.Empty, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	if req.NotificationId == 0 {
		return nil, _errors.ReturnError(400, "notification_id is required")
	}

	log.Printf("[MarkNotificationRead] UserID: %d, NotificationID: %d", userID, req.NotificationId)

	if err := h.notifUsecase.MarkRead(ctx, userID, req.NotificationId); err != nil {
		log.Printf("[MarkNotificationRead] Error: %v", err)
		return nil, _errors.ReturnError(500, err.Error())
	}
	return &emptypb.Empty{}, nil
}

// MarkAllNotificationsRead - Đánh dấu đã đọc tất cả thông báo
func (h *SubscriptionGrpcHandler) MarkAllNotificationsRead(
	ctx context.Context,
	_ *emptypb.Empty,
) (*emptypb.Empty, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	log.Printf("[MarkAllNotificationsRead] UserID: %d", userID)

	if err := h.notifUsecase.MarkAllRead(ctx, userID); err != nil {
		log.Printf("[MarkAllNotificationsRead] Error: %v", err)
		return nil, _errors.ReturnError(500, err.Error())
	}
	return &emptypb.Empty{}, nil
}

// =====================================================
// ADMIN APIs
// =====================================================

// AdminListSubscriptions - Admin: Danh sách subscriptions
func (h *SubscriptionGrpcHandler) AdminListSubscriptions(
	ctx context.Context,
	req *tqdpb.AdminListSubscriptionsRequest,
) (*tqdpb.ListSubscriptionsAdminResponse, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	log.Printf("[AdminListSubscriptions] AdminUserID: %d", userID)

	// Sử dụng Pagable
	pagable := _dto.NewPagableFromGrpc(req.Page, req.Size, nil)
	page := int(pagable.GetPage())
	limit := pagable.GetLimit()

	log.Printf("[AdminListSubscriptions] Page: %d, Limit: %d", page, limit)

	subs, total, err := h.subUsecase.AdminList(ctx, req.UserId, req.TargetType, (*uint32)(req.Status), page, limit)
	if err != nil {
		log.Printf("[AdminListSubscriptions] Error: %v", err)
		return nil, _errors.ReturnError(500, err.Error())
	}

	items := make([]*tqdpb.SubscriptionAdminItem, 0, len(subs))
	for _, sub := range subs {
		scope, _ := dto.SubscriptionScopeJSONToProto(sub.SubscriptionScope)
		items = append(items, &tqdpb.SubscriptionAdminItem{
			Id:                sub.ID,
			UserId:            sub.UserID,
			TargetType:        sub.TargetType,
			TargetId:          sub.TargetID,
			SubscriptionScope: scope,
			Status:            uint32(sub.Status),
			CreatedAt:         sub.CreatedAt.Format(time.RFC3339),
			UpdatedAt:         sub.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &tqdpb.ListSubscriptionsAdminResponse{
		Data:  items,
		Total: total,
		Page:  int32(pagable.GetPage()),
		Size:  int32(pagable.GetSize()),
	}, nil
}

// AdminListParcelSubscriptions - Admin: danh sách thửa đất đang được theo dõi (subscription)
func (h *SubscriptionGrpcHandler) AdminListParcelSubscriptions(
	ctx context.Context,
	req *tqdpb.AdminListParcelSubscriptionsRequest,
) (*tqdpb.AdminListParcelSubscriptionsResponse, error) {
	adminUserID := _utils.GetProfileIdWithContext(ctx)
	log.Printf("[AdminListParcelSubscriptions] AdminUserID: %d", adminUserID)

	pagable := _dto.NewPagableFromGrpc(req.Page, req.Size, nil)
	page := int(pagable.GetPage())
	limit := pagable.GetLimit()

	subs, total, err := h.subUsecase.AdminListParcelSubscriptions(ctx, req.UserId, req.Status, page, limit)
	if err != nil {
		log.Printf("[AdminListParcelSubscriptions] Error: %v", err)
		return nil, _errors.ReturnError(500, err.Error())
	}

	if len(subs) == 0 {
		return &tqdpb.AdminListParcelSubscriptionsResponse{
			Data:  []*tqdpb.AdminParcelSubscriptionItem{},
			Total: 0,
			Page:  int32(pagable.GetPage()),
			Size:  int32(pagable.GetSize()),
		}, nil
	}

	parcelIDs := make([]uint64, 0, len(subs))
	for _, sub := range subs {
		parcelIDs = append(parcelIDs, sub.TargetID)
	}

	parcels, err := h.parcelUsecase.GetByIDs(ctx, parcelIDs)
	if err != nil {
		parcels = []qh_domain.Parcel{}
	}
	parcelMap := make(map[uint64]qh_domain.Parcel, len(parcels))
	for _, p := range parcels {
		parcelMap[p.ID] = p
	}

	parcelInfoMap := make(map[uint64]*dto.ParcelInfoResponse, len(parcelIDs))
	for _, parcelID := range parcelIDs {
		info, infoErr := h.parcelUsecase.GetParcelInfo(ctx, parcelID)
		if infoErr == nil && info != nil && info.ParcelID > 0 {
			parcelInfoMap[parcelID] = info
		}
	}

	activeStatus := uint32(10)
	followerCounts, err := h.subUsecase.CountParcelFollowers(ctx, parcelIDs, &activeStatus)
	if err != nil {
		log.Printf("[AdminListParcelSubscriptions] CountFollowers error: %v", err)
		followerCounts = map[uint64]int64{}
	}

	userIDs := make([]uint64, 0, len(subs))
	seenUserIDs := make(map[uint64]struct{}, len(subs))
	for _, sub := range subs {
		if sub.UserID == 0 {
			continue
		}
		if _, ok := seenUserIDs[sub.UserID]; ok {
			continue
		}
		seenUserIDs[sub.UserID] = struct{}{}
		userIDs = append(userIDs, sub.UserID)
	}

	profileMap := map[uint64]*dto.UserProfileDTO{}
	if h.userProvider != nil && len(userIDs) > 0 {
		if profiles, profileErr := h.userProvider.GetProfilesByIDs(ctx, userIDs); profileErr == nil {
			profileMap = profiles
		} else {
			log.Printf("[AdminListParcelSubscriptions] GetProfilesByIDs error: %v", profileErr)
		}
	}

	items := make([]*tqdpb.AdminParcelSubscriptionItem, 0, len(subs))
	for _, sub := range subs {
		scope, _ := dto.SubscriptionScopeJSONToProto(sub.SubscriptionScope)
		item := &tqdpb.AdminParcelSubscriptionItem{
			SubscriptionId:    sub.ID,
			UserId:            sub.UserID,
			ParcelId:          sub.TargetID,
			SubscriptionScope: scope,
			Status:            uint32(sub.Status),
			CreatedAt:         sub.CreatedAt.Format(time.RFC3339),
			FollowerCount:     followerCounts[sub.TargetID],
		}

		if parcel, ok := parcelMap[sub.TargetID]; ok {
			parcelProto := h.parcelMapper.DomainToProto(&parcel)
			item.Geometry = parcelProto.Geometry
			lat := parcel.Lat
			lng := parcel.Lng
			item.Lat = &lat
			item.Lng = &lng
			if parcel.SeoID != nil && *parcel.SeoID > 0 {
				seoID := *parcel.SeoID
				item.SeoId = &seoID
			}
		}
		if info, ok := parcelInfoMap[sub.TargetID]; ok {
			if info.AddressText != "" {
				addr := info.AddressText
				item.Address = &addr
			}
			if info.MapNumber != "" || info.LandNumber != "" {
				name := info.MapNumber + "/" + info.LandNumber
				item.ParcelName = &name
			}
		}

		if profile := profileMap[sub.UserID]; profile != nil {
			if fn := strings.TrimSpace(profile.FullName); fn != "" {
				item.UserName = &fn
			}
			if av := strings.TrimSpace(profile.Avatar); av != "" {
				item.UserAvatar = &av
			}
		}

		items = append(items, item)
	}

	return &tqdpb.AdminListParcelSubscriptionsResponse{
		Data:  items,
		Total: total,
		Page:  int32(pagable.GetPage()),
		Size:  int32(pagable.GetSize()),
	}, nil
}

// AdminDeleteSubscription - Admin: Xóa subscription
func (h *SubscriptionGrpcHandler) AdminDeleteSubscription(
	ctx context.Context,
	req *tqdpb.AdminDeleteSubscriptionRequest,
) (*emptypb.Empty, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	log.Printf("[AdminDeleteSubscription] AdminUserID: %d, SubscriptionID: %d", userID, req.SubscriptionId)

	if req.SubscriptionId == 0 {
		return nil, _errors.ReturnError(400, "subscription_id is required")
	}

	if err := h.subUsecase.AdminDelete(ctx, req.SubscriptionId); err != nil {
		log.Printf("[AdminDeleteSubscription] Error: %v", err)
		return nil, _errors.ReturnError(500, err.Error())
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDSubscriptionList, 0), t.UnixMilli())

	return &emptypb.Empty{}, nil
}
