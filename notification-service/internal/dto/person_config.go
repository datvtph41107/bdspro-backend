package dto

import (
	_enum "common/domain/enum"
	"fmt"
)

type PersonConfigDTO struct {
	ID        uint64
	UserID    uint64
	Key       string
	Checked   bool
	IsDefault bool
	Channel   _enum.EChannelNotification
}

func (p *PersonConfigDTO) Validate() error {
	if p == nil {
		return fmt.Errorf("person config không hợp lệ")
	}
	if p.Key == "" {
		return fmt.Errorf("key không được để trống")
	}
	if !p.Channel.IsValid() {
		return fmt.Errorf("channel không hợp lệ")
	}
	// if p.UserID == 0 {
	// 	return fmt.Errorf("userId không được để trống")
	// }
	return nil
}
