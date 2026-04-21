package enums

type RegionType uint32

const (
	RegionTypeRoad     RegionType = 1
	RegionTypePlanning RegionType = 2
	RegionTypeBoundary RegionType = 3
	RegionTypeZone     RegionType = 4
)

func (r RegionType) String() string {
	switch r {
	case RegionTypeRoad:
		return "road"
	case RegionTypePlanning:
		return "planning"
	case RegionTypeBoundary:
		return "boundary"
	case RegionTypeZone:
		return "zone"
	default:
		return "unknown"
	}
}

func (r RegionType) IsValid() bool {
	return r >= 1 && r <= 4
}

type RegionStatus uint32

const (
	RegionStatusActive   RegionStatus = 10
	RegionStatusInactive RegionStatus = 20
	RegionStatusPlanned  RegionStatus = 30
)

func (r RegionStatus) String() string {
	switch r {
	case RegionStatusActive:
		return "active"
	case RegionStatusInactive:
		return "inactive"
	case RegionStatusPlanned:
		return "planned"
	default:
		return "unknown"
	}
}

func (r RegionStatus) IsValid() bool {
	return r >= 10 && r <= 30
}
