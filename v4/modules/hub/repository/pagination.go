package repository

// pageOffset maps 1-based FE page/size onto SQL OFFSET/LIMIT.
// page<=0 is treated as page 1 (first page).
func pageOffset(page, size int) (offset, limit int) {
	limit = size
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	return (page - 1) * limit, limit
}
