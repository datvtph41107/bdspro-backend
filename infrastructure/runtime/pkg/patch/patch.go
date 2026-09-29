// Package patch giải quyết bài toán:
// "Chỉ write field thực sự thay đổi vào DB — không write thừa."
//
// # Bản chất
//
// Builder tích lũy (col → value) chỉ khi đồng thời thỏa 3 điều kiện:
//  1. mask.Allows(path)    → field được phép update (authorization)
//  2. incoming != nil      → client có gửi field (presence)
//  3. *incoming != existing → giá trị thực sự khác DB (diff)
//
// Ba điều kiện này map trực tiếp lên 3 concerns độc lập trong kiến trúc.
//
// # Method convention — mỗi method = 1 loại field = 1 quy ước zero/null
//
//	String        *string   nil=absent | ""=set empty (hợp lệ)
//	Int32         *int32    nil=absent | 0=set zero (hợp lệ, KHÔNG filter)
//	Bool          *bool     nil=absent | false=set false (hợp lệ)
//	ForeignKey    *uint64   nil=absent | 0=CLEAR FK (NULL) | >0=set id
//	NullableFloat *float64  nil=absent | present=set (0.0 hợp lệ)
//	NullableUint32 *uint32  nil=absent | present=set
//	DateString    *string   nil=absent | ""=CLEAR date (NULL) | RFC3339=set
//
// # Error handling
//
// Builder tích lũy error — không panic, không drop silently.
// Gọi Result() để lấy (map, error) thay vì Build() khi cần error check.
// Trong thực tế: DateString đã được validate ở handler → error ở đây là bug,
// nhưng vẫn propagate để không mất data silently.
//
// # Thread safety
//
// Builder KHÔNG thread-safe — 1 instance per request.
// Designed cho single-request lifecycle: New() → chain → Result() → discard.
//
// # Usage
//
//	updates, err := patch.WithMask(mask).
//	    String("info.title", "title", dto.Title, existing.Title).
//	    ForeignKey("info.project_id", "project_id", dto.ProjectID, existing.ProjectID).
//	    DateString("land_info.expiry.expired_land", "expired_land", dto.ExpiredLand, existing.ExpiredLand).
//	    Result()
//	if err != nil { return err }
//	if len(updates) == 0 { return nil } // nothing changed
//	return repo.UpdateFields(ctx, id, updates)
package patch

import (
	"common/pkg/fieldmask"
	"fmt"
	"time"
)

// Builder tích lũy field updates.
//
// Khởi tạo qua New() hoặc WithMask().
// Không dùng zero value trực tiếp.
type Builder struct {
	m    map[string]any
	mask fieldmask.FieldMask
	err  error // first error — subsequent errors dropped
}

// New tạo Builder không có mask — allow all.
// Dùng cho: Create operations, internal calls, tests.
func New() *Builder {
	return &Builder{
		m: make(map[string]any),
		// mask zero value = AllowAll
	}
}

// WithMask tạo Builder với FieldMask từ client request.
// Dùng cho: Update (PATCH) operations.
func WithMask(mask fieldmask.FieldMask) *Builder {
	return &Builder{
		m:    make(map[string]any),
		mask: mask,
	}
}

// Has trả về true nếu có ít nhất 1 field cần update.
// Dùng để skip repo call khi không có thay đổi:
//
//	b := patch.WithMask(mask).String(...).ForeignKey(...)
//	if !b.Has() { return nil }
//	return repo.UpdateFields(ctx, id, b.Build())
func (b *Builder) Has() bool { return len(b.m) > 0 }

// Build trả về map để truyền vào repo.UpdateFields.
// Nếu có error (từ DateString parse), Build() vẫn trả về map hiện tại.
// Dùng Result() nếu cần error check.
func (b *Builder) Build() map[string]any { return b.m }

// Result trả về (map, error).
// Dùng khi cần handle error từ Builder (e.g. DateString parse fail).
func (b *Builder) Result() (map[string]any, error) {
	return b.m, b.err
}

// ─────────────────────────────────────────────────────────────────────────────
// guard — điểm kết hợp duy nhất: mask + presence
// Tất cả method đều đi qua đây trước khi diff check.
// ─────────────────────────────────────────────────────────────────────────────

func (b *Builder) guard(path string, present bool) bool {
	return b.mask.Allows(path) && present
}

// ─────────────────────────────────────────────────────────────────────────────
// String
// Convention: nil=absent | ""=set empty (hợp lệ) | "v"=set v
// Bài toán: user clear title → gửi "" → server phải set "" (không filter zero)
// ─────────────────────────────────────────────────────────────────────────────

// String: existing là non-nullable string column.
func (b *Builder) String(path, col string, in *string, ex string) *Builder {
	if b.guard(path, in != nil) && *in != ex {
		b.m[col] = *in
	}
	return b
}

// StringPtr: existing là nullable string column (*string).
func (b *Builder) StringPtr(path, col string, in *string, ex *string) *Builder {
	if !b.guard(path, in != nil) {
		return b
	}
	exVal := ""
	if ex != nil {
		exVal = *ex
	}
	if *in != exVal {
		b.m[col] = *in
	}
	return b
}

// ─────────────────────────────────────────────────────────────────────────────
// Numeric — non-nullable columns
// Convention: nil=absent | 0=set zero (hợp lệ, KHÔNG filter)
// Bài toán: scope=0 là enum value hợp lệ → phải cho phép set về 0
// ─────────────────────────────────────────────────────────────────────────────

// Int32: existing là int32 column.
func (b *Builder) Int32(path, col string, in *int32, ex int32) *Builder {
	if b.guard(path, in != nil) && *in != ex {
		b.m[col] = *in
	}
	return b
}

// Uint32: existing là uint32 column.
func (b *Builder) Uint32(path, col string, in *uint32, ex uint32) *Builder {
	if b.guard(path, in != nil) && *in != ex {
		b.m[col] = *in
	}
	return b
}

// Int64: existing là int64 column.
func (b *Builder) Int64(path, col string, in *int64, ex int64) *Builder {
	if b.guard(path, in != nil) && *in != ex {
		b.m[col] = *in
	}
	return b
}

// ─────────────────────────────────────────────────────────────────────────────
// Bool
// Convention: nil=absent | false=set false (hợp lệ) | true=set true
// ─────────────────────────────────────────────────────────────────────────────

func (b *Builder) Bool(path, col string, in *bool, ex bool) *Builder {
	if b.guard(path, in != nil) && *in != ex {
		b.m[col] = *in
	}
	return b
}

// ─────────────────────────────────────────────────────────────────────────────
// ForeignKey — *uint64 nullable FK
// Convention: nil=absent | 0=CLEAR(NULL) | >0=set id
// Invariant: uint64 ID hợp lệ trong DB luôn > 0 (auto-increment từ 1)
// Bài toán: user bỏ project → gửi project_id=0 → server set NULL
// ─────────────────────────────────────────────────────────────────────────────

func (b *Builder) ForeignKey(path, col string, in *uint64, ex *uint64) *Builder {
	if !b.guard(path, in != nil) {
		return b
	}

	// CASE 1: clear FK (0 → NULL)
	if *in == 0 {
		// chỉ write nếu DB đang có value
		if ex != nil {
			b.m[col] = nil
		}
		return b
	}

	// CASE 2: set FK (>0)
	if ex == nil || *in != *ex {
		b.m[col] = *in
	}

	return b
}

// ─────────────────────────────────────────────────────────────────────────────
// Nullable numeric — *float64, *uint32
// Convention: nil=absent | present=set (zero là giá trị hợp lệ)
// Bài toán: area_total=0.0 là diện tích hợp lệ → không dùng 0=clear
// ─────────────────────────────────────────────────────────────────────────────

// NullableFloat: nil=absent | set value (kể cả 0.0).
func (b *Builder) NullableFloat(path, col string, in *float64, ex *float64) *Builder {
	if !b.guard(path, in != nil) {
		return b
	}
	if ex == nil || *in != *ex {
		b.m[col] = *in
	}
	return b
}

// NullableUint32: nil=absent | set value (kể cả 0).
func (b *Builder) NullableUint32(path, col string, in *uint32, ex *uint32) *Builder {
	if !b.guard(path, in != nil) {
		return b
	}
	if ex == nil || *in != *ex {
		b.m[col] = *in
	}
	return b
}

// ─────────────────────────────────────────────────────────────────────────────
// DateString — *string với dual-intent convention
// Convention: nil=absent | ""=CLEAR(NULL) | RFC3339=set date
// Bài toán: user xóa ngày hết hạn → gửi "" → server set NULL
//
// Note: format validation PHẢI chạy ở handler trước (validate.DateStringOpt).
// Parse error ở đây là programming error (validation không chạy) →
// propagate qua b.err thay vì silent drop.
// ─────────────────────────────────────────────────────────────────────────────

func (b *Builder) DateString(path, col string, in *string, ex *time.Time) *Builder {
	if !b.guard(path, in != nil) {
		return b
	}
	if *in == "" {
		// Clear intent: chỉ write nếu DB chưa NULL
		if ex != nil {
			b.m[col] = nil
		}
		return b
	}
	t, err := time.Parse(time.RFC3339, *in)
	if err != nil {
		// Không drop silently — đây là programming error
		// (validate.DateStringOpt phải đã chạy trước ở handler)
		if b.err == nil {
			b.err = fmt.Errorf("patch: %s: invalid date %q: %w", path, *in, err)
		}
		return b
	}
	if ex == nil || !t.Equal(*ex) {
		b.m[col] = t
	}
	return b
}
