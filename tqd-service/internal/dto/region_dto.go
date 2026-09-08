package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
)

type CreateRegionRequest struct {
	Name        string                 `json:"name" binding:"required"`
	DisplayName string                 `json:"display_name" binding:"required"`
	Description string                 `json:"description"`
	Type        uint32                 `json:"type" binding:"required,min=1,max=4"`
	Status      uint32                 `json:"status" binding:"required,min=10,max=30"`
	LayerID     *uint64                `json:"layer_id"`
	Properties  map[string]interface{} `json:"properties"`

	// Geometry inputs
	GeoJSON   string        `json:"geo_json"`
	WKT       string        `json:"wkt"`
	Lat       float64       `json:"lat"`
	Lng       float64       `json:"lng"`
	Points    [][]float64   `json:"points"`
	Rings     [][][]float64 `json:"rings"`
	MinX      float64       `json:"min_x"`
	MinY      float64       `json:"min_y"`
	MaxX      float64       `json:"max_x"`
	MaxY      float64       `json:"max_y"`
	CenterLat float64       `json:"center_lat"`
	CenterLng float64       `json:"center_lng"`
	Radius    float64       `json:"radius"`
	SourceID  uint64        `json:"source_id"`
	Distance  float64       `json:"distance"`
	Segments  int           `json:"segments"`
}

type UpdateRegionRequest struct {
	Name        *string                `json:"name"`
	DisplayName *string                `json:"display_name"`
	Description *string                `json:"description"`
	Type        *uint32                `json:"type"`
	Status      *uint32                `json:"status"`
	LayerID     uint64                 `json:"layer_id"`
	Properties  map[string]interface{} `json:"properties"`
	Geometry    *string                `json:"geometry"`
}

type ListRegionRequest struct {
	Type    *uint32 `form:"type"`
	Status  *uint32 `form:"status"`
	LayerID *uint64 `form:"layer_id"`
	Search  *string `form:"search"`
	Page    int     `form:"page"`
	Limit   int     `form:"limit"`
}

type PointQueryRequest struct {
	Lat     float64  `form:"lat" binding:"required"`
	Lng     float64  `form:"lng" binding:"required"`
	Type    *uint32  `form:"type"`
	LayerID *uint64  `form:"layer_id"`
	Buffer  *float64 `form:"buffer"`
}

type BBoxQueryRequest struct {
	MinX    float64 `form:"min_x" binding:"required"`
	MinY    float64 `form:"min_y" binding:"required"`
	MaxX    float64 `form:"max_x" binding:"required"`
	MaxY    float64 `form:"max_y" binding:"required"`
	Type    *uint32 `form:"type"`
	LayerID *uint64 `form:"layer_id"`
}

type IntersectionQueryRequest struct {
	ID     uint64  `uri:"id" binding:"required"`
	Buffer float64 `form:"buffer"`
}

type ImportGeoJSONRequest struct {
	Name        string  `json:"name" binding:"required"`
	DisplayName string  `json:"display_name" binding:"required"`
	Description string  `json:"description"`
	Type        uint32  `json:"type" binding:"required,min=1,max=4"`
	Status      uint32  `json:"status" binding:"required,min=10,max=30"`
	LayerID     *uint64 `json:"layer_id"`
	GeoJSON     string  `json:"geo_json" binding:"required"`
}

type RegionResponse struct {
	ID          uint64                 `json:"id"`
	Name        string                 `json:"name"`
	DisplayName string                 `json:"display_name"`
	Description string                 `json:"description"`
	Type        uint32                 `json:"type"`
	TypeName    string                 `json:"type_name"`
	Status      uint32                 `json:"status"`
	StatusName  string                 `json:"status_name"`
	Geometry    string                 `json:"geometry"`
	Centroid    string                 `json:"centroid"`
	BBox        string                 `json:"bbox"`
	Area        float64                `json:"area"`
	Length      float64                `json:"length"`
	LayerID     uint64                 `json:"layer_id"`
	Properties  map[string]interface{} `json:"properties"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

type ListRegionResponse struct {
	Data  []RegionResponse `json:"data"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}

type ImportGeoJSONResponse struct {
	Region       RegionResponse `json:"region"`
	FeatureCount int            `json:"feature_count"`
	Message      string         `json:"message"`
}

// CreateRegionBatchItem — một phần tử trong batch (geometry + layerId bắt buộc)
type CreateRegionBatchItem struct {
	Geometry    []byte                 `json:"geometry" binding:"required"`
	LayerID     uint64                 `json:"layerId" binding:"required"`
	Name        string                 `json:"name"`
	DisplayName string                 `json:"displayName"`
	Type        uint32                 `json:"type"`   // mặc định 1 (Road) nếu bỏ trống
	Status      uint32                 `json:"status"` // mặc định 10 (Active) nếu bỏ trống
	Properties  map[string]interface{} `json:"properties"`
}

// CreateRegionBatchRequest — body request cho batch create
type CreateRegionBatchRequest struct {
	Items []CreateRegionBatchItem `json:"items" binding:"required,min=1"`
}

// CreateRegionBatchResultItem — kết quả của từng phần tử trong batch
type CreateRegionBatchResultItem struct {
	Index   int             `json:"index"`
	Success bool            `json:"success"`
	Region  *RegionResponse `json:"region,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// CreateRegionBatchResponse — response tổng hợp batch create
type CreateRegionBatchResponse struct {
	Items        []CreateRegionBatchResultItem `json:"items"`
	Total        int                           `json:"total"`
	SuccessCount int                           `json:"successCount"`
	FailureCount int                           `json:"failureCount"`
}

type GeometryParser struct{}

func NewGeometryParser() *GeometryParser {
	return &GeometryParser{}
}

func (p *GeometryParser) Parse(req *CreateRegionRequest) (string, error) {
	if req.GeoJSON != "" {
		return req.GeoJSON, nil
	}
	if req.WKT != "" {
		return p.wktToGeoJSON(req.WKT)
	}
	if len(req.Rings) > 0 {
		return p.ringsToGeoJSON(req.Rings)
	}
	if len(req.Points) > 0 {
		return p.pointsToGeoJSON(req.Points)
	}
	if req.MinX != 0 && req.MinY != 0 && req.MaxX != 0 && req.MaxY != 0 {
		return p.bboxToGeoJSON(req.MinX, req.MinY, req.MaxX, req.MaxY)
	}
	if req.CenterLat != 0 && req.CenterLng != 0 && req.Radius > 0 {
		segments := req.Segments
		if segments <= 0 {
			segments = 36
		}
		return p.circleToGeoJSON(req.CenterLat, req.CenterLng, req.Radius, segments)
	}
	if req.SourceID > 0 && req.Distance > 0 {
		return "", errors.New("buffer requires source region to be loaded first")
	}
	if req.Lat != 0 && req.Lng != 0 {
		return p.pointToGeoJSON(req.Lat, req.Lng)
	}
	return "", errors.New("no valid geometry provided")
}

// func (p *GeometryParser) wktToGeoJSON(wkt string) (string, error) {
// 	geom := map[string]interface{}{"type": "Point", "coordinates": []float64{0, 0}}
// 	bytes, err := json.Marshal(geom)
// 	return string(bytes), err
// }

// func (p *GeometryParser) ringsToGeoJSON(rings [][][]float64) (string, error) {
// 	if len(rings) == 0 {
// 		return "", errors.New("rings cannot be empty")
// 	}
// 	var geom interface{}
// 	if len(rings) == 1 {
// 		geom = map[string]interface{}{"type": "Polygon", "coordinates": rings[0]}
// 	} else {
// 		geom = map[string]interface{}{"type": "MultiPolygon", "coordinates": rings}
// 	}
// 	bytes, err := json.Marshal(geom)
// 	return string(bytes), err
// }

// func (p *GeometryParser) pointsToGeoJSON(points [][]float64) (string, error) {
// 	if len(points) < 2 {
// 		return "", errors.New("points need at least 2 points")
// 	}
// 	geom := map[string]interface{}{"type": "LineString", "coordinates": points}
// 	bytes, err := json.Marshal(geom)
// 	return string(bytes), err
// }

func (p *GeometryParser) bboxToGeoJSON(minX, minY, maxX, maxY float64) (string, error) {
	// Cách 1: Dùng cấu trúc đúng
	coordinates := [][][]float64{
		{
			{minX, minY},
			{maxX, minY},
			{maxX, maxY},
			{minX, maxY},
			{minX, minY},
		},
	}

	geom := map[string]interface{}{
		"type":        "Polygon",
		"coordinates": coordinates,
	}

	bytes, err := json.Marshal(geom)
	return string(bytes), err
}

func (p *GeometryParser) circleToGeoJSON(centerLat, centerLng, radiusMeters float64, segments int) (string, error) {
	radiusDeg := radiusMeters / 111319.9
	points := make([][]float64, segments+1)

	for i := 0; i <= segments; i++ {
		angle := 2 * math.Pi * float64(i) / float64(segments)
		x := centerLng + radiusDeg*math.Cos(angle)
		y := centerLat + radiusDeg*math.Sin(angle)
		points[i] = []float64{x, y}
	}

	// Đóng vòng tròn (thêm điểm đầu tiên vào cuối)
	points[segments] = points[0]

	coordinates := [][][]float64{points}

	geom := map[string]interface{}{
		"type":        "Polygon",
		"coordinates": coordinates,
	}

	bytes, err := json.Marshal(geom)
	return string(bytes), err
}

func (p *GeometryParser) pointsToGeoJSON(points [][]float64) (string, error) {
	if len(points) < 2 {
		return "", errors.New("points need at least 2 points for LineString")
	}

	geom := map[string]interface{}{
		"type":        "LineString",
		"coordinates": points,
	}

	bytes, err := json.Marshal(geom)
	return string(bytes), err
}

func (p *GeometryParser) ringsToGeoJSON(rings [][][]float64) (string, error) {
	if len(rings) == 0 {
		return "", errors.New("rings cannot be empty")
	}

	var geom interface{}
	if len(rings) == 1 {
		geom = map[string]interface{}{
			"type":        "Polygon",
			"coordinates": rings[0],
		}
	} else {
		geom = map[string]interface{}{
			"type":        "MultiPolygon",
			"coordinates": rings,
		}
	}

	bytes, err := json.Marshal(geom)
	return string(bytes), err
}

func (p *GeometryParser) pointToGeoJSON(lat, lng float64) (string, error) {
	geom := map[string]interface{}{
		"type":        "Point",
		"coordinates": []float64{lng, lat},
	}

	bytes, err := json.Marshal(geom)
	return string(bytes), err
}

func (p *GeometryParser) wktToGeoJSON(wkt string) (string, error) {
	// Parse WKT to GeoJSON
	// Đây là phiên bản đơn giản, bạn có thể dùng thư viện để parse WKT chính xác hơn
	wkt = strings.TrimSpace(wkt)

	if strings.HasPrefix(wkt, "POINT") {
		var x, y float64
		fmt.Sscanf(wkt, "POINT(%f %f)", &x, &y)
		return p.pointToGeoJSON(y, x) // lat=y, lng=x
	}

	if strings.HasPrefix(wkt, "LINESTRING") {
		// Parse LINESTRING(x1 y1, x2 y2, ...)
		coordsPart := strings.TrimPrefix(wkt, "LINESTRING(")
		coordsPart = strings.TrimSuffix(coordsPart, ")")
		pointsStr := strings.Split(coordsPart, ",")

		points := make([][]float64, 0, len(pointsStr))
		for _, pointStr := range pointsStr {
			var x, y float64
			fmt.Sscanf(strings.TrimSpace(pointStr), "%f %f", &x, &y)
			points = append(points, []float64{x, y})
		}
		return p.pointsToGeoJSON(points)
	}

	if strings.HasPrefix(wkt, "POLYGON") {
		// Parse POLYGON((x1 y1, x2 y2, ...))
		coordsPart := strings.TrimPrefix(wkt, "POLYGON((")
		coordsPart = strings.TrimSuffix(coordsPart, "))")
		pointsStr := strings.Split(coordsPart, ",")

		rings := make([][]float64, 0, len(pointsStr))
		for _, pointStr := range pointsStr {
			var x, y float64
			fmt.Sscanf(strings.TrimSpace(pointStr), "%f %f", &x, &y)
			rings = append(rings, []float64{x, y})
		}

		// Đóng polygon
		if len(rings) > 0 {
			rings = append(rings, rings[0])
		}

		coordinates := [][][]float64{rings}
		geom := map[string]interface{}{
			"type":        "Polygon",
			"coordinates": coordinates,
		}
		bytes, err := json.Marshal(geom)
		return string(bytes), err
	}

	return "", fmt.Errorf("unsupported WKT type: %s", wkt)
}
