package processing

import "errors"

var ErrClaimLost = errors.New("report job claim is no longer owned")
