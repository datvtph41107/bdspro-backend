// generics.go
package property_usecases

import (
	"bdspro/internal/modules"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Handler interface {
	Name() string
	FKColumn() string
	Present() bool
	Insert(ctx context.Context) (uint64, error)
	Update(ctx context.Context, id uint64) error
}

func RunDirect(ctx context.Context, currentFK *uint64, h Handler) (*uint64, error) {
	if !h.Present() {
		return nil, nil
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
	return nil, nil
}

func RunPropose(ctx context.Context, h Handler) (uint64, error) {
	if !h.Present() {
		return 0, nil
	}
	id, err := h.Insert(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s.insert: %w", h.Name(), err)
	}
	return id, nil
}

type FKPatch map[string]any

func (p FKPatch) Set(col string, id *uint64) {
	if id != nil {
		p[col] = *id
	}
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return errors.Is(err, gorm.ErrDuplicatedKey)
}

// ImpactRequiredError lỗi yêu cầu xác nhận
type ImpactRequiredError struct {
	Evaluation *modules.ImpactEvaluation
}

func (e *ImpactRequiredError) Error() string {
	return "impact evaluation requires confirmation"
}

// Helper functions
func ParseDate(v *string) *time.Time {
	if v == nil || *v == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *v)
	if err != nil {
		return nil
	}
	return &t
}

func derefStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func derefUint32(p *uint32) uint32 {
	if p == nil {
		return 0
	}
	return *p
}

func derefInt32(p *int32) int32 {
	if p == nil {
		return 0
	}
	return *p
}

func derefUint64Ptr(p *uint64) *uint64 {
	if p == nil || *p == 0 {
		return nil
	}
	return p
}

func derefInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

// Helper function để safe log
func safeDerefUint64(p *uint64) string {
	if p == nil {
		return "nil"
	}
	return fmt.Sprintf("%d", *p)
}

func safeDerefFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func derefUint64(p *uint64) uint64 {
	if p == nil {
		return 0
	}
	return *p
}

func marshalJSON(v any) string {
	switch val := v.(type) {
	case []string:
		if len(val) == 0 {
			return ""
		}
	case []uint64:
		if len(val) == 0 {
			return ""
		}
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	s := string(b)
	if s == "null" || s == "[]" {
		return ""
	}
	return s
}
