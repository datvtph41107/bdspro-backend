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
