package handler_grpc

import (
	_utils "common/utils"
	"context"
	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"time"
	"tqd/internal/enums"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

type TqdPublicService struct {
	tqdpb.UnimplementedTqdPublicServiceServer
	StartUpTime int64
}

func NewBdsproPublicService() *TqdPublicService {
	return &TqdPublicService{

		StartUpTime: time.Now().UnixMilli(),
	}
}

var enumProviders = map[string][]*sharepb.EnumItem{
	"problem_report": _utils.EnumMapToList(enums.EProblemReportNames),
	"legal_status":   _utils.EnumMapToList(enums.LegalStatusMap),
	"legal_level":    _utils.EnumMapToList(enums.LegalLevelMap),
	"scope_apply":    _utils.EnumMapToList(enums.ScopeApplyNames),
}

func (s *TqdPublicService) GetEnums(ctx context.Context, req *sharepb.SyncRequest) (*sharepb.GetEnumsResponse, error) {
	if req.Timestamp > 0 && req.Timestamp < s.StartUpTime {
		return &sharepb.GetEnumsResponse{
			Timestamp: s.StartUpTime,
			Lastest:   false,
		}, nil
	}
	result := []*sharepb.GetEnumResponse{}
	for k, v := range enumProviders {
		result = append(result, &sharepb.GetEnumResponse{
			Key:  k,
			Data: v,
		})
	}
	return &sharepb.GetEnumsResponse{
		Enums:     result,
		Timestamp: s.StartUpTime,
		Lastest:   true,
	}, nil
}

func (s *TqdPublicService) GetQhConfigApp(ctx context.Context, _ *emptypb.Empty) (*structpb.Struct, error) {
	payload := map[string]interface{}{
		"backgroundMaps": []interface{}{
			map[string]interface{}{"id": "basemap-light", "name": "Bản đồ nền sáng", "image": ""},
			map[string]interface{}{"id": "basemap-dark", "name": "Bản đồ nền tối", "image": ""},
			map[string]interface{}{"id": "satellite", "name": "Ảnh vệ tinh", "image": ""},
		},
		"qhLayers": []interface{}{
			map[string]interface{}{"id": "planning-use", "name": "Quy hoạch sử dụng đất", "image": ""},
			map[string]interface{}{"id": "planning-transport", "name": "Quy hoạch giao thông", "image": ""},
		},
	}
	st, err := structpb.NewStruct(payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "qh config-app: %v", err)
	}
	return st, nil
}
