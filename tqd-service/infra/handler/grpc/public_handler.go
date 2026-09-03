package handler_grpc

import (
	_utils "common/utils"
	"context"
	"fmt"
	"log"
	"os"
	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"regexp"
	"strings"
	"sync"
	"time"
	"tqd/config"
	"tqd/internal/enums"
	atlasflowcatalog "tqd/internal/usecase/atlasflow/catalog"

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

	"planning_document_type": _utils.EnumMapToList(enums.PlanningDocumentTypeLabelMap),
	"planning_type":          _utils.EnumMapToList(enums.PlanningTypeMap),
	"planning_level":         _utils.EnumMapToList(enums.PlanningLevelMap),
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

var qhConfigNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func (s *TqdPublicService) GetQhConfigApp(ctx context.Context, req *tqdpb.GetQhConfigAppRequest) (*tqdpb.GetQhConfigAppResponse, error) {
	if req == nil || !qhConfigNamePattern.MatchString(req.Config) {
		return nil, status.Error(codes.InvalidArgument, "invalid config name")
	}

	jsonPath, err := s.ensureQhConfigProcessed(req.Config)
	if err != nil {
		return nil, err
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

// ensureQhConfigProcessed owns template expansion and atomic publication.
//
// Configuration endpoints are read frequently by map clients, while template
// expansion is required only once per process. The single mutex prevents two
// concurrent requests from writing the same generated JSON file, and the
// processed flag is committed only after the atomic rename succeeds.
func (s *TqdPublicService) ensureQhConfigProcessed(configName string) (string, error) {
	jsonPath := fmt.Sprintf("files/json/%s.json", configName)

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.processedConfigs[configName] {
		if _, err := os.Stat(jsonPath); err == nil {
			return jsonPath, nil
		}
		// A generated file may have been removed outside the process. Rebuild it
		// instead of keeping a stale in-memory success flag.
		delete(s.processedConfigs, configName)
	}

	templatePath := fmt.Sprintf("template/%s.json", configName)
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		log.Printf("failed to read template file: %v", err)
		return "", status.Errorf(codes.NotFound, "template file not found: %v", err)
	}

	qhpro := config.AppConfig.ConfigApp.Qhpro
	content := string(templateData)
	content = strings.ReplaceAll(content, "{{domain_api}}", qhpro.DomainAPI)
	content = strings.ReplaceAll(content, "{{domain_file}}", qhpro.DomainFile)
	content = strings.ReplaceAll(content, "{{domain_pmtile}}", qhpro.DomainPMTile)
	content = strings.ReplaceAll(content, "{{domain_style}}", qhpro.DomainStyle)

	if configName == atlasflowcatalog.TemplateName {
		if err := atlasflowcatalog.ValidateJSON([]byte(content)); err != nil {
			log.Printf("invalid AtlasFlow map catalog: %v", err)
			return "", status.Errorf(codes.Internal, "invalid AtlasFlow map catalog: %v", err)
		}
	}

	if err := os.MkdirAll("files/json", 0755); err != nil {
		log.Printf("failed to create json directory: %v", err)
		return "", status.Errorf(codes.Internal, "failed to create json directory: %v", err)
	}

	temporaryFile, err := os.CreateTemp("files/json", configName+"-*.tmp")
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to create temporary config: %v", err)
	}
	temporaryPath := temporaryFile.Name()
	defer os.Remove(temporaryPath)

	if _, err := temporaryFile.Write([]byte(content)); err != nil {
		temporaryFile.Close()
		return "", status.Errorf(codes.Internal, "failed to write temporary config: %v", err)
	}
	if err := temporaryFile.Chmod(0644); err != nil {
		temporaryFile.Close()
		return "", status.Errorf(codes.Internal, "failed to set config permissions: %v", err)
	}
	if err := temporaryFile.Sync(); err != nil {
		temporaryFile.Close()
		return "", status.Errorf(codes.Internal, "failed to sync temporary config: %v", err)
	}
	if err := temporaryFile.Close(); err != nil {
		return "", status.Errorf(codes.Internal, "failed to close temporary config: %v", err)
	}
	if err := os.Rename(temporaryPath, jsonPath); err != nil {
		return "", status.Errorf(codes.Internal, "failed to publish processed config: %v", err)
	}

	s.processedConfigs[configName] = true
	return jsonPath, nil
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
	qhpro := config.AppConfig.ConfigApp.Qhpro
	return &tqdpb.GetPublicMapsResponse{
		Data: []*tqdpb.MapSource{
			{
				Name:   "Giao thông",
				Source: strings.TrimRight(qhpro.DomainAPI, "/") + "/v2/tqd/qh/public/json/style",
				Type:   "style",
			},
			{
				Name:   "Vệ tinh",
				Source: "https://mt1.google.com/vt/lyrs=s&x={x}&y={y}&z={z}",
				Type:   "raster",
			},
		},
	}, nil
}
