package repo

import (
	cmodels "common/models"
	"map/models"

	"gorm.io/gorm"
)

type MapRepo struct {
	cmodels.BaseRepository
}

func NewMapRepo(db *gorm.DB) *MapRepo {
	return &MapRepo{
		BaseRepository: cmodels.BaseRepository{DB: db},
	}
}

// Find locations near a point
func (r *MapRepo) FindLocationNearby(lat, lng, radius float64) ([]models.MapPoint, error) {
	var points []models.MapPoint
	query := `
        SELECT id, ST_Y(geom::geometry) as lat, ST_X(geom::geometry) as lng, name,
        ST_DistanceSphere(ST_MakePoint(lng, lat), ST_MakePoint(?, ?)) AS distance
        FROM map_points 
        WHERE ST_DistanceSphere(ST_MakePoint(lng, lat), ST_MakePoint(?, ?)) <= ?
        ORDER BY distance;
    `
	result := r.DB.Raw(query, lng, lat, lng, lat, radius).Scan(&points)
	return points, result.Error
}

// Find locations inside a polygon
func (r *MapRepo) FindLocationInPolygon(polygon string) ([]models.MapPointDTO, error) {
	var points []models.MapPointDTO
	query := `
        SELECT ST_Y(geom::geometry) as lat, 
               ST_X(geom::geometry) as lng, 
               name 
        FROM get_partitioned_map_points(ST_GeogFromText(?))
        WHERE ST_Within(geom::geometry, ST_SetSRID(?::geometry, 4326))
        LIMIT 200;
    `
	result := r.DB.Raw(query, polygon, polygon).Scan(&points)
	return points, result.Error
}

// Create new map point
func (r *MapRepo) CreateNewPoint(lat float64, lng float64, name string) error {
	query := `SELECT create_point(CAST(? AS NUMERIC), CAST(? AS NUMERIC), CAST(? AS VARCHAR));`
	result := r.DB.Exec(query, lat, lng, name)
	return result.Error
}
