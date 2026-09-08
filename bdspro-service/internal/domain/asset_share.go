package domain

import "strings"

type AssetShare struct {
	AssetID     uint64 `json:"assetId" binding:"required"`
	TargetID    uint64 `json:"targetId" binding:"required"`
	TypeShare   uint32 `json:"typeShare" binding:"required"`
	Permissions string `json:"permissions"`
}

func (AssetShare) TableName() string {
	return "asset_shares"
}

// Chuyển đổi khi lưu vào DB
func (p *AssetShare) SetFields(fields []string) {
	p.Permissions = strings.Join(fields, ",")
}

// Chuyển đổi khi đọc từ DB
func (p *AssetShare) GetFields() []string {
	if p.Permissions == "" {
		return []string{}
	}
	return strings.Split(p.Permissions, ",")
}
