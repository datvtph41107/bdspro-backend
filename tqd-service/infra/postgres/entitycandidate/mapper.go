package entitycandidate

import (
	"encoding/json"
	"regexp"
	"strings"

	discoverydomain "tqd/internal/domain/discovery/model"
)

// Row is the shared read-model used by spatial repositories when projecting a
// database row into the canonical Discovery EntityCandidate contract.
//
// Discovery and Related Entity must use the same projection rules for URLs,
// capabilities, source metadata and spatial summaries. Keeping those rules in
// one package prevents the same canonical entity from having different
// behaviour depending on which endpoint returned it.
type Row struct {
	ID             string  `gorm:"column:id"`
	Title          string  `gorm:"column:title"`
	Subtitle       string  `gorm:"column:subtitle"`
	Description    string  `gorm:"column:description"`
	Address        string  `gorm:"column:address"`
	Province       string  `gorm:"column:province"`
	ProvinceCode   string  `gorm:"column:province_code"`
	Ward           string  `gorm:"column:ward"`
	WardCode       string  `gorm:"column:ward_code"`
	GeometryType   string  `gorm:"column:geometry_type"`
	CenterLat      float64 `gorm:"column:center_lat"`
	CenterLon      float64 `gorm:"column:center_lon"`
	MinLon         float64 `gorm:"column:min_lon"`
	MinLat         float64 `gorm:"column:min_lat"`
	MaxLon         float64 `gorm:"column:max_lon"`
	MaxLat         float64 `gorm:"column:max_lat"`
	GeoJSON        string  `gorm:"column:geo_json"`
	LayerID        string  `gorm:"column:layer_id"`
	MapNumber      string  `gorm:"column:map_number"`
	LandNumber     string  `gorm:"column:land_number"`
	AreaSqm        float64 `gorm:"column:area_sqm"`
	Relationship   string  `gorm:"column:relationship"`
	DistanceMeters float64 `gorm:"column:distance_meters"`
	Score          float64 `gorm:"column:score"`
}

type BuildOptions struct {
	MatchType   string
	DataQuality string
}

func Build(
	kind discoverydomain.EntityKind,
	row Row,
	options BuildOptions,
) discoverydomain.EntityCandidate {
	ref := discoverydomain.NewEntityRef(kind, strings.TrimSpace(row.ID))
	title := strings.TrimSpace(row.Title)
	if title == "" {
		title = string(kind) + " #" + ref.ID
	}

	spatial := &discoverydomain.SpatialSummary{
		GeometryType: strings.TrimSpace(row.GeometryType),
		GeometryRef:  ref.Key,
	}
	if row.CenterLat != 0 || row.CenterLon != 0 {
		spatial.Centroid = &discoverydomain.Point{
			Latitude:  row.CenterLat,
			Longitude: row.CenterLon,
		}
	}
	if row.MinLon != 0 || row.MinLat != 0 || row.MaxLon != 0 || row.MaxLat != 0 {
		spatial.Bounds = &discoverydomain.Bounds{
			MinLongitude: row.MinLon,
			MinLatitude:  row.MinLat,
			MaxLongitude: row.MaxLon,
			MaxLatitude:  row.MaxLat,
		}
	}
	if value := GeometryPreview(row.GeoJSON); value != nil {
		spatial.GeometryPreview = value
	}

	attributes := map[string]any{}
	if row.LayerID != "" {
		attributes["layerId"] = row.LayerID
	}
	if row.MapNumber != "" {
		attributes["mapNumber"] = row.MapNumber
	}
	if row.LandNumber != "" {
		attributes["landNumber"] = row.LandNumber
	}
	if row.AreaSqm > 0 {
		attributes["areaSqm"] = row.AreaSqm
	}
	if row.DistanceMeters > 0 {
		attributes["distanceMeters"] = row.DistanceMeters
	}

	detailPath := DetailURL(kind, row)
	capabilities := discoverydomain.Capabilities{
		CanOpenDetail: detailPath != "",
		CanFocusMap:   spatial.Centroid != nil || spatial.Bounds != nil,
		CanShare:      detailPath != "",
	}
	switch kind {
	case discoverydomain.EntityKindParcel:
		capabilities.CanFollow = true
		capabilities.CanCreateReport = true
		capabilities.CanCompare = true
	case discoverydomain.EntityKindPlanningRegion:
		capabilities.CanCreateReport = true
	case discoverydomain.EntityKindPlanningProject:
		capabilities.CanCreateReport = true
		capabilities.CanCompare = true
	}

	links := discoverydomain.Links{
		DetailURL:   detailPath,
		GeometryURL: "/v2/tqd/discovery/entities/" + ref.Key,
	}
	if IsPublicCanonicalPath(detailPath) {
		links.CanonicalURL = detailPath
	}

	matchType := strings.TrimSpace(options.MatchType)
	if matchType == "" {
		matchType = "entity_projection"
	}
	dataQuality := strings.TrimSpace(options.DataQuality)
	if dataQuality == "" {
		dataQuality = "unknown"
	}

	return discoverydomain.EntityCandidate{
		Entity: ref,
		Presentation: discoverydomain.Presentation{
			Title:       title,
			Subtitle:    strings.TrimSpace(row.Subtitle),
			Description: strings.TrimSpace(row.Description),
			Status:      "unknown",
		},
		Spatial: spatial,
		Location: &discoverydomain.Location{
			Address:      strings.TrimSpace(row.Address),
			Province:     strings.TrimSpace(row.Province),
			ProvinceCode: strings.TrimSpace(row.ProvinceCode),
			Ward:         strings.TrimSpace(row.Ward),
			WardCode:     strings.TrimSpace(row.WardCode),
		},
		Source: discoverydomain.Source{
			System:      "tqd",
			Dataset:     Dataset(kind),
			Authority:   "",
			DataQuality: dataQuality,
		},
		Capabilities: capabilities,
		Links:        links,
		Match: discoverydomain.Match{
			Type:        matchType,
			ReasonCodes: []string{matchType},
		},
		Rank: discoverydomain.Rank{
			Score:       row.Score,
			ReasonCodes: []string{matchType},
		},
		Attributes: attributes,
	}
}

func GeometryPreview(raw string) any {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var value any
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return nil
	}
	return value
}

func Dataset(kind discoverydomain.EntityKind) string {
	switch kind {
	case discoverydomain.EntityKindParcel:
		return "parcels"
	case discoverydomain.EntityKindPlanningRegion:
		return "qh_regions"
	case discoverydomain.EntityKindPlanningProject:
		return "qh_planning_projects"
	case discoverydomain.EntityKindAdministrativeUnit:
		return "province_v2/ward_v2"
	case discoverydomain.EntityKindPOI:
		return "pois"
	default:
		return ""
	}
}

func DetailURL(kind discoverydomain.EntityKind, row Row) string {
	id := strings.TrimSpace(row.ID)
	if id == "" {
		return ""
	}
	switch kind {
	case discoverydomain.EntityKindParcel:
		return "/v2/tqd/parcels/" + id + "/detail"
	case discoverydomain.EntityKindPlanningRegion:
		return "/vung-quy-hoach/" + id
	case discoverydomain.EntityKindPlanningProject:
		return "/do-an-quy-hoach/" + id
	case discoverydomain.EntityKindAdministrativeUnit:
		return administrativeDetailURL(row)
	case discoverydomain.EntityKindPOI:
		return "/v2/tqd/pois/" + id
	default:
		return ""
	}
}

func IsPublicCanonicalPath(path string) bool {
	return strings.HasPrefix(path, "/dia-ban/") ||
		strings.HasPrefix(path, "/do-an-quy-hoach/") ||
		strings.HasPrefix(path, "/vung-quy-hoach/")
}

var nonSlugCharacters = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	replacer := strings.NewReplacer(
		"à", "a", "á", "a", "ạ", "a", "ả", "a", "ã", "a",
		"â", "a", "ầ", "a", "ấ", "a", "ậ", "a", "ẩ", "a", "ẫ", "a",
		"ă", "a", "ằ", "a", "ắ", "a", "ặ", "a", "ẳ", "a", "ẵ", "a",
		"è", "e", "é", "e", "ẹ", "e", "ẻ", "e", "ẽ", "e",
		"ê", "e", "ề", "e", "ế", "e", "ệ", "e", "ể", "e", "ễ", "e",
		"ì", "i", "í", "i", "ị", "i", "ỉ", "i", "ĩ", "i",
		"ò", "o", "ó", "o", "ọ", "o", "ỏ", "o", "õ", "o",
		"ô", "o", "ồ", "o", "ố", "o", "ộ", "o", "ổ", "o", "ỗ", "o",
		"ơ", "o", "ờ", "o", "ớ", "o", "ợ", "o", "ở", "o", "ỡ", "o",
		"ù", "u", "ú", "u", "ụ", "u", "ủ", "u", "ũ", "u",
		"ư", "u", "ừ", "u", "ứ", "u", "ự", "u", "ử", "u", "ữ", "u",
		"ỳ", "y", "ý", "y", "ỵ", "y", "ỷ", "y", "ỹ", "y",
		"đ", "d",
	)
	value = replacer.Replace(value)
	return strings.Trim(nonSlugCharacters.ReplaceAllString(value, "-"), "-")
}

func administrativeDetailURL(row Row) string {
	code := strings.TrimSpace(row.ProvinceCode)
	if strings.HasPrefix(strings.TrimSpace(row.ID), "ward:") {
		code = strings.TrimSpace(row.WardCode)
	}
	nameSlug := slugify(row.Title)
	codeSlug := slugify(code)
	if nameSlug == "" || codeSlug == "" {
		return ""
	}
	return "/dia-ban/" + nameSlug + "-" + codeSlug
}
