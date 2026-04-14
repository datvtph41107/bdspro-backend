package config

const (
	// Giới hạn hình học
	MaxPolygonAreaKm2  = 100.0 // diện tích tối đa (km²)
	MaxPolygonRadiusKm = 50.0  // bán kính tối đa từ tâm (km)
	MaxPolygonVertices = 200   // số đỉnh tối đa
	MinPolygonVertices = 3     // số đỉnh tối thiểu

	// Giới hạn truy vấn và hiệu năng
	MaxParcelReturn = 5000 // tổng số parcel tối đa cho phép trả về (dù phân trang)
	DefaultPageSize = 50
	MaxPageSize     = 200

	// Cache
	CacheTTLSeconds = 300 // thời gian cache polygon (giây)
)
