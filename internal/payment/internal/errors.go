package service

import (
	_errors "common/errors"

	"google.golang.org/grpc/codes"
)

var (
	AdminPageSizeOutOfRange = _errors.MustSpec(
		510001,
		"PAYMENT_ADMIN_PAGE_SIZE_OUT_OF_RANGE",
		"page size must be between 1 and 100",
		codes.InvalidArgument,
		_errors.LegacyProblemCode("payment.admin.page_size_out_of_range"),
	)
)
