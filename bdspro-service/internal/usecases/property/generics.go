package property_usecases

import (
	"context"
	"fmt"
	"time"
)

type Handler interface {
	Name() string
	FKColumn() string
	Present() bool
	Insert(ctx context.Context) (uint64, error)
	Update(ctx context.Context, id uint64) error
}

func RunDirect(ctx context.Context, currentFK *uint64, h Handler) (newFK *uint64, err error) {
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

type Fields map[string]any

func (f Fields) SetStr(col string, v *string) {
	if v != nil {
		f[col] = *v
	}
}
func (f Fields) SetInt32(col string, v *int32) {
	if v != nil {
		f[col] = *v
	}
}
func (f Fields) SetUint32(col string, v *uint32) {
	if v != nil {
		f[col] = *v
	}
}
func (f Fields) SetFloat64(col string, v *float64) {
	if v != nil {
		f[col] = *v
	}
}
func (f Fields) SetUint64(col string, v *uint64) {
	if v != nil {
		f[col] = *v
	}
}
func (f Fields) SetDate(col string, v *string) {
	if v == nil {
		return
	}
	t, err := time.Parse(time.RFC3339, *v)
	if err != nil {
		return
	}
	f[col] = t
}

func DerefStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
func DerefInt32(p *int32) int32 {
	if p == nil {
		return 0
	}
	return *p
}
func DerefUint32(p *uint32) uint32 {
	if p == nil {
		return 0
	}
	return *p
}
func ParseDate(v *string) *time.Time {
	if v == nil {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *v)
	if err != nil {
		return nil
	}
	return &t
}
