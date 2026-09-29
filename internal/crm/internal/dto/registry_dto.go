package dto

import (
	_dto "common/domain/dto"
)

// RegistryCreateRequest represents request to create a registry entry
type RegistryCreateRequest struct {
	RegistryKey  uint32 `json:"registryKey" binding:"required"`
	RegistryName string `json:"registryName" binding:"required,max=20"`
	Port         int    `json:"port" binding:"required"`
	Service      string `json:"service" binding:"required"`
}

// RegistryUpdateRequest represents request to update a registry entry
type RegistryUpdateRequest struct {
	RegistryKey  uint32 `json:"registryKey" binding:"required"`
	RegistryName string `json:"registryName" binding:"required,max=20"`
	Port         int    `json:"port" binding:"required"`
	Service      string `json:"service" binding:"required"`
}

// RegistryResponse represents response for registry entry
type RegistryResponse struct {
	RegistryKey  uint32 `json:"registryKey"`
	RegistryName string `json:"registryName"`
	Port         int    `json:"port"`
	Service      string `json:"service"`
	_dto.Pagable
}

// RegistryListRequest represents request to get registry list
type RegistryListRequest struct {
	RegistryKey uint32 `json:"registryKey,omitempty"`
	Service     string `json:"service,omitempty"`
	_dto.Pagable
}

// RegistryListResponse represents response for registry list
type RegistryListResponse struct {
	Data  []RegistryResponse `json:"data"`
	Total int64              `json:"total"`
}