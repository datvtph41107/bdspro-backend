// common/quota/errors.go
package quota

import "errors"

var (
    ErrQuotaExceeded        = errors.New("quota exceeded")
    ErrReservationNotFound  = errors.New("reservation not found")
    ErrInvalidAction        = errors.New("invalid action")
    ErrUnknownService       = errors.New("unknown service")
)