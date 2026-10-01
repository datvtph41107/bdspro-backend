package geometry

import (
	"encoding/json"
	"math"
)

type GeometryInfo struct {
	GeoJSON   string
	Centroid  Point
	Area      float64
	Length    float64
	BBox      BoundingBox
	IsPoint   bool
	IsLine    bool
	IsPolygon bool
}

type Point struct {
	Lat float64
	Lng float64
}

type BoundingBox struct {
	MinX float64 `json:"min_x"`
	MinY float64 `json:"min_y"`
	MaxX float64 `json:"max_x"`
	MaxY float64 `json:"max_y"`
}

func CalculateBBox(coords interface{}) BoundingBox {
	bbox := BoundingBox{
		MinX: math.Inf(1),
		MinY: math.Inf(1),
		MaxX: math.Inf(-1),
		MaxY: math.Inf(-1),
	}

	flattenCoords(coords, func(x, y float64) {
		if x < bbox.MinX {
			bbox.MinX = x
		}
		if y < bbox.MinY {
			bbox.MinY = y
		}
		if x > bbox.MaxX {
			bbox.MaxX = x
		}
		if y > bbox.MaxY {
			bbox.MaxY = y
		}
	})

	if bbox.MinX == math.Inf(1) {
		bbox.MinX = 0
		bbox.MinY = 0
		bbox.MaxX = 0
		bbox.MaxY = 0
	}

	return bbox
}

func ParseGeometry(geoJSON string) (*GeometryInfo, error) {
	var geom map[string]interface{}
	if err := json.Unmarshal([]byte(geoJSON), &geom); err != nil {
		return nil, err
	}

	info := &GeometryInfo{GeoJSON: geoJSON}
	geomType, _ := geom["type"].(string)
	coordinates, _ := geom["coordinates"].(interface{})

	switch geomType {
	case "Point":
		info.IsPoint = true
		coords, ok := coordinates.([]interface{})
		if ok && len(coords) >= 2 {
			info.Centroid = Point{Lng: coords[0].(float64), Lat: coords[1].(float64)}
			info.BBox = CalculateBBox(coordinates)
		}
	case "LineString":
		info.IsLine = true
		info.Length = calculateLength(coordinates)
		info.Centroid = calculateLineCentroid(coordinates)
		info.BBox = CalculateBBox(coordinates)
	case "Polygon":
		info.IsPolygon = true
		info.Area = calculatePolygonArea(coordinates)
		info.Length = calculatePolygonPerimeter(coordinates)
		info.Centroid = calculatePolygonCentroid(coordinates)
		info.BBox = CalculateBBox(coordinates)
	case "MultiPolygon":
		info.IsPolygon = true
		info.Area = calculateMultiPolygonArea(coordinates)
		info.Centroid = calculateMultiPolygonCentroid(coordinates)
		info.BBox = CalculateBBox(coordinates)
	}

	return info, nil
}

func PointToGeoJSON(lat, lng float64) string {
	point := map[string]interface{}{
		"type":        "Point",
		"coordinates": []float64{lng, lat},
	}
	bytes, _ := json.Marshal(point)
	return string(bytes)
}

func CheckPointInGeometry(geometryGeoJSON string, lat, lng float64) (bool, error) {
	var geom map[string]interface{}
	if err := json.Unmarshal([]byte(geometryGeoJSON), &geom); err != nil {
		return false, err
	}

	geomType, _ := geom["type"].(string)
	coordinates, _ := geom["coordinates"].(interface{})

	switch geomType {
	case "Polygon":
		return pointInPolygon(lat, lng, coordinates), nil
	case "MultiPolygon":
		polygons, ok := coordinates.([]interface{})
		if !ok {
			return false, nil
		}
		for _, poly := range polygons {
			if pointInPolygon(lat, lng, poly) {
				return true, nil
			}
		}
	}
	return false, nil
}

// Helper functions
func flattenCoords(coords interface{}, fn func(x, y float64)) {
	switch v := coords.(type) {
	case []interface{}:
		for _, item := range v {
			flattenCoords(item, fn)
		}
	case []float64:
		if len(v) >= 2 {
			fn(v[0], v[1])
		}
	}
}

func calculateLength(coords interface{}) float64 {
	points, ok := coords.([]interface{})
	if !ok || len(points) < 2 {
		return 0
	}
	length := 0.0
	for i := 1; i < len(points); i++ {
		p1, _ := points[i-1].([]interface{})
		p2, _ := points[i].([]interface{})
		if len(p1) >= 2 && len(p2) >= 2 {
			dx := p2[0].(float64) - p1[0].(float64)
			dy := p2[1].(float64) - p1[1].(float64)
			length += math.Sqrt(dx*dx+dy*dy) * 111319.9
		}
	}
	return length
}

func calculatePolygonArea(coords interface{}) float64 {
	rings, ok := coords.([]interface{})
	if !ok || len(rings) == 0 {
		return 0
	}
	points, ok := rings[0].([]interface{})
	if !ok {
		return 0
	}
	area := 0.0
	n := len(points)
	for i := 0; i < n; i++ {
		p1, _ := points[i].([]interface{})
		p2, _ := points[(i+1)%n].([]interface{})
		if len(p1) >= 2 && len(p2) >= 2 {
			area += p1[0].(float64) * p2[1].(float64)
			area -= p2[0].(float64) * p1[1].(float64)
		}
	}
	area = math.Abs(area) / 2
	return area * 111319.9 * 111319.9
}

func calculatePolygonPerimeter(coords interface{}) float64 {
	rings, ok := coords.([]interface{})
	if !ok || len(rings) == 0 {
		return 0
	}
	return calculateLength(rings[0])
}

func calculateLineCentroid(coords interface{}) Point {
	points, ok := coords.([]interface{})
	if !ok || len(points) == 0 {
		return Point{}
	}
	var sumLat, sumLng float64
	for _, p := range points {
		pt, _ := p.([]interface{})
		if len(pt) >= 2 {
			sumLng += pt[0].(float64)
			sumLat += pt[1].(float64)
		}
	}
	return Point{Lat: sumLat / float64(len(points)), Lng: sumLng / float64(len(points))}
}

func calculatePolygonCentroid(coords interface{}) Point {
	rings, ok := coords.([]interface{})
	if !ok || len(rings) == 0 {
		return Point{}
	}
	points, ok := rings[0].([]interface{})
	if !ok || len(points) == 0 {
		return Point{}
	}
	var sumLat, sumLng float64
	for _, p := range points {
		pt, _ := p.([]interface{})
		if len(pt) >= 2 {
			sumLng += pt[0].(float64)
			sumLat += pt[1].(float64)
		}
	}
	return Point{Lat: sumLat / float64(len(points)), Lng: sumLng / float64(len(points))}
}

func calculateMultiPolygonArea(coords interface{}) float64 {
	polygons, ok := coords.([]interface{})
	if !ok {
		return 0
	}
	totalArea := 0.0
	for _, poly := range polygons {
		totalArea += calculatePolygonArea(poly)
	}
	return totalArea
}

func calculateMultiPolygonCentroid(coords interface{}) Point {
	polygons, ok := coords.([]interface{})
	if !ok || len(polygons) == 0 {
		return Point{}
	}
	var sumLat, sumLng float64
	for _, poly := range polygons {
		centroid := calculatePolygonCentroid(poly)
		sumLat += centroid.Lat
		sumLng += centroid.Lng
	}
	return Point{Lat: sumLat / float64(len(polygons)), Lng: sumLng / float64(len(polygons))}
}

func pointInPolygon(lat, lng float64, polygon interface{}) bool {
	rings, ok := polygon.([]interface{})
	if !ok || len(rings) == 0 {
		return false
	}
	points, ok := rings[0].([]interface{})
	if !ok {
		return false
	}
	intersections := 0
	n := len(points)
	for i := 0; i < n; i++ {
		p1, _ := points[i].([]interface{})
		p2, _ := points[(i+1)%n].([]interface{})
		if len(p1) < 2 || len(p2) < 2 {
			continue
		}
		y1 := p1[1].(float64)
		y2 := p2[1].(float64)
		if (y1 > lat) != (y2 > lat) {
			x1 := p1[0].(float64)
			x2 := p2[0].(float64)
			xIntersect := (x2-x1)*(lat-y1)/(y2-y1) + x1
			if xIntersect < lng {
				intersections++
			}
		}
	}
	return intersections%2 == 1
}
