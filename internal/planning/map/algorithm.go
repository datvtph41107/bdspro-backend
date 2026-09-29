package _map

import "math"

type Point struct {
	ID   uint32
	Type uint8
	Lat  float64
	Lon  float64
}

type Grid struct {
	points []Point
}

func MakeGrid() *Grid {
	return &Grid{
		points: make([]Point, 0),
	}
}

func (g *Grid) Add(p Point) {
	g.points = append(g.points, p)
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {

	const R = 6371000 // bán kính trái đất (m)

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	lat1 = lat1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180

	a :=
		math.Sin(dLat/2)*math.Sin(dLat/2) +
			math.Cos(lat1)*math.Cos(lat2)*
				math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func (g *Grid) Nearby(lat, lon float64, r int32) []Point {

	var result []Point

	for _, p := range g.points {

		d := haversine(lat, lon, p.Lat, p.Lon)

		if d <= float64(r) {
			result = append(result, p)
		}
	}

	return result
}
