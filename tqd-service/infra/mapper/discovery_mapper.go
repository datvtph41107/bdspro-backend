package mapper

import (
	"encoding/json"
	"fmt"

	tqdpb "pb/types/tqd"
	"tqd/internal/domain/discovery/model"

	"google.golang.org/protobuf/types/known/structpb"
)

// DiscoveryMapper quản lý toàn bộ conversion giữa:
//
//   - generated Discovery protobuf contract;
//   - Discovery domain model.
//
// Handler không được trực tiếp biết cách ánh xạ enum, spatial data,
// EntityCandidate, attributes hoặc aggregate response.
type DiscoveryMapper struct{}

// NewDiscoveryMapper khởi tạo Discovery mapper.
func NewDiscoveryMapper() *DiscoveryMapper {
	return &DiscoveryMapper{}
}

// FromProtoIdentifyRequest chuyển typed protobuf request thành domain request
// dùng bởi Discovery application service.
func (m *DiscoveryMapper) FromProtoIdentifyRequest(
	input *tqdpb.IdentifyAtPointRequest,
) domain.IdentifyRequest {
	request := domain.IdentifyRequest{}

	if input == nil {
		return request
	}

	if point := input.GetPoint(); point != nil {
		request.Point = domain.Point{
			Latitude:  point.GetLatitude(),
			Longitude: point.GetLongitude(),
		}
	}

	if contextValue := input.GetContext(); contextValue != nil {
		request.Context = domain.MapContext{
			Zoom:    contextValue.GetZoom(),
			MapMode: contextValue.GetMapMode(),
			ActiveLayerIDs: cloneStrings(
				contextValue.GetActiveLayerIds(),
			),
			PreferredKinds: m.entityKindsFromProto(
				contextValue.GetPreferredKinds(),
			),
			Viewport: m.boundsFromProto(
				contextValue.GetViewport(),
			),
		}
	}

	if options := input.GetOptions(); options != nil {
		request.Options = domain.IdentifyOptions{
			ToleranceMeters: options.GetToleranceMeters(),
			Limit:           int(options.GetLimit()),
			IncludeGeometryPreview: options.
				GetIncludeGeometryPreview(),
		}
	}

	return request
}

// FromProtoSearchRequest chuyển typed protobuf SearchDiscoveryRequest
// thành domain SearchRequest.
func (m *DiscoveryMapper) FromProtoSearchRequest(
	input *tqdpb.SearchDiscoveryRequest,
) domain.SearchRequest {
	if input == nil {
		return domain.SearchRequest{}
	}

	return domain.SearchRequest{
		Query: input.GetQ(),
		Mode:  input.GetMode(),
		EntityKinds: m.entityKindsFromProto(
			input.GetKinds(),
		),
		MapContext: domain.MapContext{
			MapMode: input.GetMapMode(),
			ActiveLayerIDs: cloneStrings(
				input.GetActiveLayerIds(),
			),
		},
		Cursor: input.GetCursor(),
		Limit:  int(input.GetLimit()),
	}
}

// ToProtoIdentifyResponse chuyển toàn bộ IdentifyResponse domain
// thành typed protobuf response.
func (m *DiscoveryMapper) ToProtoIdentifyResponse(
	value *domain.IdentifyResponse,
) (*tqdpb.IdentifyAtPointResponse, error) {
	if value == nil {
		return nil, fmt.Errorf("identify response is nil")
	}

	candidates, err := m.candidatesToProto(value.Candidates)
	if err != nil {
		return nil, fmt.Errorf(
			"map identify candidates: %w",
			err,
		)
	}

	var primary *tqdpb.EntityCandidate

	if value.Primary != nil {
		primary, err = m.ToProtoEntityCandidate(value.Primary)
		if err != nil {
			return nil, fmt.Errorf(
				"map identify primary candidate: %w",
				err,
			)
		}
	}

	return &tqdpb.IdentifyAtPointResponse{
		RequestId: value.RequestID,
		Resolution: &tqdpb.DiscoveryResolution{
			Status:     value.Resolution.Status,
			PrimaryKey: value.Resolution.PrimaryKey,
			Confidence: value.Resolution.Confidence,
			Ambiguous:  value.Resolution.Ambiguous,
			Strategy:   value.Resolution.Strategy,
		},
		Primary:    primary,
		Candidates: candidates,
		Partial:    value.Partial,
		Warnings:   cloneStrings(value.Warnings),
	}, nil
}

// ToProtoSearchResponse chuyển grouped domain search results
// thành typed protobuf response.
func (m *DiscoveryMapper) ToProtoSearchResponse(
	value *domain.SearchResponse,
) (*tqdpb.SearchDiscoveryResponse, error) {
	if value == nil {
		return nil, fmt.Errorf("search response is nil")
	}

	groups := make(
		[]*tqdpb.SearchGroup,
		0,
		len(value.Groups),
	)

	for index := range value.Groups {
		group := &value.Groups[index]

		items, err := m.candidatesToProto(group.Items)
		if err != nil {
			return nil, fmt.Errorf(
				"map search group %q: %w",
				group.Label,
				err,
			)
		}

		groups = append(
			groups,
			&tqdpb.SearchGroup{
				Kind:  m.entityKindToProto(group.Kind),
				Label: group.Label,
				Total: int32(group.Total),
				Items: items,
			},
		)
	}

	return &tqdpb.SearchDiscoveryResponse{
		RequestId: value.RequestID,
		Interpretation: &tqdpb.SearchInterpretation{
			OriginalQuery: value.Interpretation.
				OriginalQuery,
			NormalizedQuery: value.Interpretation.
				NormalizedQuery,
			DetectedIntents: cloneStrings(
				value.Interpretation.DetectedIntents,
			),
			StructuredTerms: cloneStringMap(
				value.Interpretation.StructuredTerms,
			),
		},
		Groups:     groups,
		NextCursor: value.NextCursor,
		Partial:    value.Partial,
		Warnings:   cloneStrings(value.Warnings),
	}, nil
}

// ToProtoEntityCandidate chuyển một domain EntityCandidate thành
// protobuf EntityCandidate.
//
// Method được public vì GetDiscoveryEntity trả trực tiếp EntityCandidate.
func (m *DiscoveryMapper) ToProtoEntityCandidate(
	value *domain.EntityCandidate,
) (*tqdpb.EntityCandidate, error) {
	if value == nil {
		return nil, fmt.Errorf("entity candidate is nil")
	}

	if !value.Entity.Kind.IsValid() {
		return nil, fmt.Errorf(
			"invalid entity kind %q",
			value.Entity.Kind,
		)
	}

	if value.Entity.ID == "" {
		return nil, fmt.Errorf(
			"entity id is required",
		)
	}

	if value.Entity.Key == "" {
		return nil, fmt.Errorf(
			"entity canonical key is required",
		)
	}

	attributes, err := structFromObject(value.Attributes)
	if err != nil {
		return nil, fmt.Errorf(
			"encode attributes for %s: %w",
			value.Entity.Key,
			err,
		)
	}

	spatial, err := m.spatialToProto(
		value.Entity.Key,
		value.Spatial,
	)
	if err != nil {
		return nil, err
	}

	return &tqdpb.EntityCandidate{
		Entity: &tqdpb.EntityRef{
			Kind: m.entityKindToProto(
				value.Entity.Kind,
			),
			Id:  value.Entity.ID,
			Key: value.Entity.Key,
		},

		Presentation: &tqdpb.EntityPresentation{
			Title:       value.Presentation.Title,
			Subtitle:    value.Presentation.Subtitle,
			Description: value.Presentation.Description,
			Status:      value.Presentation.Status,
			Badges: cloneStrings(
				value.Presentation.Badges,
			),
		},

		Spatial: m.mustSpatialValue(spatial),

		Location: m.locationToProto(
			value.Location,
		),

		Source: &tqdpb.EntitySource{
			System:      value.Source.System,
			Dataset:     value.Source.Dataset,
			Authority:   value.Source.Authority,
			UpdatedAt:   value.Source.UpdatedAt,
			DataQuality: value.Source.DataQuality,
		},

		Capabilities: &tqdpb.EntityCapabilities{
			CanOpenDetail: value.Capabilities.
				CanOpenDetail,
			CanFocusMap: value.Capabilities.
				CanFocusMap,
			CanFollow: value.Capabilities.
				CanFollow,
			CanCreateReport: value.Capabilities.
				CanCreateReport,
			CanCompare: value.Capabilities.
				CanCompare,
			CanShare: value.Capabilities.
				CanShare,
			CanDownload: value.Capabilities.
				CanDownload,
		},

		Links: &tqdpb.EntityLinks{
			CanonicalUrl: value.Links.CanonicalURL,
			DetailUrl:    value.Links.DetailURL,
			GeometryUrl:  value.Links.GeometryURL,
		},

		Match: &tqdpb.EntityMatch{
			Type: value.Match.Type,
			MatchedFields: cloneStrings(
				value.Match.MatchedFields,
			),
			Highlights: cloneStrings(
				value.Match.Highlights,
			),
			ReasonCodes: cloneStrings(
				value.Match.ReasonCodes,
			),
		},

		Rank: &tqdpb.EntityRank{
			Score: value.Rank.Score,
			ReasonCodes: cloneStrings(
				value.Rank.ReasonCodes,
			),
		},

		Attributes: attributes,
	}, nil
}

// entityKindFromProto chuyển protobuf EntityKind thành domain EntityKind.
func (m *DiscoveryMapper) entityKindFromProto(
	kind tqdpb.EntityKind,
) domain.EntityKind {
	switch kind {
	case tqdpb.EntityKind_ENTITY_KIND_PARCEL:
		return domain.EntityKindParcel

	case tqdpb.EntityKind_ENTITY_KIND_PLANNING_REGION:
		return domain.EntityKindPlanningRegion

	case tqdpb.EntityKind_ENTITY_KIND_ADMINISTRATIVE_UNIT:
		return domain.EntityKindAdministrativeUnit

	case tqdpb.EntityKind_ENTITY_KIND_POI:
		return domain.EntityKindPOI

	case tqdpb.EntityKind_ENTITY_KIND_PLANNING_PROJECT:
		return domain.EntityKindPlanningProject

	case tqdpb.EntityKind_ENTITY_KIND_LEGAL_DOCUMENT:
		return domain.EntityKindLegalDocument

	case tqdpb.EntityKind_ENTITY_KIND_PLANNING_MAP:
		return domain.EntityKindPlanningMap

	case tqdpb.EntityKind_ENTITY_KIND_MAP_LAYER:
		return domain.EntityKindMapLayer

	case tqdpb.EntityKind_ENTITY_KIND_ANALYSIS_REPORT:
		return domain.EntityKindAnalysisReport

	case tqdpb.EntityKind_ENTITY_KIND_NEWS_ARTICLE:
		return domain.EntityKindNewsArticle

	default:
		return ""
	}
}

// entityKindToProto chuyển domain EntityKind thành protobuf EntityKind.
func (m *DiscoveryMapper) entityKindToProto(
	kind domain.EntityKind,
) tqdpb.EntityKind {
	switch kind {
	case domain.EntityKindParcel:
		return tqdpb.EntityKind_ENTITY_KIND_PARCEL

	case domain.EntityKindPlanningRegion:
		return tqdpb.EntityKind_ENTITY_KIND_PLANNING_REGION

	case domain.EntityKindAdministrativeUnit:
		return tqdpb.EntityKind_ENTITY_KIND_ADMINISTRATIVE_UNIT

	case domain.EntityKindPOI:
		return tqdpb.EntityKind_ENTITY_KIND_POI

	case domain.EntityKindPlanningProject:
		return tqdpb.EntityKind_ENTITY_KIND_PLANNING_PROJECT

	case domain.EntityKindLegalDocument:
		return tqdpb.EntityKind_ENTITY_KIND_LEGAL_DOCUMENT

	case domain.EntityKindPlanningMap:
		return tqdpb.EntityKind_ENTITY_KIND_PLANNING_MAP

	case domain.EntityKindMapLayer:
		return tqdpb.EntityKind_ENTITY_KIND_MAP_LAYER

	case domain.EntityKindAnalysisReport:
		return tqdpb.EntityKind_ENTITY_KIND_ANALYSIS_REPORT

	case domain.EntityKindNewsArticle:
		return tqdpb.EntityKind_ENTITY_KIND_NEWS_ARTICLE

	default:
		return tqdpb.EntityKind_ENTITY_KIND_UNSPECIFIED
	}
}

// entityKindsFromProto chuyển danh sách protobuf EntityKind
// thành danh sách domain EntityKind.
//
// Giá trị ENTITY_KIND_UNSPECIFIED hoặc enum không hỗ trợ sẽ bị bỏ qua.
func (m *DiscoveryMapper) entityKindsFromProto(
	values []tqdpb.EntityKind,
) []domain.EntityKind {
	if len(values) == 0 {
		return nil
	}

	result := make(
		[]domain.EntityKind,
		0,
		len(values),
	)

	seen := make(map[domain.EntityKind]struct{}, len(values))

	for _, value := range values {
		kind := m.entityKindFromProto(value)
		if kind == "" {
			continue
		}

		if _, exists := seen[kind]; exists {
			continue
		}

		seen[kind] = struct{}{}
		result = append(result, kind)
	}

	return result
}

// boundsFromProto chuyển protobuf DiscoveryBounds thành domain Bounds.
func (m *DiscoveryMapper) boundsFromProto(
	value *tqdpb.DiscoveryBounds,
) *domain.Bounds {
	if value == nil {
		return nil
	}

	return &domain.Bounds{
		MinLongitude: value.GetMinLongitude(),
		MinLatitude:  value.GetMinLatitude(),
		MaxLongitude: value.GetMaxLongitude(),
		MaxLatitude:  value.GetMaxLatitude(),
	}
}

// candidatesToProto chuyển danh sách domain candidate sang protobuf.
func (m *DiscoveryMapper) candidatesToProto(
	values []domain.EntityCandidate,
) ([]*tqdpb.EntityCandidate, error) {
	if values == nil {
		return nil, nil
	}

	result := make(
		[]*tqdpb.EntityCandidate,
		0,
		len(values),
	)

	for index := range values {
		candidate, err := m.ToProtoEntityCandidate(
			&values[index],
		)
		if err != nil {
			return nil, fmt.Errorf(
				"map candidate at index %d: %w",
				index,
				err,
			)
		}

		result = append(result, candidate)
	}

	return result, nil
}

// spatialToProto chuyển domain SpatialSummary thành protobuf.
func (m *DiscoveryMapper) spatialToProto(
	entityKey string,
	value *domain.SpatialSummary,
) (*tqdpb.DiscoverySpatialSummary, error) {
	if value == nil {
		return nil, nil
	}

	geometryPreview, err := structFromObject(
		value.GeometryPreview,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"encode geometry preview for %s: %w",
			entityKey,
			err,
		)
	}

	return &tqdpb.DiscoverySpatialSummary{
		GeometryType: value.GeometryType,
		Centroid: pointToProto(
			value.Centroid,
		),
		Bounds: boundsToProto(
			value.Bounds,
		),
		GeometryRef:     value.GeometryRef,
		GeometryPreview: geometryPreview,
	}, nil
}

// mustSpatialValue chỉ giữ code khởi tạo candidate dễ đọc.
// Giá trị nil vẫn hợp lệ đối với entity không có spatial data.
func (m *DiscoveryMapper) mustSpatialValue(
	value *tqdpb.DiscoverySpatialSummary,
) *tqdpb.DiscoverySpatialSummary {
	return value
}

// locationToProto chuyển domain Location thành protobuf EntityLocation.
func (m *DiscoveryMapper) locationToProto(
	value *domain.Location,
) *tqdpb.EntityLocation {
	if value == nil {
		return nil
	}

	return &tqdpb.EntityLocation{
		Address:      value.Address,
		Province:     value.Province,
		ProvinceCode: value.ProvinceCode,
		Ward:         value.Ward,
		WardCode:     value.WardCode,
	}
}

// pointToProto chuyển domain Point thành protobuf DiscoveryPoint.
func pointToProto(
	value *domain.Point,
) *tqdpb.DiscoveryPoint {
	if value == nil {
		return nil
	}

	return &tqdpb.DiscoveryPoint{
		Latitude:  value.Latitude,
		Longitude: value.Longitude,
	}
}

// boundsToProto chuyển domain Bounds thành protobuf DiscoveryBounds.
func boundsToProto(
	value *domain.Bounds,
) *tqdpb.DiscoveryBounds {
	if value == nil {
		return nil
	}

	return &tqdpb.DiscoveryBounds{
		MinLongitude: value.MinLongitude,
		MinLatitude:  value.MinLatitude,
		MaxLongitude: value.MaxLongitude,
		MaxLatitude:  value.MaxLatitude,
	}
}

// structFromObject chuẩn hóa repository-specific aliases,
// json.RawMessage và các giá trị JSON-serializable trước khi tạo
// google.protobuf.Struct.
//
// Trường protobuf đích yêu cầu object. Scalar hoặc array sẽ bị từ chối
// thay vì âm thầm đổi cấu trúc dữ liệu.
func structFromObject(
	value any,
) (*structpb.Struct, error) {
	if value == nil {
		return nil, nil
	}

	raw, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf(
			"marshal object: %w",
			err,
		)
	}

	var normalized map[string]any

	if err := json.Unmarshal(raw, &normalized); err != nil {
		return nil, fmt.Errorf(
			"decode object: %w",
			err,
		)
	}

	result, err := structpb.NewStruct(normalized)
	if err != nil {
		return nil, fmt.Errorf(
			"construct protobuf struct: %w",
			err,
		)
	}

	return result, nil
}

// cloneStrings tạo slice mới để protobuf response không chia sẻ
// backing array với domain model.
func cloneStrings(
	values []string,
) []string {
	if values == nil {
		return nil
	}

	return append([]string(nil), values...)
}

// cloneStringMap tạo map mới để protobuf response không chia sẻ
// map reference với domain model.
func cloneStringMap(
	values map[string]string,
) map[string]string {
	if values == nil {
		return nil
	}

	result := make(
		map[string]string,
		len(values),
	)

	for key, value := range values {
		result[key] = value
	}

	return result
}
