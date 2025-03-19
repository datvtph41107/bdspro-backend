package _db

// Page represents pagination information
type Page struct {
	Limit int
	Skip  int
	Sort  string
}

// NewPage creates a new Page instance
func NewPage(skip, limit int, sort string) *Page {
	return &Page{Limit: limit, Skip: skip, Sort: sort}
}
