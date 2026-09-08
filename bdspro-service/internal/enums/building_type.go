package enums

type EBuildingType uint32

const (
	EBuildingTypeResidential EBuildingType = 10
	EBuildingTypeFactory     EBuildingType = 20
)

var BuildingTypeNames = map[EBuildingType]string{
	EBuildingTypeResidential: "Nhà ở",
	EBuildingTypeFactory:     "Nhà xưởng",
}
