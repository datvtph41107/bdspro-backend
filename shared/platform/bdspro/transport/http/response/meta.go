package http_resp

import (
	"context"
	"time"

	corectx "bdspro/core/context"
	"bdspro/core/pkg/clock"
	"bdspro/core/pkg/query"
)

type MetaBuilder struct {
	clk clock.Clock
}

func NewMetaBuilder(clk clock.Clock) *MetaBuilder {
	return &MetaBuilder{clk: clk}
}

func (b *MetaBuilder) ItemMeta(ctx context.Context) ItemMeta {
	return ItemMeta{
		RequestID: corectx.RequestID(ctx),
		Timestamp: b.clk.Now().Format(time.RFC3339),
	}
}

func (b *MetaBuilder) ListMeta(ctx context.Context, p query.PageInfo) ListMeta {
	return ListMeta{
		RequestID:  corectx.RequestID(ctx),
		Timestamp:  b.clk.Now().Format(time.RFC3339),
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      p.Total,
		TotalPages: p.TotalPages(),
		HasNext:    p.HasNext(),
		HasPrev:    p.HasPrev(),
	}
}
