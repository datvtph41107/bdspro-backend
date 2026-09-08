package dto

import (
	"bdspro/internal/enums"
	_dto "common/domain/dto"
	"time"
)

type TxTransactionFilters struct {
	_dto.Pagable
	TransactionType   *string
	ApprovalStatus    *string
	TransactionStatus *uint32
	StartDate         *time.Time
	EndDate           *time.Time
	ToID              *uint64
	ToOf              enums.TxOwnerType
}

// type TransactionAction
type TxProcessDTO struct {
	Action        uint32     `gorm:"column:action"`      // mã action
	StatusName    string     `gorm:"column:status_name"` // tên trạng thái (map từ ActionNames)
	Timestamp     *time.Time `gorm:"column:created_at"`  // thời gian
	TransactionId uint64     `gorm:"column:tx_id"`       // id giao dịch
	Note          string     `gorm:"column:note"`        // ghi chú (nếu muốn lấy)
	Message       string     `gorm:"column:message"`     // message (nếu muốn lấy)
	Color         string     `gorm:"column:color"`       // custom (có thể join/mapping)
	BgColor       string     `gorm:"column:bg_color"`
	BorderColor   string     `gorm:"column:border_color"`
}
