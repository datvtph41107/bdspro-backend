package handler_grpc

import (
	_utils "common/utils"
	"context"
	"fmt"
	"log"
	"os"
	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"strings"
	"sync"
	"time"
	"tqd/config"
	"tqd/internal/enums"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

type TqdPublicService struct {
	tqdpb.UnimplementedTqdPublicServiceServer
	StartUpTime      int64
	processedConfigs map[string]bool
	mu               sync.RWMutex
}

func NewBdsproPublicService() *TqdPublicService {
	return &TqdPublicService{
		StartUpTime:      time.Now().UnixMilli(),
		processedConfigs: make(map[string]bool),
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

func (s *TqdPublicService) GetQhConfigApp(ctx context.Context, req *tqdpb.GetQhConfigAppRequest) (*tqdpb.GetQhConfigAppResponse, error) {
	s.mu.RLock()
	processed := s.processedConfigs[req.Config]
	s.mu.RUnlock()

	jsonPath := fmt.Sprintf("files/json/%s.json", req.Config)

	if !processed {
		templatePath := fmt.Sprintf("template/%s.json", req.Config)
		templateData, err := os.ReadFile(templatePath)
		if err != nil {
			log.Printf("failed to read template file: %v", err)
			return nil, status.Errorf(codes.NotFound, "template file not found: %v", err)
		}

		// Load domains from hub-service/config/develop.yml
		qhpro := config.AppConfig.ConfigApp.Qhpro
		content := string(templateData)
		content = strings.ReplaceAll(content, "{{domain_api}}", qhpro.DomainAPI)
		content = strings.ReplaceAll(content, "{{domain_file}}", qhpro.DomainFile)
		content = strings.ReplaceAll(content, "{{domain_pmtile}}", qhpro.DomainPMTile)
		content = strings.ReplaceAll(content, "{{domain_style}}", qhpro.DomainStyle)

		s.mu.Lock()
		s.processedConfigs[req.Config] = true
		s.mu.Unlock()

		if err := os.MkdirAll("files/json", 0755); err != nil {
			log.Printf("failed to create json directory: %v", err)
			return nil, status.Errorf(codes.Internal, "failed to create json directory: %v", err)
		}

		if err := os.WriteFile(jsonPath, []byte(content), 0644); err != nil {
			log.Printf("failed to save processed config: %v", err)
			return nil, status.Errorf(codes.Internal, "failed to save processed config: %v", err)
		}
	}

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "config file not found: %v", err)
	}

	st := &structpb.Struct{}
	if err := st.UnmarshalJSON(data); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshal config: %v", err)
	}
	return &tqdpb.GetQhConfigAppResponse{
		Code: 8500,
		Data: st,
	}, nil
}

func (s *TqdPublicService) GetFont(ctx context.Context, req *tqdpb.GetFontRequest) (*sharepb.CommonResponse, error) {
	path := fmt.Sprintf("files/fonts/%s/%s", req.Fontstack, req.Range)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "font file not found: %v", err)
	}

	return &sharepb.CommonResponse{
		Code: 8501,
		Data: data,
	}, nil
}

func (s *TqdPublicService) loadDomains() (map[string]string, error) {

	qhpro := config.AppConfig.ConfigApp.Qhpro
	result := make(map[string]string)
	result["domain_api"] = qhpro.DomainAPI
	result["domain_file"] = qhpro.DomainFile
	result["domain_pmtile"] = qhpro.DomainPMTile
	result["domain_style"] = qhpro.DomainStyle
	return result, nil
}

func (s *TqdPublicService) GetPublicMaps(ctx context.Context, _ *emptypb.Empty) (*tqdpb.GetPublicMapsResponse, error) {
	return &tqdpb.GetPublicMapsResponse{
		Data: []*tqdpb.MapSource{
			{
				Name:   "Giao thông",
				Source: "https://data.bdspro.com/v1/map/config/style.json",
				Type:   "vector",
			},
			{
				Name:   "Vệ tinh",
				Source: "https://mt1.google.com/vt/lyrs=s&x={x}&y={y}&z={z}",
				Type:   "raster",
			},
		},
	}, nil
}
