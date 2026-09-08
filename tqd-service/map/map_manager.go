package _map

import (
	"context"
	"math/rand"
	"time"
)

type MapManager struct {
	grids *Grid
}

func NewMapManager() *MapManager {
	grid := MakeGrid()

	rand.Seed(time.Now().UnixNano())

	const N = 1_000_000

	// bounding box VN
	minLat := 8.18
	maxLat := 23.39

	minLon := 102.14
	maxLon := 109.46

	for i := 0; i < N; i++ {

		lat := minLat + rand.Float64()*(maxLat-minLat)
		lon := minLon + rand.Float64()*(maxLon-minLon)

		p := Point{
			ID:   uint32(i),
			Type: uint8(rand.Intn(5)),
			Lat:  lat,
			Lon:  lon,
		}

		grid.Add(p)
	}
	return &MapManager{
		grids: grid,
	}
}

func (mm *MapManager) Nearby(ctx context.Context, lat, lon float64, r int32) ([]Point, error) {
	return mm.grids.Nearby(lat, lon, r), nil
}
