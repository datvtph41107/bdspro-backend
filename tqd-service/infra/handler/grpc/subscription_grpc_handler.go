package handler_grpc

import (
	"context"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"google.golang.org/protobuf/types/known/emptypb"
)

type SubscriptionGrpcHandler struct {
	tqdpb.UnimplementedSubscriptionServiceServer
	subUsecase    usecase.SubscriptionUsecase
	notifUsecase  usecase.NotificationUsecase
	regionUsecase usecase.RegionUsecase
	parcelUsecase usecase.ParcelUsecase
	parcelMapper  *mapper.ParcelMapper
}

func NewSubscriptionGrpcHandler(
	subUsecase usecase.SubscriptionUsecase,
	notifUsecase usecase.NotificationUsecase,
	regionUsecase usecase.RegionUsecase,
	parcelUsecase usecase.ParcelUsecase,
) *SubscriptionGrpcHandler {
	return &SubscriptionGrpcHandler{
		subUsecase:    subUsecase,
		notifUsecase:  notifUsecase,
		regionUsecase: regionUsecase,
		parcelUsecase: parcelUsecase,
		parcelMapper:  mapper.NewParcelMapper(),
	}
}

func (h *SubscriptionGrpcHandler) ListParcelSubscriptions(
	ctx context.Context,
	req *tqdpb.ListParcelSubscriptionsRequest,
) (*tqdpb.ListParcelSubscriptionsResponse, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}
	page := int32(1)
	limit := int32(20)
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}

	subs, total, err := h.subUsecase.ListParcelSubscriptions(ctx, userID, req.Status, int(page), int(limit))
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	if len(subs) == 0 {
		return &tqdpb.ListParcelSubscriptionsResponse{
			Data:  []*tqdpb.ParcelSubscriptionItem{},
			Total: 0,
			Page:  page,
			Size:  limit,
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

	parcelMap := make(map[uint64]qh_domain.Parcel)
	for _, p := range parcels {
		parcelMap[p.ID] = p
	}

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
			item.Lat = &parcel.Lat
			item.Lng = &parcel.Lng
		}

		items = append(items, item)
	}

	return &tqdpb.ListParcelSubscriptionsResponse{
		Data:  items,
		Total: total,
		Page:  page,
		Size:  limit,
	}, nil
}

func (h *SubscriptionGrpcHandler) GetParcelSubscriptionStatus(
	ctx context.Context,
	req *tqdpb.GetParcelSubscriptionStatusRequest,
) (*tqdpb.ParcelSubscriptionStatusResponse, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}
	isFollowing, sub, err := h.subUsecase.GetParcelStatus(ctx, userID, req.ParcelId)
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	resp := &tqdpb.ParcelSubscriptionStatusResponse{IsFollowing: isFollowing}
	if sub != nil {
		subscriptionId := sub.ID
		resp.SubscriptionId = &subscriptionId

		// Convert từ datatypes.JSON sang proto
		scope, _ := dto.SubscriptionScopeJSONToProto(sub.SubscriptionScope)
		resp.SubscriptionScope = scope

		createdAt := sub.CreatedAt.Format(time.RFC3339)
		resp.CreatedAt = &createdAt
	}
	return resp, nil
}

func (h *SubscriptionGrpcHandler) CreateParcelSubscription(
	ctx context.Context,
	req *tqdpb.CreateParcelSubscriptionRequest,
) (*tqdpb.SubscriptionResponse, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	sub, err := h.subUsecase.CreateParcelSubscription(ctx, userID, req.ParcelId, req.SubscriptionScope, req.GetTriggerConfig())
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}

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

func (h *SubscriptionGrpcHandler) DeleteParcelSubscription(
	ctx context.Context,
	req *tqdpb.DeleteParcelSubscriptionRequest,
) (*emptypb.Empty, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	if err := h.subUsecase.DeleteParcelSubscription(ctx, userID, req.ParcelId); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	return &emptypb.Empty{}, nil
}

// ==================== Region ====================

func (h *SubscriptionGrpcHandler) ListRegionSubscriptions(
	ctx context.Context,
	req *tqdpb.ListRegionSubscriptionsRequest,
) (*tqdpb.ListRegionSubscriptionsResponse, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	page := int32(1)
	limit := int32(20)
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}

	subs, regions, total, err := h.subUsecase.ListUserRegionSubscriptions(ctx, userID, int(page), int(limit))
	if err != nil {
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
			item.Geometry = regions[i].Geometry
		}
		items = append(items, item)
	}

	return &tqdpb.ListRegionSubscriptionsResponse{
		Data:  items,
		Total: total,
		Page:  page,
		Size:  limit,
	}, nil
}

func (h *SubscriptionGrpcHandler) CreateRegionSubscription(
	ctx context.Context,
	req *tqdpb.CreateRegionSubscriptionRequest,
) (*tqdpb.SubscriptionResponse, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	_, err := h.regionUsecase.GetByID(ctx, req.RegionId)
	if err != nil {
		return nil, _errors.ReturnError(404, "region not found")
	}

	sub, err := h.subUsecase.CreateRegionSubscription(ctx, userID, req.RegionId, req.SubscriptionScope, req.GetTriggerConfig())
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}

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

func (h *SubscriptionGrpcHandler) DeleteRegionSubscription(
	ctx context.Context,
	req *tqdpb.DeleteRegionSubscriptionRequest,
) (*emptypb.Empty, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	if err := h.subUsecase.DeleteRegionSubscription(ctx, userID, req.SubscriptionId); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	return &emptypb.Empty{}, nil
}

// ==================== Notifications ====================

func (h *SubscriptionGrpcHandler) ListNotifications(
	ctx context.Context,
	req *tqdpb.ListNotificationsRequest,
) (*tqdpb.ListNotificationsResponse, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	page := int32(1)
	limit := int32(20)
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}

	notifs, total, unread, err := h.notifUsecase.List(ctx, userID, req.Type, req.Read, int(page), int(limit))
	if err != nil {
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
		if n.ActionLink != nil {
			item.ActionLink = n.ActionLink
		}
		items = append(items, item)
	}

	return &tqdpb.ListNotificationsResponse{
		Data:        items,
		Total:       total,
		UnreadCount: unread,
		Page:        page,
		Size:        limit,
	}, nil
}

func (h *SubscriptionGrpcHandler) MarkNotificationRead(
	ctx context.Context,
	req *tqdpb.MarkNotificationReadRequest,
) (*emptypb.Empty, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	if err := h.notifUsecase.MarkRead(ctx, userID, req.NotificationId); err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h *SubscriptionGrpcHandler) MarkAllNotificationsRead(
	ctx context.Context,
	_ *emptypb.Empty,
) (*emptypb.Empty, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	if err := h.notifUsecase.MarkAllRead(ctx, userID); err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}
	return &emptypb.Empty{}, nil
}

// ==================== Admin ====================

func (h *SubscriptionGrpcHandler) AdminListSubscriptions(
	ctx context.Context,
	req *tqdpb.AdminListSubscriptionsRequest,
) (*tqdpb.ListSubscriptionsAdminResponse, error) {
	page := int32(1)
	limit := int32(20)
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}

	subs, total, err := h.subUsecase.AdminList(ctx, req.UserId, req.TargetType, (*uint32)(req.Status), int(page), int(limit))
	if err != nil {
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
		Page:  page,
		Size:  limit,
	}, nil
}

func (h *SubscriptionGrpcHandler) AdminDeleteSubscription(
	ctx context.Context,
	req *tqdpb.AdminDeleteSubscriptionRequest,
) (*emptypb.Empty, error) {
	if err := h.subUsecase.AdminDelete(ctx, req.SubscriptionId); err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}
	return &emptypb.Empty{}, nil
}
