package enums

type LocationUnitType uint32
type LocationDataSource uint32

const (
	LocationUnitTypeProvince    LocationUnitType = 10 // Tỉnh
	LocationUnitTypeCity        LocationUnitType = 20 // Thành phố
	LocationUnitTypeSpecialZone LocationUnitType = 30 // Đặc khu
)

const (
	LocationDataSourceInternal LocationDataSource = 10 // Nội bộ
	LocationDataSourceOfficial LocationDataSource = 20 // Chính thức
	LocationDataSourceGIS      LocationDataSource = 30 // GIS
	LocationDataSourceImport   LocationDataSource = 40 // Import
)

var LocationUnitTypeMap = map[LocationUnitType]string{
	LocationUnitTypeProvince:    "Tỉnh",
	LocationUnitTypeCity:        "Phường",
	LocationUnitTypeSpecialZone: "Đặc khu",
}

var LocationDataSourceMap = map[LocationDataSource]string{
	LocationDataSourceInternal: "Nội bộ",
	LocationDataSourceOfficial: "Chính thức",
	LocationDataSourceGIS:      "GIS",
	LocationDataSourceImport:   "Chính phủ",
}

func (e LocationUnitType) String() string {
	if name, ok := LocationUnitTypeMap[e]; ok {
		return name
	}
	return ""
}

func (e LocationUnitType) IsValid() bool {
	_, ok := LocationUnitTypeMap[e]
	return ok
}

func (e LocationDataSource) String() string {
	if name, ok := LocationDataSourceMap[e]; ok {
		return name
	}
	return ""
}

func (e LocationDataSource) IsValid() bool {
	_, ok := LocationDataSourceMap[e]
	return ok
}
