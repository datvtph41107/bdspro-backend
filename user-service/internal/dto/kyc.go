package dto

import (
	_dto "common/domain/dto"
)

// SubmitKYCRequest DTO cho submit KYC
type SubmitKYCRequest struct {
	FullName     string `json:"fullName" binding:"required"`
	IdentityCard string `json:"identityCard" binding:"required"`
	FrontImage   string `json:"frontImage" binding:"required"`
	BackImage    string `json:"backImage" binding:"required"`
	SelfieImage  string `json:"selfieImage"`
	IDNumber     string `json:"idNumber"`
	DateOfBirth  string `json:"dateOfBirth"`
	ExpiryDate   string `json:"expiryDate"`
}

// ListKYCsRequest DTO cho list KYC (Admin)
type ListKYCsRequest struct {
	_dto.Pagable
	Search    string `json:"search"`
	Status    uint32 `json:"status"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"`
}

// ApproveKYCRequest DTO cho approve KYC (Admin)
type ApproveKYCRequest struct {
	ID uint64 `json:"id" uri:"id" binding:"required"`
}

// RejectKYCRequest DTO cho reject KYC (Admin)
type RejectKYCRequest struct {
	ID     uint64 `json:"id" uri:"id" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}
