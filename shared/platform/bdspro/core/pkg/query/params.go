package query

type Filter interface {
	isFilter() // unexported — chỉ để enforce implement explicit
}

// Params là contract giữa transport và domain.
// Handler build Params → truyền vào usecase → repo dùng để query.
// Không bao giờ truyền *http.Request vào usecase/repo.
type Params struct {
	Pagination Pagination
	Sort       Sort
	// Filters mở rộng per-domain bằng embedding
}

// WithTotal tạo PageInfo để pass vào response.Writer.
func (p Params) WithTotal(total int64) PageInfo {
	return PageInfo{
		Page:    p.Pagination.Page,
		PerPage: p.Pagination.PerPage,
		Total:   total,
	}
}

// PageInfo là output — truyền vào response layer.
type PageInfo struct {
	Page    int
	PerPage int
	Total   int64
}

func (p PageInfo) TotalPages() int {
	if p.PerPage <= 0 {
		return 0
	}
	return int((p.Total + int64(p.PerPage) - 1) / int64(p.PerPage))
}
func (p PageInfo) HasNext() bool { return p.Page < p.TotalPages() }
func (p PageInfo) HasPrev() bool { return p.Page > 1 }
func (p PageInfo) Offset() int64 { return int64((p.Page - 1) * p.PerPage) }
