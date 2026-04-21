// internal/enums/layer_enums.go
package enums

// LayerStatus trạng thái lớp bản đồ (uint32)
type LayerStatus uint32

const (
	LayerStatusActive   LayerStatus = 10 // Hoạt động
	LayerStatusInactive LayerStatus = 20 // Không hoạt động
	LayerStatusDeleted  LayerStatus = 30 // Đã xóa
)

// String trả về tên status
func (s LayerStatus) String() string {
	switch s {
	case LayerStatusActive:
		return "active"
	case LayerStatusInactive:
		return "inactive"
	case LayerStatusDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}

// IsValid kiểm tra status hợp lệ
func (s LayerStatus) IsValid() bool {
	switch s {
	case LayerStatusActive, LayerStatusInactive, LayerStatusDeleted:
		return true
	default:
		return false
	}
}

// Value trả về giá trị uint32
func (s LayerStatus) Value() uint32 {
	return uint32(s)
}

// LayerType loại dữ liệu bản đồ (uint32)
type LayerType uint32

const (
	LayerTypeVector LayerType = 1 // Vector data
	LayerTypeRaster LayerType = 2 // Raster data
	LayerTypeWMS    LayerType = 3 // WMS service
	LayerTypeWMTS   LayerType = 4 // WMTS service
	LayerTypeTMS    LayerType = 5 // TMS service
	LayerTypeXYZ    LayerType = 6 // XYZ tiles
)

// String trả về tên type
func (t LayerType) String() string {
	switch t {
	case LayerTypeVector:
		return "vector"
	case LayerTypeRaster:
		return "raster"
	case LayerTypeWMS:
		return "wms"
	case LayerTypeWMTS:
		return "wmts"
	case LayerTypeTMS:
		return "tms"
	case LayerTypeXYZ:
		return "xyz"
	default:
		return "unknown"
	}
}

// IsValid kiểm tra type hợp lệ
func (t LayerType) IsValid() bool {
	switch t {
	case LayerTypeVector, LayerTypeRaster, LayerTypeWMS, LayerTypeWMTS, LayerTypeTMS, LayerTypeXYZ:
		return true
	default:
		return false
	}
}

// Value trả về giá trị uint32
func (t LayerType) Value() uint32 {
	return uint32(t)
}

// FromString chuyển từ string sang enum
func (t *LayerType) FromString(str string) bool {
	switch str {
	case "vector":
		*t = LayerTypeVector
	case "raster":
		*t = LayerTypeRaster
	case "wms":
		*t = LayerTypeWMS
	case "wmts":
		*t = LayerTypeWMTS
	case "tms":
		*t = LayerTypeTMS
	case "xyz":
		*t = LayerTypeXYZ
	default:
		return false
	}
	return true
}
