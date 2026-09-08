package enums

type OwnerOf int32

const (
	OwnerOfUser         OwnerOf = 10
	OwnerOfGroup        OwnerOf = 20
	OwnerOfOrganization OwnerOf = 30
)

var OwnerTypeNames = map[OwnerOf]string{
	OwnerOfUser:         "user",
	OwnerOfGroup:        "group",
	OwnerOfOrganization: "organization",
}
