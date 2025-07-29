package dto

type CheckUserInOrganization struct {
	UserId   uint32 `json:"userId"`
	IsExist  bool   `json:"isExist"`
	JoinedAt string `json:"joinedAt"`
}
