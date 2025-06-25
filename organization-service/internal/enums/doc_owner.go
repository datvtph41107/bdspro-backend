package enums

type DocOwner uint

const (
	DocOwnerDeal  DocOwner = 10
	DocOwnerGroup DocOwner = 20
)

var DocOwnerMap = map[DocOwner]string{
	DocOwnerDeal:  "Thương vụ",
	DocOwnerGroup: "Nhóm",
}
