package setting

import "hub/internal/dto"

func NewSystemConfigPersist() dto.SystemConfigPersist {
	return dto.SystemConfigPersist{
		Configs: map[string]string{
			"date_format":     "DD/MM/YYYY",
			"price_format":    "10",
			"currency":        "VND",
			"timezone":        "Asia/Ho_Chi_Minh",
			"language":        "vi",
			"items_per_page":  "20",
			"max_upload_size": "10",
		},
	}
}
