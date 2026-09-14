package admin

import "common/fault"

var ErrPageSizeOutOfRange = fault.Validation(
	"payment.admin.page_size_out_of_range",
	"page size must be between 1 and 100",
	fault.FieldViolation{Field: "page_size", Description: "must be between 1 and 100"},
)
