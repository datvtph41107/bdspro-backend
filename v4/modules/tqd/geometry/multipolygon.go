package geometry

import (
	"encoding/json"
	"fmt"
	"strings"
)

// stripGeoJSONExtras removes keys that can break ST_GeomFromGeoJSON / cast on older stacks.
func stripGeoJSONExtras(geom map[string]interface{}) map[string]interface{} {
	if geom == nil {
		return nil
	}
	out := make(map[string]interface{}, len(geom))
	for k, v := range geom {
		switch strings.ToLower(k) {
		case "crs", "bbox":
			continue
		default:
			out[k] = v
		}
	}
	return out
}

// IsPolygonType trả về true nếu geometry type là Polygon hoặc MultiPolygon.
func IsPolygonType(geom map[string]interface{}) bool {
	if geom == nil {
		return false
	}
	t, _ := geom["type"].(string)
	t = strings.TrimSpace(strings.ToLower(t))
	return t == "polygon" || t == "multipolygon"
}

// GetGeomType trả về geometry type string từ GeoJSON object.
func GetGeomType(geom map[string]interface{}) string {
	if geom == nil {
		return ""
	}
	t, _ := geom["type"].(string)
	return strings.TrimSpace(t)
}

// ToRawGeoJSONBytes serialize geometry as-is (chỉ strip CRS/BBox) không convert sang MultiPolygon.
func ToRawGeoJSONBytes(geom map[string]interface{}) ([]byte, error) {
	if geom == nil {
		return nil, fmt.Errorf("nil geometry")
	}
	return json.Marshal(stripGeoJSONExtras(geom))
}

// ToMultiPolygonGeoJSONBytes converts GeoJSON Polygon to MultiPolygon (single shell)
// and passes MultiPolygon through. qh_regions.geometry is geometry(MultiPolygon,4326).
func ToMultiPolygonGeoJSONBytes(geom map[string]interface{}) ([]byte, error) {
	if geom == nil {
		return nil, fmt.Errorf("nil geometry")
	}
	g := stripGeoJSONExtras(geom)

	t, _ := g["type"].(string)
	t = strings.TrimSpace(t)
	coords := g["coordinates"]

	switch t {
	case "Polygon":
		mp := map[string]interface{}{
			"type":        "MultiPolygon",
			"coordinates": []interface{}{coords},
		}
		return json.Marshal(mp)
	case "MultiPolygon":
		return json.Marshal(g)
	default:
		return nil, fmt.Errorf("geometry type %q is not Polygon or MultiPolygon (cannot store in MultiPolygon column)", t)
	}
}
