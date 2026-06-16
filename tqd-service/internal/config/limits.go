package config

const (
	MaxPolygonAreaKm2  = 100.0
	MaxPolygonVertices = 200
	MinPolygonVertices = 3
	DefaultPageSize    = 50
	MaxPageSize        = 200
	MaxParcelReturn    = 5000

	// Polygon draw ở zoom >= ngưỡng này mới trả parcel.
	// Zoom thấp hơn trả region để khớp cấp hiển thị của bản đồ.
	ParcelPolygonMinZ = 14
)