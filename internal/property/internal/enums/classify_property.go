package enums

type EClassifyProperty uint32

const (
	EClassifyPropertyLand  EClassifyProperty = 10
	EClassifyPropertyHouse EClassifyProperty = 20
	EClassifyPropertyCell  EClassifyProperty = 30
	EClassifyPropertyOther EClassifyProperty = 40
)

var EClassifyPropertyNames = map[EClassifyProperty]string{
	EClassifyPropertyLand:  "Đất",
	EClassifyPropertyHouse: "Nhà",
	EClassifyPropertyCell:  "Phòng",
	EClassifyPropertyOther: "Khác",
}

func (e EClassifyProperty) String() string {
	return EClassifyPropertyNames[e]
}

func (e EClassifyProperty) IsValid() bool {
	return e == EClassifyPropertyCell || e == EClassifyPropertyHouse || e == EClassifyPropertyLand || e == EClassifyPropertyOther
}
