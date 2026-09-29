package _enum

type SubjectTypeProperty uint32

const (
	SUBJECT_TYPE_DISTRIBUTE SubjectTypeProperty = 1
	SUBJECT_TYPE_PRODUCT    SubjectTypeProperty = 2
	SUBJECT_TYPE_PROPERTY   SubjectTypeProperty = 3
)

var SubjectTypePropertyLabel = map[SubjectTypeProperty]string{
	SUBJECT_TYPE_DISTRIBUTE: "Phân phối",
	SUBJECT_TYPE_PRODUCT:    "Sản phẩm",
	SUBJECT_TYPE_PROPERTY:   "Bất động sản",
}

func (t SubjectTypeProperty) String() string {
	if v, ok := SubjectTypePropertyLabel[t]; ok {
		return v
	}
	return "Unknown"
}
