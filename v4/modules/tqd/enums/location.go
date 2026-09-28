package enums

type LocationUnitType uint32
type LocationDataSource uint32

const (
	LocationUnitTypeProvince    LocationUnitType = 10
	LocationUnitTypeCity        LocationUnitType = 20
	LocationUnitTypeSpecialZone LocationUnitType = 30
)

const (
	LocationDataSourceInternal LocationDataSource = 10
	LocationDataSourceOfficial LocationDataSource = 20
	LocationDataSourceGIS      LocationDataSource = 30
	LocationDataSourceImport   LocationDataSource = 40
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

func (e LocationDataSource) String() string {
	if name, ok := LocationDataSourceMap[e]; ok {
		return name
	}
	return ""
}
