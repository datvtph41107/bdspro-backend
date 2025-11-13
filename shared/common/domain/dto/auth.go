package _dto

type OTPV3DTO struct {
	Otp      string `json:"otp"`
	Phone    string `json:"phone"`
	Fullname string `json:"fullname"`
	Mode     string `json:"mode"`
	AuthId   uint64 `json:"authId"`
}
