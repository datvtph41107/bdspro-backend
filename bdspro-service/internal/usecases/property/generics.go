package property_usecases

import (
	"context"
	"fmt"
	"time"
)

// ─────────────────────────────────────────────────────────────────────────────
// Handler interface — contract cho mỗi sub-domain
//
// Bản chất: mỗi sub-domain có 2 operation:
//   Insert: tạo record mới → trả về ID để cập nhật FK lên lineage
//   Update: update record đã tồn tại → partial, chỉ field được gửi
//
// Handler KHÔNG biết về FieldMask — mask được truyền vào lúc khởi tạo
// qua patch.WithMask(mask) trong từng handler.
// ─────────────────────────────────────────────────────────────────────────────

// Handler là contract cho mỗi sub-domain có FK trên PropertyLineage.
type Handler interface {
	Name() string     // tên domain, dùng cho error wrapping: "info.insert"
	FKColumn() string // tên cột FK trên lineage: "property_info_id"
	Present() bool    // client có gửi sub-domain này không
	Insert(ctx context.Context) (uint64, error)
	Update(ctx context.Context, id uint64) error
}

// ─────────────────────────────────────────────────────────────────────────────
// RunDirect — cho DirectUpdate flow
//
// Logic:
//   !present         → skip (client không gửi sub-domain)
//   present + FK nil → Insert (record chưa tồn tại)
//   present + FK !nil → Update (record đã tồn tại)
//
// Trả về newFK != nil chỉ khi Insert thành công (cần update lineage)
// ─────────────────────────────────────────────────────────────────────────────

func RunDirect(ctx context.Context, currentFK *uint64, h Handler) (*uint64, error) {
	if !h.Present() {
		return nil, nil // client không gửi → skip hoàn toàn
	}
	if currentFK == nil {
		id, err := h.Insert(ctx)
		if err != nil {
			return nil, fmt.Errorf("%s.insert: %w", h.Name(), err)
		}
		return &id, nil
	}
	if err := h.Update(ctx, *currentFK); err != nil {
		return nil, fmt.Errorf("%s.update: %w", h.Name(), err)
	}
	return nil, nil // update không tạo FK mới
}

// ─────────────────────────────────────────────────────────────────────────────
// RunPropose — cho ProposeChange flow
//
// Logic: luôn Insert mới nếu present (bất kể FK cũ có không)
// Block không gửi → trả 0 (giữ FK từ bản gốc)
// ─────────────────────────────────────────────────────────────────────────────

func RunPropose(ctx context.Context, h Handler) (uint64, error) {
	if !h.Present() {
		return 0, nil // giữ FK từ origin lineage
	}
	id, err := h.Insert(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s.insert: %w", h.Name(), err)
	}
	return id, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// FKPatch — map FK cần update lên lineage
//
// Tích lũy trong directUpdate/proposeChange,
// flush 1 lần duy nhất vào LineageRepo.UpdateFields.
// ─────────────────────────────────────────────────────────────────────────────

// FKPatch tích lũy FK mới cần update lên lineage.
type FKPatch map[string]any

// Set thêm FK vào patch nếu id != nil.
func (p FKPatch) Set(col string, id *uint64) {
	if id != nil {
		p[col] = *id
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers — parse date, deref pointer
//
// ParseDate: dùng trong Insert (cần *time.Time cho domain entity)
// applyDateField: dùng trong Update (dùng patch.Builder.DateString thay thế)
// ─────────────────────────────────────────────────────────────────────────────

// ParseDate parse RFC3339 string → *time.Time.
// nil input → nil output (absent).
// "" input  → nil output (clear intent — caller xử lý riêng qua patch.Builder.DateString).
// Invalid   → nil (đã validate ở handler).
func ParseDate(v *string) *time.Time {
	if v == nil || *v == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *v)
	if err != nil {
		return nil // đã validate ở handler, err ở đây là programming error
	}
	return &t
}

// derefStr safe deref *string → string.
func derefStr(p *string) string {
	// "name": "" (define -> clear set value db = "")
	if p == nil {
		return ""
	}
	return *p
}

// derefUint32 safe deref *uint32 → uint32.
func derefUint32(p *uint32) uint32 {
	// "name": "" (define -> clear set value db = 0)
	if p == nil {
		return 0
	}
	return *p
}

// derefUint64Ptr safe deref *uint64 → *uint64 cho nullable FK.
// nil = absent → nil (không set)
// &0  = clear  → nil (clear FK)
// &>0 = set    → *uint64
func derefUint64Ptr(p *uint64) *uint64 {
	if p == nil || *p == 0 {
		return nil // absent hoặc clear → NULL trong DB
	}
	return p
}
