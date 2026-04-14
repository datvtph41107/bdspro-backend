package geoutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
)

// ValidateGeoJSON kiểm tra tính hợp lệ của GeoJSON polygon:
// - đúng cấu trúc
// - số đỉnh tối thiểu
// - không tự cắt
func ValidateGeoJSON(geoJSON string) error {
	var poly orb.Polygon
	if err := json.Unmarshal([]byte(geoJSON), &poly); err != nil {
		return fmt.Errorf("invalid GeoJSON: %v", err)
	}
	if len(poly) == 0 {
		return errors.New("polygon must have at least one ring")
	}
	outerRing := poly[0]
	if len(outerRing) < 4 { // cần ít nhất 4 điểm (3 + lặp lại điểm đầu)
		return errors.New("outer ring must have at least 3 distinct vertices")
	}

	// Kiểm tra tự cắt (self-intersection)
	if isSelfIntersecting(poly) {
		return errors.New("polygon self-intersects")
	}

	return nil
}

// isSelfIntersecting kiểm tra xem một polygon có tự cắt hay không.
// Hàm này chỉ kiểm tra vòng ngoài (outer ring).
func isSelfIntersecting(poly orb.Polygon) bool {
	ring := poly[0]
	n := len(ring)

	// Kiểm tra tất cả các cặp cạnh không kề nhau
	for i := 0; i < n-1; i++ {
		for j := i + 2; j < n-1; j++ {
			// Bỏ qua các cạnh kề nhau
			if i == 0 && j == n-2 {
				continue
			}
			if intersect(ring[i], ring[i+1], ring[j], ring[j+1]) {
				return true
			}
		}
	}
	return false
}

// intersect kiểm tra hai đoạn thẳng (a,b) và (c,d) có giao nhau tại một điểm nằm giữa hay không.
func intersect(a, b, c, d orb.Point) bool {
	// Hàm tính hướng của ba điểm
	orientation := func(p, q, r orb.Point) int {
		val := (q[1]-p[1])*(r[0]-q[0]) - (q[0]-p[0])*(r[1]-q[1])
		if val == 0 {
			return 0 // thẳng hàng
		}
		if val > 0 {
			return 1 // chiều kim đồng hồ
		}
		return 2 // ngược chiều kim đồng hồ
	}

	o1 := orientation(a, b, c)
	o2 := orientation(a, b, d)
	o3 := orientation(c, d, a)
	o4 := orientation(c, d, b)

	// Trường hợp tổng quát
	if o1 != o2 && o3 != o4 {
		return true
	}

	// Trường hợp đặc biệt: thẳng hàng và chồng lấn
	if o1 == 0 && onSegment(a, c, b) {
		return true
	}
	if o2 == 0 && onSegment(a, d, b) {
		return true
	}
	if o3 == 0 && onSegment(c, a, d) {
		return true
	}
	if o4 == 0 && onSegment(c, b, d) {
		return true
	}

	return false
}

// onSegment kiểm tra điểm q có nằm trên đoạn thẳng pr không.
func onSegment(p, q, r orb.Point) bool {
	return q[0] <= max(p[0], r[0]) && q[0] >= min(p[0], r[0]) &&
		q[1] <= max(p[1], r[1]) && q[1] >= min(p[1], r[1])
}

// CalcAreaKm2 tính diện tích polygon (km²) trên mặt cầu trái đất (haversine-based).
// Input: GeoJSON Polygon string.
// Output: diện tích km².
func CalcAreaKm2(geoJSON string) (float64, error) {
	var poly orb.Polygon
	if err := json.Unmarshal([]byte(geoJSON), &poly); err != nil {
		return 0, err
	}
	if len(poly) == 0 {
		return 0, errors.New("empty polygon")
	}
	// Tính diện tích outer ring trừ đi diện tích các inner rings (holes)
	area := ringAreaKm2(poly[0])
	for i := 1; i < len(poly); i++ {
		area -= ringAreaKm2(poly[i])
	}
	if area < 0 {
		area = -area
	}
	return area, nil
}

// ringAreaKm2 tính diện tích một ring (đóng) theo công thức haversine.
func ringAreaKm2(ring orb.Ring) float64 {
	if len(ring) < 3 {
		return 0
	}
	const earthRadiusKm = 6371.0
	rad := func(deg float64) float64 { return deg * math.Pi / 180.0 }

	area := 0.0
	n := len(ring)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		lat1 := rad(ring[i].Y())
		lng1 := rad(ring[i].X())
		lat2 := rad(ring[j].Y())
		lng2 := rad(ring[j].X())
		area += (lng2 - lng1) * math.Sin((lat1+lat2)/2)
	}
	area = area * earthRadiusKm * earthRadiusKm
	if area < 0 {
		area = -area
	}
	return area
}

// ScalePolygon scale polygon quanh centroid với factor.
// Input: GeoJSON Polygon string, factor > 0.
// Output: GeoJSON Polygon đã scale.
func ScalePolygon(geoJSON string, factor float64) (string, error) {
	var poly orb.Polygon
	if err := json.Unmarshal([]byte(geoJSON), &poly); err != nil {
		return "", err
	}
	if factor <= 0 {
		return "", errors.New("scale factor must be positive")
	}
	if len(poly) == 0 {
		return "", errors.New("empty polygon")
	}

	// Sử dụng planar.CentroidArea để lấy centroid chính xác
	centroid, _ := planar.CentroidArea(poly)

	newPoly := make(orb.Polygon, len(poly))
	for r, ring := range poly {
		newRing := make(orb.Ring, len(ring))
		for i, p := range ring {
			dx := p.X() - centroid.X()
			dy := p.Y() - centroid.Y()
			newRing[i] = orb.Point{centroid.X() + dx*factor, centroid.Y() + dy*factor}
		}
		newPoly[r] = newRing
	}

	// Đảm bảo ring cuối trùng với điểm đầu (khép kín)
	for r := range newPoly {
		if len(newPoly[r]) > 1 {
			newPoly[r][len(newPoly[r])-1] = newPoly[r][0]
		}
	}

	b, err := json.Marshal(newPoly)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CalcRadiusKm tính bán kính tương đương (khoảng cách từ centroid đến đỉnh xa nhất) (km).
func CalcRadiusKm(geoJSON string) (float64, error) {
	var poly orb.Polygon
	if err := json.Unmarshal([]byte(geoJSON), &poly); err != nil {
		return 0, err
	}
	if len(poly) == 0 {
		return 0, errors.New("empty polygon")
	}
	centroid, _ := planar.CentroidArea(poly)
	maxDist := 0.0
	for _, p := range poly[0] {
		dist := haversineKm(centroid.Y(), centroid.X(), p.Y(), p.X())
		if dist > maxDist {
			maxDist = dist
		}
	}
	return maxDist, nil
}

// haversineKm tính khoảng cách giữa hai điểm (lat1, lon1) và (lat2, lon2) theo km.
func haversineKm(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0
	rad := math.Pi / 180.0
	phi1 := lat1 * rad
	phi2 := lat2 * rad
	deltaPhi := (lat2 - lat1) * rad
	deltaLambda := (lon2 - lon1) * rad

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

// Hàm min và max trợ giúp cho onSegment
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
