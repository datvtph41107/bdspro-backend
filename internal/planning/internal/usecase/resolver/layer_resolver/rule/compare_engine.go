package rule

import (
	"context"
	"fmt"

	"tqd/internal/usecase/resolver/layer_resolver/types"
)

// CompareRepository — interface cho Compare Engine
type CompareRepository interface {
	GetParcelLayerRowsAt(ctx context.Context, parcelID uint64, asOfDate string) ([]*CompareRow, error)
}

// CompareRow — dữ liệu layer tại 1 thời điểm cho compare
type CompareRow struct {
	LayerID          uint64
	LayerName        string
	LayerDisplayName string
	LegalStatus      uint32
	OverlapPct       float64
	OverlapAreaSqm   float64
	LandUseCode      string
	LandUseName      string
}

// CompareEngine — Engine so sánh layer giữa 2 thời điểm
type CompareEngine struct {
	config      *types.CompareConfig
	compareRepo CompareRepository
}

// NewCompareEngine tạo CompareEngine
func NewCompareEngine(config *types.CompareConfig, compareRepo CompareRepository) *CompareEngine {
	if config == nil {
		config = types.DefaultCompareConfig()
	}
	return &CompareEngine{config: config, compareRepo: compareRepo}
}

// Name trả về tên engine
func (e *CompareEngine) Name() string { return "CompareEngine" }

// Compare thực hiện so sánh giữa 2 thời điểm
func (e *CompareEngine) Compare(ctx context.Context, req *types.CompareRequest) (*types.CompareResult, error) {
	baselineRows, err := e.compareRepo.GetParcelLayerRowsAt(ctx, req.ParcelID, req.BaselineAt)
	if err != nil {
		return nil, fmt.Errorf("get baseline rows: %w", err)
	}

	targetRows, err := e.compareRepo.GetParcelLayerRowsAt(ctx, req.ParcelID, req.TargetAt)
	if err != nil {
		return nil, fmt.Errorf("get target rows: %w", err)
	}

	baselineMap := e.buildRowMap(baselineRows)
	targetMap := e.buildRowMap(targetRows)

	var changes []*types.LayerChange
	allLayerIDs := e.mergeKeys(baselineMap, targetMap)

	for _, layerID := range allLayerIDs {
		base, inBaseline := baselineMap[layerID]
		target, inTarget := targetMap[layerID]

		switch {
		case inBaseline && !inTarget:
			changes = append(changes, &types.LayerChange{
				LayerID:    layerID,
				LayerName:  base.LayerName,
				ChangeType: types.ChangeRemoved,
				ChangeDetail: &types.ChangeDetail{
					ImpactSummary: "Layer đã bị thu hồi hoặc hết hiệu lực",
				},
			})
		case !inBaseline && inTarget:
			changes = append(changes, &types.LayerChange{
				LayerID:    layerID,
				LayerName:  target.LayerName,
				ChangeType: types.ChangeAdded,
				ChangeDetail: &types.ChangeDetail{
					ImpactSummary: "Layer mới được ban hành",
				},
			})
		case inBaseline && inTarget:
			detail := e.detectModifications(base, target)
			if detail != nil {
				changes = append(changes, &types.LayerChange{
					LayerID:      layerID,
					LayerName:    target.LayerName,
					ChangeType:   types.ChangeModified,
					ChangeDetail: detail,
				})
			} else {
				changes = append(changes, &types.LayerChange{
					LayerID:    layerID,
					LayerName:  target.LayerName,
					ChangeType: types.ChangeUnchanged,
				})
			}
		}
	}

	summary := e.buildSummary(changes)

	return &types.CompareResult{
		ParcelID:          req.ParcelID,
		BaselineAt:        req.BaselineAt,
		TargetAt:          req.TargetAt,
		TotalLayersBefore: len(baselineMap),
		TotalLayersAfter:  len(targetMap),
		Changes:           changes,
		Summary:           summary,
	}, nil
}

// buildRowMap xây dựng map layerID → row
func (e *CompareEngine) buildRowMap(rows []*CompareRow) map[uint64]*CompareRow {
	m := make(map[uint64]*CompareRow, len(rows))
	for _, r := range rows {
		m[r.LayerID] = r
	}
	return m
}

// mergeKeys hợp nhất tất cả layerID từ 2 map
func (e *CompareEngine) mergeKeys(a, b map[uint64]*CompareRow) []uint64 {
	seen := make(map[uint64]bool)
	var keys []uint64
	for k := range a {
		if !seen[k] {
			seen[k] = true
			keys = append(keys, k)
		}
	}
	for k := range b {
		if !seen[k] {
			seen[k] = true
			keys = append(keys, k)
		}
	}
	return keys
}

// detectModifications phát hiện thay đổi giữa 2 version của cùng 1 layer
func (e *CompareEngine) detectModifications(base, target *CompareRow) *types.ChangeDetail {
	var fieldChanges []types.FieldChange

	if base.LegalStatus != target.LegalStatus {
		fieldChanges = append(fieldChanges, types.FieldChange{
			Field:    "legal_status",
			OldValue: fmt.Sprintf("%d", base.LegalStatus),
			NewValue: fmt.Sprintf("%d", target.LegalStatus),
		})
	}
	if absDiff(base.OverlapPct, target.OverlapPct) > e.config.SignificantChangeThreshold {
		fieldChanges = append(fieldChanges, types.FieldChange{
			Field:    "overlap_percent",
			OldValue: fmt.Sprintf("%.1f%%", base.OverlapPct),
			NewValue: fmt.Sprintf("%.1f%%", target.OverlapPct),
		})
	}
	if base.LandUseCode != target.LandUseCode {
		fieldChanges = append(fieldChanges, types.FieldChange{
			Field:    "land_use_code",
			OldValue: base.LandUseCode,
			NewValue: target.LandUseCode,
		})
	}

	if len(fieldChanges) == 0 {
		return nil
	}

	return &types.ChangeDetail{
		FieldChanges:  fieldChanges,
		OverlapDelta:  target.OverlapPct - base.OverlapPct,
		ImpactSummary: fmt.Sprintf("Thay đổi %d thuộc tính", len(fieldChanges)),
	}
}

// buildSummary xây dựng tóm tắt so sánh
func (e *CompareEngine) buildSummary(changes []*types.LayerChange) *types.CompareSummary {
	s := &types.CompareSummary{}
	for _, c := range changes {
		switch c.ChangeType {
		case types.ChangeAdded:
			s.AddedCount++
		case types.ChangeRemoved:
			s.RemovedCount++
		case types.ChangeModified:
			s.ModifiedCount++
		case types.ChangeUnchanged:
			s.UnchangedCount++
		}
	}
	s.NetChange = fmt.Sprintf("Thêm %d, Xóa %d, Sửa %d", s.AddedCount, s.RemovedCount, s.ModifiedCount)
	return s
}

// Validate kiểm tra config
func (e *CompareEngine) Validate() error {
	return nil
}

func absDiff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}
