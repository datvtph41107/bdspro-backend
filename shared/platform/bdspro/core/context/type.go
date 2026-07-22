package context

import "time"

type RouteInfo struct {
	Path         string
	Service      string
	Timeout      time.Duration
	RequiresAuth bool
	RequiredRole string
}

type Principal struct {
	AuthID         uint64
	ProfileID      uint64
	OriginID       uint64
	OrganizationID *uint64
	PlanID         *uint64
	Role           string
	Type           string
	Session        uint64
}
