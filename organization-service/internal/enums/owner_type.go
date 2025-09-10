package enums

type OwnerOf uint

const (
	OwnerOfUser         OwnerOf = 10
	OwnerOfGroup        OwnerOf = 20
	OwnerOfOrganization OwnerOf = 30
)

var OwnerOfMap = map[OwnerOf]string{
	OwnerOfUser:         "Người dùng",
	OwnerOfGroup:        "Nhóm",
	OwnerOfOrganization: "Công ty",
}
