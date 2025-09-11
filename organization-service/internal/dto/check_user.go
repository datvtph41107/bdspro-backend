package dto

type CheckUserInOrganization struct {
	UserId   uint64 `json:"userId"`
	IsExist  bool   `json:"isExist"`
	JoinedAt string `json:"joinedAt"`
}
