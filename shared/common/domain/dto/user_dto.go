package _dto

import "time"

type ColorDTO struct {
	ID              uint64 `json:"id"`
	Code            string `json:"code"`
	Color           string `json:"color"`
	ContentColor    string `json:"contentColor"`
	BackgroundColor string `json:"backgroundColor"`
	ColorKey        string `json:"colorKey"`
	HexCode         string `json:"hexCode"`
}

type RoleDTO struct {
	ID         uint64    `json:"id"`
	RoleName   string    `json:"roleName"`
	RoleKey    uint32    `json:"roleKey"`
	ProfileID  uint64    `json:"profileId"`
	Color      *ColorDTO `json:"color"`
	Status     string    `json:"status"`
	StatusText string    `json:"statusText"`
}

type AuthUserProfileDTO struct {
	ProfileID uint64   `json:"profileId"`
	FullName  string   `json:"fullName"`
	Avatar    string   `json:"avatar"`
	Role      *RoleDTO `json:"role"`
	RoleID    uint64   `json:"roleId"`
	IsLocked  bool     `json:"isLocked"`
	LockType  string   `json:"lockType"`
}

type UserV3DTO struct {
	ProfileID         uint64        `json:"profileId"`
	FullName          string        `json:"fullName"`
	Avatar            string        `json:"avatar"`
	RoleRealEstate    uint32        `json:"roleRealEstate"`
	MainAreaIds       []uint64      `json:"mainAreaIds"`
	PurposeUseIDs     []uint64      `json:"purposeUseIds"`
	MainAreaNews      []ItemDTO     `json:"mainAreaNews"`
	Visibility        uint32        `json:"visibility"`
	ProfileVisibility uint32        `json:"profileVisibility"`
	StatusOnline      uint32        `json:"statusOnline"`
	FacebookURL       string        `json:"facebookUrl"`
	InstagramURL      string        `json:"instagramUrl"`
	TwitterURL        string        `json:"twitterUrl"`
	LinkedinURL       string        `json:"linkedinUrl"`
	YoutubeURL        string        `json:"youtubeUrl"`
	WebsiteURL        string        `json:"websiteUrl"`
	ZaloURL           string        `json:"zaloUrl"`
	BlogURL           string        `json:"blogUrl"`
	OtherURL          string        `json:"otherUrl"`
	Email             string        `json:"email"`
	Phone             string        `json:"phone"`
	Phone2            string        `json:"phone2"`
	TaxCode           string        `json:"taxCode"`
	BackgroundImage   string        `json:"backgroundImage"`
	Position          string        `json:"position"`
	DepartmentID      uint64        `json:"departmentId"`
	UpdateTags        bool          `json:"updateTags"`
	Tags              []ItemDTO     `json:"tags"`
	SignatureVisible  bool          `json:"signatureVisible"`
	ProvinceID        *uint64       `json:"provinceId"`
	WardID            *uint64       `json:"wardId"`
	Gender            uint32        `json:"gender"`
	Birth             *time.Time    `json:"birth"`
	Address           *AddressV3DTO `json:"address"`
}
