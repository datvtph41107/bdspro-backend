package enums

// TagType định nghĩa các loại tag mà người dùng có thể sử dụng trong hồ sơ.
type TagType uint32

const (
	TagTypeSpecialty TagType = 10
	TagTypeMainArea  TagType = 20
	TagTypeProject   TagType = 30
	TagTypeJob       TagType = 40
)

// IsValid kiểm tra giá trị tag type có hợp lệ không.
func (t TagType) IsValid() bool {
	switch t {
	case TagTypeSpecialty, TagTypeMainArea, TagTypeProject, TagTypeJob:
		return true
	default:
		return false
	}
}
