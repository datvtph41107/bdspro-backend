package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"tqd/internal/enums"
)

/**
 * buildReport kiểm tra input và dựng report trước khi bắt đầu giữ quota.
 *
 * Bước này chỉ đọc parcel/region để bổ sung dữ liệu hiển thị, chưa thay đổi state.
 */
func (s *Service) buildReport(
	ctx context.Context,
	userID uint64,
	commandKey string,
	input Input,
) (Report, error) {
	reportType := enums.GeneratedReportType(input.ReportType)
	if !reportType.IsValid() {
		return Report{}, fmt.Errorf("%w: invalid report type", ErrInvalidInput)
	}

	entityType := enums.WorkspaceEntityTypeFromUint32(input.EntityType)
	if !entityType.IsValid() {
		switch {
		case input.ParcelID > 0:
			entityType = enums.WorkspaceEntityTypeParcel
		case input.RegionID > 0:
			entityType = enums.WorkspaceEntityTypeRegion
		default:
			return Report{}, fmt.Errorf("%w: entity type is required", ErrInvalidInput)
		}
	}
	if !reportType.SupportsEntityType(entityType) {
		return Report{}, fmt.Errorf("%w: report type does not support entity type", ErrInvalidInput)
	}

	report := Report{
		UserID:       userID,
		ReportType:   reportType,
		Status:       enums.GeneratedReportStatusProcessing,
		Title:        strings.TrimSpace(input.Title),
		Subtitle:     strings.TrimSpace(input.Subtitle),
		Location:     input.Location,
		Spatial:      input.Spatial,
		Format:       reportType.DefaultFormat(),
		CommandKey:   commandKey,
		MetadataJSON: []byte("{}"),
	}

	comparison, err := json.Marshal(input.Comparison)
	if err != nil {
		return Report{}, fmt.Errorf("%w: comparison: %v", ErrInvalidInput, err)
	}
	report.ComparisonJSON = comparison

	if strings.TrimSpace(input.MetadataJSON) != "" {
		metadata := []byte(input.MetadataJSON)
		if !json.Valid(metadata) {
			return Report{}, fmt.Errorf("%w: metadata is not valid JSON", ErrInvalidInput)
		}
		report.MetadataJSON = metadata
	}

	switch entityType {
	case enums.WorkspaceEntityTypeParcel:
		parcelID := input.ParcelID
		if parcelID == 0 {
			parcelID = input.EntityID
		}
		if parcelID == 0 {
			return Report{}, fmt.Errorf("%w: parcel ID is required", ErrInvalidInput)
		}

		parcel, found, err := s.source.FindParcel(ctx, parcelID)
		if err != nil {
			return Report{}, err
		}
		if !found {
			return Report{}, ErrParcelNotFound
		}

		report.ParcelID = uint64Ptr(parcelID)
		report.Location = parcel.Location
		report.Spatial = parcel.Spatial
		if report.Title == "" {
			report.Title = fmt.Sprintf("Báo cáo quy hoạch thửa %s", parcel.LandNumber)
		}
		if report.Subtitle == "" {
			report.Subtitle = buildParcelSubtitle(parcel.MapNumber, parcel.AreaSqm)
		}

	case enums.WorkspaceEntityTypeRegion:
		regionID := input.RegionID
		if regionID == 0 {
			regionID = input.EntityID
		}
		if regionID == 0 {
			return Report{}, fmt.Errorf("%w: region ID is required", ErrInvalidInput)
		}

		region, found, err := s.source.FindRegion(ctx, regionID)
		if err != nil {
			return Report{}, err
		}
		if !found {
			return Report{}, ErrRegionNotFound
		}

		report.RegionID = uint64Ptr(regionID)
		report.Location = region.Location
		report.Spatial = region.Spatial
		if report.Title == "" {
			report.Title = "Báo cáo vùng quy hoạch: " + buildRegionTitle(region)
		}
		if report.Subtitle == "" {
			report.Subtitle = buildRegionSubtitle(region)
		}

	case enums.WorkspaceEntityTypeLocation:
		// Snapshot theo location dùng trực tiếp Location/Spatial từ request.
	}

	return report, nil
}

func buildParcelSubtitle(mapNumber string, areaSqm float64) string {
	if mapNumber == "" || mapNumber == "0" {
		if areaSqm <= 0 {
			return ""
		}
		return fmt.Sprintf("Diện tích %.2f m²", areaSqm)
	}
	if areaSqm <= 0 {
		return "Số tờ " + mapNumber
	}
	return fmt.Sprintf("Số tờ %s · Diện tích %.2f m²", mapNumber, areaSqm)
}

func buildRegionTitle(region Region) string {
	if region.DisplayName != "" {
		return region.DisplayName
	}
	if region.LandUseName != "" {
		return region.LandUseName
	}
	if region.Name != "" {
		return region.Name
	}
	return "Vùng quy hoạch"
}

func buildRegionSubtitle(region Region) string {
	parts := make([]string, 0, 3)
	if region.LayerDisplayName != "" {
		parts = append(parts, region.LayerDisplayName)
	} else if region.LayerName != "" {
		parts = append(parts, region.LayerName)
	}
	if region.LandUseName != "" {
		parts = append(parts, region.LandUseName)
	}
	if region.AreaSqm > 0 {
		parts = append(parts, fmt.Sprintf("Diện tích %.2f m²", region.AreaSqm))
	}
	return strings.Join(parts, " · ")
}

func uint64Ptr(value uint64) *uint64 {
	return &value
}
