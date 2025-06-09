package dto

type UserProfile struct {
	ID       uint64 `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Gender   uint32 `json:"gender"`
}
