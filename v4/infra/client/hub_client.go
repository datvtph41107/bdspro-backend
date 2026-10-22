package client

import (
	"context"
	"fmt"
	"strconv"

	"mdn4/modules/hub/repository"
	hubpb "pb/types/hub"
)

// HubClient serves in-process hub lookups from the v4 DSN (no hub-service TCP).
// @bind: mdn4/internal/interface/providers.HubProvider
type HubClient struct {
	location   repository.LocationRepository
	locationV2 repository.LocationV2Repository
	userGuide  repository.UserGuideRepository
}

func NewHubClient(
	location repository.LocationRepository,
	locationV2 repository.LocationV2Repository,
	userGuide repository.UserGuideRepository,
) *HubClient {
	return &HubClient{
		location:   location,
		locationV2: locationV2,
		userGuide:  userGuide,
	}
}

// GetLocationsByIds lấy thông tin province/district/ward (v1) theo IDs.
func (c *HubClient) GetLocationsByIds(ctx context.Context, ids []uint64) (map[uint64]*hubpb.Location, error) {
	out := make(map[uint64]*hubpb.Location)
	if c == nil || len(ids) == 0 {
		return out, nil
	}
	if c.location != nil {
		if provinces, err := c.location.GetProvincesByIDs(ctx, ids); err == nil {
			for i := range provinces {
				p := &provinces[i]
				out[p.ID] = &hubpb.Location{
					Id: p.ID, Name: p.Name, Type: p.Type, TypeText: p.TypeText, Level: 1,
				}
			}
		}
		if districts, err := c.location.GetDistrictsByIDs(ctx, ids); err == nil {
			for i := range districts {
				d := &districts[i]
				parent := d.ProvinceID
				out[d.ID] = &hubpb.Location{
					Id: d.ID, Name: d.Name, Type: d.Type, TypeText: d.TypeText, Level: 2, ParentId: &parent,
				}
			}
		}
		if wards, err := c.location.GetWardsByIDs(ctx, ids); err == nil {
			for i := range wards {
				w := &wards[i]
				parent := w.DistrictID
				out[w.ID] = &hubpb.Location{
					Id: w.ID, Name: w.Name, Type: w.Type, TypeText: w.TypeText, Level: 3, ParentId: &parent,
				}
			}
		}
	}
	return out, nil
}

// GetProvinceAndWardNames lấy tên province và ward (v2) theo IDs.
func (c *HubClient) GetProvinceAndWardNames(ctx context.Context, provinceId, wardId *uint64) (provinceName, wardName string) {
	if c == nil || c.locationV2 == nil {
		return "", ""
	}
	if provinceId != nil && *provinceId > 0 {
		if p, err := c.locationV2.GetProvinceByID(ctx, *provinceId); err == nil && p != nil {
			provinceName = p.Name
		}
	}
	if wardId != nil && *wardId > 0 {
		if w, err := c.locationV2.GetWardByID(ctx, *wardId); err == nil && w != nil {
			wardName = w.Name
		}
	}
	return provinceName, wardName
}

// GetUserGuideByKey lấy user guide theo key từ local DB.
func (c *HubClient) GetUserGuideByKey(ctx context.Context, key string) (*hubpb.UserGuideDetail, error) {
	if c == nil || c.userGuide == nil {
		return nil, fmt.Errorf("user guide store not available")
	}
	row, err := c.userGuide.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, fmt.Errorf("user guide not found")
	}
	out := &hubpb.UserGuideDetail{
		Id: row.ID, Title: row.Title, Description: row.Description,
		GroupKey: row.GroupKey, Key: row.Key, Mode: hubpb.UserGuideMode(row.Mode),
	}
	if row.CreatedAt != nil {
		out.CreatedAt = strconv.FormatInt(row.CreatedAt.UnixMilli(), 10)
	}
	if row.UpdatedAt != nil {
		out.UpdatedAt = strconv.FormatInt(row.UpdatedAt.UnixMilli(), 10)
	}
	for _, step := range row.Steps {
		out.Steps = append(out.Steps, &hubpb.UserGuideStep{
			Id: step.ID, UserGuideId: step.UserGuideID,
			StepOrder: int32(step.StepOrder), Image: step.Image, Content: step.Content,
		})
	}
	return out, nil
}
