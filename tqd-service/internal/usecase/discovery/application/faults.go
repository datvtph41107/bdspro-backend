package application

import (
	_errors "common/errors"
	"tqd/internal"
)

func invalidCoordinateFault() error {
	return _errors.ReturnError(service.DiscoveryCoordinateInvalid)
}
