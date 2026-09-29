package grpcadapter

import (
	"common/request"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"tqd/internal/enums"
	"tqd/internal/usecase/generatedreport/application"
	"tqd/internal/usecase/quota"

	tqdpb "pb/types/tqd"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**
 * Adapter chuyển proto request/response sang report service.
 *
 * File này chỉ làm việc ở gRPC boundary, không chứa quota hay DB logic.
 */
type Adapter struct {
	service *application.Service
}

func NewAdapter(service *application.Service) *Adapter {
	return &Adapter{service: service}
}

/**
 * CreateGeneratedReport nhận proto request và gọi đúng report flow mới.
 */
func (a *Adapter) CreateGeneratedReport(
	ctx context.Context,
	req *tqdpb.CreateGeneratedReportRequest,
) (*tqdpb.CreateGeneratedReportResponse, error) {
	if a == nil || a.service == nil {
		return nil, status.Error(codes.Internal, "report service is not configured")
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	ctx, err := bindRequestCommandKeyCompatibility(ctx, requestCommandKeyCompatibility(req))
	if err != nil {
		return nil, err
	}

	result, err := a.service.CreateGeneratedReport(ctx, fromProto(req))
	if err != nil {
		return nil, toStatus(err)
	}

	message := "report generation started"
	if result.Outcome == application.OutcomeExisting {
		message = "report request already accepted"
	}

	return &tqdpb.CreateGeneratedReportResponse{
		Success: true,
		Message: message,
		Item:    toProto(result.Report),
	}, nil
}

// requestCommandKeyCompatibility reads protobuf field #11 by number rather
// than generated accessor name. CP03 renames the proto source field from
// clientRequestId to commandKey while generated code is intentionally deferred
// to the generated-contract checkpoint. Field number 11 remains wire-stable.
func requestCommandKeyCompatibility(req *tqdpb.CreateGeneratedReportRequest) string {
	if req == nil {
		return ""
	}
	message := req.ProtoReflect()
	field := message.Descriptor().Fields().ByNumber(11)
	if field == nil {
		return ""
	}
	return strings.TrimSpace(message.Get(field).String())
}

// bindRequestCommandKeyCompatibility preserves old body-only clients during
// the cutover. Canonical identity is Idempotency-Key; if field #11 is also
// supplied, both transports must agree.
func bindRequestCommandKeyCompatibility(ctx context.Context, wireValue string) (context.Context, error) {
	wireValue = strings.TrimSpace(wireValue)
	if wireValue == "" {
		return ctx, nil
	}
	if existing, ok := request.IdempotencyKeyFromContext(ctx); ok {
		if existing != wireValue {
			return ctx, status.Error(codes.InvalidArgument, "idempotency key conflicts with report command key")
		}
		return ctx, nil
	}

	bound, err := request.BindIdempotencyKey(ctx, wireValue)
	if err != nil {
		return ctx, status.Error(codes.InvalidArgument, "report command key is invalid")
	}
	return bound, nil
}

func fromProto(req *tqdpb.CreateGeneratedReportRequest) application.Input {
	input := application.Input{
		ReportType:   req.GetReportType(),
		EntityType:   req.GetEntityType(),
		EntityID:     req.GetEntityId(),
		ParcelID:     req.GetParcelId(),
		RegionID:     req.GetRegionId(),
		Title:        req.GetTitle(),
		Subtitle:     req.GetSubtitle(),
		MetadataJSON: req.GetMetadataJson(),
	}

	if value := req.GetLocation(); value != nil {
		input.Location = application.Location{
			Address:      value.GetAddress(),
			Province:     value.GetProvince(),
			ProvinceCode: value.GetProvinceCode(),
			WardCode:     value.GetWardCode(),
		}
	}

	if value := req.GetSpatial(); value != nil {
		if bounds := value.GetBounds(); bounds != nil {
			input.Spatial.Bounds = application.Bounds{
				MinLon: bounds.GetMinLon(),
				MinLat: bounds.GetMinLat(),
				MaxLon: bounds.GetMaxLon(),
				MaxLat: bounds.GetMaxLat(),
			}
		}
		if center := value.GetCentroid(); center != nil {
			input.Spatial.Centroid = application.Point{
				Lat: center.GetLat(),
				Lon: center.GetLon(),
			}
		}
	}

	if value := req.GetComparison(); value != nil {
		input.Comparison = application.Comparison{
			FromPlanName: value.GetFromPlanName(),
			ToPlanName:   value.GetToPlanName(),
			FromYear:     value.GetFromYear(),
			ToYear:       value.GetToYear(),
		}
	}

	return input
}

func toProto(item application.Report) *tqdpb.GeneratedReportPreview {
	comparison := application.Comparison{}
	if len(item.ComparisonJSON) > 0 {
		_ = json.Unmarshal(item.ComparisonJSON, &comparison)
	}

	entityType := enums.WorkspaceEntityTypeReport.Uint32()
	entityID := item.ID
	var parcelID uint64
	var regionID uint64

	if item.ParcelID != nil {
		parcelID = *item.ParcelID
		entityType = enums.WorkspaceEntityTypeParcel.Uint32()
		entityID = parcelID
	}
	if item.RegionID != nil {
		regionID = *item.RegionID
		entityType = enums.WorkspaceEntityTypeRegion.Uint32()
		entityID = regionID
	}

	canDownload := item.Status.CanDownload()
	return &tqdpb.GeneratedReportPreview{
		Id:         fmt.Sprintf("report_%d", item.ID),
		Type:       enums.WorkspaceEntityTypeReport.Uint32(),
		ReportId:   item.ID,
		UserId:     item.UserID,
		ReportType: uint32(item.ReportType),
		Status:     item.Status.Uint32(),
		Title:      item.Title,
		Subtitle:   item.Subtitle,
		Location: &tqdpb.WorkspaceLocation{
			Address:      item.Location.Address,
			Province:     item.Location.Province,
			ProvinceCode: item.Location.ProvinceCode,
			WardCode:     item.Location.WardCode,
		},
		Spatial: &tqdpb.ReportSpatialPreview{
			Bounds: &tqdpb.SpatialBounds{
				MinLon: item.Spatial.Bounds.MinLon,
				MinLat: item.Spatial.Bounds.MinLat,
				MaxLon: item.Spatial.Bounds.MaxLon,
				MaxLat: item.Spatial.Bounds.MaxLat,
			},
			Centroid: &tqdpb.SpatialPoint{
				Lat: item.Spatial.Centroid.Lat,
				Lon: item.Spatial.Centroid.Lon,
			},
		},
		Assets: &tqdpb.ReportAssets{
			ThumbnailUrl: item.ThumbnailURL,
			ImageUrl:     item.ImageURL,
			PdfUrl:       item.PDFURL,
			ShareUrl:     item.ShareURL,
		},
		ReportMeta: &tqdpb.ReportMeta{
			CreatedAt: formatTime(item.CreatedAt),
			UpdatedAt: formatTime(item.UpdatedAt),
			ExpiresAt: formatTimePointer(item.ExpiresAt),
			FileSize:  item.FileSize,
			Format:    item.Format,
		},
		Comparison: &tqdpb.ReportComparison{
			FromPlanName: comparison.FromPlanName,
			ToPlanName:   comparison.ToPlanName,
			FromYear:     comparison.FromYear,
			ToYear:       comparison.ToYear,
		},
		Actions: &tqdpb.ReportActions{
			Focus:         hasBounds(item.Spatial.Bounds),
			DownloadImage: canDownload && item.ImageURL != "",
			DownloadPdf:   canDownload && item.PDFURL != "",
			Share:         item.Status.CanShare(),
			Regenerate:    item.Status.CanRegenerate(),
			Remove:        item.Status.CanRemove(),
		},
		EntityType: entityType,
		EntityId:   entityID,
		ParcelId:   parcelID,
		RegionId:   regionID,
	}
}

func toStatus(err error) error {
	var exhausted *quota.ExhaustedError
	var denied *quota.AccessDeniedError
	switch {
	case errors.Is(err, application.ErrActorMissing),
		errors.Is(err, application.ErrProfileMissing):
		return status.Error(codes.Unauthenticated, "profile identity is required")
	case errors.Is(err, application.ErrCommandConflict):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, application.ErrInvalidInput),
		errors.Is(err, application.ErrCommandKeyMissing),
		errors.Is(err, application.ErrOperationIDMissing):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, application.ErrParcelNotFound),
		errors.Is(err, application.ErrRegionNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.As(err, &denied):
		return accessDeniedStatus(denied)
	case errors.Is(err, quota.ErrAccessDenied):
		return status.Error(codes.PermissionDenied, "commercial entitlement required")
	case errors.As(err, &exhausted):
		return quotaExhaustedStatus(exhausted)
	case errors.Is(err, quota.ErrQuotaExceeded):
		return status.Error(codes.ResourceExhausted, "quota exhausted")
	case errors.Is(err, application.ErrAcceptanceUnavailable):
		return status.Error(
			codes.Unavailable,
			application.ErrAcceptanceUnavailable.Error(),
		)
	case errors.Is(err, application.ErrAccessUnavailable):
		return status.Error(
			codes.Unavailable,
			application.ErrAccessUnavailable.Error(),
		)
	case errors.Is(err, application.ErrProcessingUnavailable):
		return status.Error(codes.Unavailable, err.Error())
	default:
		return status.Error(codes.Internal, "cannot create generated report")
	}
}

func accessDeniedStatus(denied *quota.AccessDeniedError) error {
	if denied == nil {
		return status.Error(codes.PermissionDenied, "commercial entitlement required")
	}
	result, detailsErr := status.New(codes.PermissionDenied, "commercial entitlement required").WithDetails(
		&errdetails.ErrorInfo{
			Reason: "ENTITLEMENT_DENIED",
			Domain: "tqd.qhpro",
			Metadata: map[string]string{
				"subject_type": string(denied.Subject.Type),
				"subject_id":   denied.Subject.ID,
				"operation":    string(denied.Operation),
			},
		},
	)
	if detailsErr != nil {
		return status.Error(codes.PermissionDenied, "commercial entitlement required")
	}
	return result.Err()
}

func quotaExhaustedStatus(exhausted *quota.ExhaustedError) error {
	if exhausted == nil {
		return status.Error(codes.ResourceExhausted, "quota exhausted")
	}
	metadata := map[string]string{
		"subject_type": string(exhausted.Subject.Type),
		"subject_id":   exhausted.Subject.ID,
		"operation":    string(exhausted.Operation),
		"meter":        string(exhausted.MeterCode),
		"limit":        strconv.FormatInt(exhausted.Limit, 10),
		"used":         strconv.FormatInt(exhausted.Used, 10),
		"reserved":     strconv.FormatInt(exhausted.Reserved, 10),
		"remaining":    strconv.FormatInt(exhausted.Remaining, 10),
	}
	if !exhausted.PeriodStart.IsZero() {
		metadata["period_start"] = exhausted.PeriodStart.UTC().Format(time.RFC3339)
	}
	if !exhausted.PeriodEnd.IsZero() {
		metadata["period_end"] = exhausted.PeriodEnd.UTC().Format(time.RFC3339)
	}
	result, detailsErr := status.New(codes.ResourceExhausted, "quota exhausted").WithDetails(
		&errdetails.ErrorInfo{
			Reason:   "QUOTA_EXHAUSTED",
			Domain:   "tqd.qhpro",
			Metadata: metadata,
		},
	)
	if detailsErr != nil {
		return status.Error(codes.ResourceExhausted, "quota exhausted")
	}
	return result.Err()
}

func hasBounds(value application.Bounds) bool {
	return value.MinLon != 0 || value.MinLat != 0 || value.MaxLon != 0 || value.MaxLat != 0
}

func formatTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}

func formatTimePointer(value *time.Time) string {
	if value == nil {
		return ""
	}
	return formatTime(*value)
}
