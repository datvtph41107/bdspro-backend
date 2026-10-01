package services

import (
	"context"
	"time"

	_utils "common/utils"
	sharepb "pb/types/shared"

	"userv4/modules/tqd/enums"
)

const syncKeyTqdEnums = "tqd:enums"

// PublicService owns guest/public TQD catalog flows (enums, later maps/config).
type PublicService interface {
	GetEnums(ctx context.Context, timestamp int64) (*sharepb.GetEnumsResponse, error)
}

// LocalPublicService serves TQD public enum catalogs from in-process maps.
type LocalPublicService struct {
	sync      *_utils.SyncUtil
	startTime int64
}

func NewLocalPublicService(sync *_utils.SyncUtil) *LocalPublicService {
	return &LocalPublicService{
		sync:      sync,
		startTime: time.Now().UnixMilli(),
	}
}

func (s *LocalPublicService) GetEnums(ctx context.Context, timestamp int64) (*sharepb.GetEnumsResponse, error) {
	catalogTS := s.startTime
	if s.sync != nil {
		catalogTS = s.sync.EnsureTimestamp(ctx, syncKeyTqdEnums)
	}

	updated := true
	if s.sync != nil {
		updated = s.sync.HasUpdated(ctx, syncKeyTqdEnums, timestamp)
	} else if timestamp > 0 && timestamp >= catalogTS {
		updated = false
	}
	if !updated {
		return &sharepb.GetEnumsResponse{
			Timestamp: catalogTS,
			Lastest:   false,
		}, nil
	}

	result := make([]*sharepb.GetEnumResponse, 0, len(enumProviders))
	for k, v := range enumProviders {
		result = append(result, &sharepb.GetEnumResponse{
			Key:  k,
			Data: v,
		})
	}
	return &sharepb.GetEnumsResponse{
		Enums:     result,
		Timestamp: catalogTS,
		Lastest:   true,
	}, nil
}

// Mirrors tqd-service/infra/handler/grpc/public_handler.go enumProviders.
var enumProviders = map[string][]*sharepb.EnumItem{
	"problem_report": _utils.EnumMapToList(enums.EProblemReportNames),
	"legal_status":   _utils.EnumMapToList(enums.LegalStatusMap),
	"legal_level":    _utils.EnumMapToList(enums.LegalLevelMap),
	"scope_apply":    _utils.EnumMapToList(enums.ScopeApplyNames),

	"planning_document_type": _utils.EnumMapToList(enums.PlanningDocumentTypeLabelMap),
	"planning_type":          _utils.EnumMapToList(enums.PlanningTypeMap),
	"planning_level":         _utils.EnumMapToList(enums.PlanningLevelMap),

	"location_unit_type":   _utils.EnumMapToList(enums.LocationUnitTypeMap),
	"location_data_source": _utils.EnumMapToList(enums.LocationDataSourceMap),
}
