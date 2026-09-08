package config

const (
	MaxPolygonAreaKm2  = 5.0
	MaxPolygonVertices = 200
	MinPolygonVertices = 3
	DefaultPageSize    = 50
	MaxPageSize        = 200
	MaxParcelReturn    = 5000

	// Polygon draw ở zoom >= ngưỡng này mới trả parcel.
	// Zoom 14 trở xuống trả region để khớp cấp hiển thị của bản đồ.
	ParcelPolygonMinZ = 15
)
