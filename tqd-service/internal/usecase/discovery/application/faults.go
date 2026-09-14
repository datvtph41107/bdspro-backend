package application

import "common/fault"

func invalidCoordinateFault() error {
	return fault.Validation(
		"tqd.discovery.coordinate_invalid",
		"invalid latitude/longitude",
	)
}
