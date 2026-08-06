package enums

// TagType phân loại tag dùng chung trong TQD.
type TagType uint32

const (
	TagTypeScope TagType = 10 // Khu vực (area)
)

var TagTypeLabelMap = map[TagType]string{
	TagTypeScope: "Khu vực",
}

func (t TagType) DisplayName() string {
	if label, ok := TagTypeLabelMap[t]; ok {
		return label
	}
	return ""
}

func GetTagTypeLabel(tagType uint32) string {
	return TagType(tagType).DisplayName()
}
