package _dto

import "time"

type ItemDTO struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	ItemType uint32 `json:"itemType"`

	FileName *string `json:"fileName"`
	FileURL  *string `json:"fileUrl"`
	FileType *string `json:"fileType"`
	Status   *uint32 `json:"status"`
	Issuer   *uint64 `json:"issuer"`
	IssuedAt *string `json:"issuedAt"`
}

type FriendV3DTO struct {
	ID         uint64  `json:"id"`
	ReceiverID uint64  `json:"receiverId"`
	CreatedBy  *uint64 `json:"createdBy"`
	Status     uint32  `json:"status"`
}

type RateStateV3DTO struct {
	TotalReviews uint64  `json:"totalReviews"`
	AverageScore float32 `json:"averageScore"`
	Star1Count   uint64  `json:"star1Count"`
	Star2Count   uint64  `json:"star2Count"`
	Star3Count   uint64  `json:"star3Count"`
	Star4Count   uint64  `json:"star4Count"`
	Star5Count   uint64  `json:"star5Count"`
}

type KYCV3DTO struct {
	ID           uint64     `json:"id"`
	ProfileID    uint64     `json:"profileId"`
	FullName     string     `json:"fullName"`
	IdentityCard string     `json:"identityCard"`
	FrontImage   string     `json:"frontImage"`
	BackImage    string     `json:"backImage"`
	SelfieImage  string     `json:"selfieImage"`
	Status       uint32     `json:"status"`
	RejectReason string     `json:"rejectReason"`
	ReviewedBy   *uint64    `json:"reviewedBy"`
	ReviewedAt   *time.Time `json:"reviewedAt"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
	IDNumber     string     `json:"idNumber"`
	DateOfBirth  string     `json:"dateOfBirth"`
	ExpiryDate   string     `json:"expiryDate"`
}

type TotalResponseDTO struct {
	User         *UserV3DTO
	KYC          *KYCV3DTO
	FriendStatus *FriendV3DTO
	RateStats    *RateStateV3DTO
}
