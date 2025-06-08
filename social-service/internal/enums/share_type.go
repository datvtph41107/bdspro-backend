package enums

type ShareType int

const (
	ShareTypeNewsFeed     ShareType = 10
	ShareTypeFriend       ShareType = 20
	ShareTypeGroup        ShareType = 30
	ShareTypeOrganization ShareType = 40
)

var ShareTypeMap = map[ShareType]string{
	ShareTypeNewsFeed:     "Bài viết",
	ShareTypeFriend:       "Bạn bè",
	ShareTypeGroup:        "Nhóm",
	ShareTypeOrganization: "Tổ chức",
}
