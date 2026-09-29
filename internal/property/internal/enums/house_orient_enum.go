package enums

type EHouseOrient uint32

const (
	EHouseOrientNone      EHouseOrient = 0  // Chưa xác định
	EHouseOrientEast      EHouseOrient = 10 // Đông
	EHouseOrientWest      EHouseOrient = 20 // Tây
	EHouseOrientSouth     EHouseOrient = 30 // Nam
	EHouseOrientNorth     EHouseOrient = 40 // Bắc
	EHouseOrientEastSouth EHouseOrient = 50 // Đông Nam
	EHouseOrientEastNorth EHouseOrient = 60 // Đông Bắc
	EHouseOrientWestSouth EHouseOrient = 70 // Tây Nam
	EHouseOrientWestNorth EHouseOrient = 80 // Tây Bắc
)

var EHouseOrientNames = map[EHouseOrient]string{
	EHouseOrientNone:      "Chưa xác định",
	EHouseOrientEast:      "Đông",
	EHouseOrientWest:      "Tây",
	EHouseOrientSouth:     "Nam",
	EHouseOrientNorth:     "Bắc",
	EHouseOrientEastSouth: "Đông Nam",
	EHouseOrientEastNorth: "Đông Bắc",
	EHouseOrientWestSouth: "Tây Nam",
	EHouseOrientWestNorth: "Tây Bắc",
}
