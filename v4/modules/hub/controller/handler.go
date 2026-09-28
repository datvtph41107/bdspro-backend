package controller

import (
	"net/http"

	"userv4/internal/requeststats"
	"userv4/modules/hub/services"
)

// Handler is the modules/hub HTTP entry for location, user guide, FAQ, version, and request-stats.
type Handler struct {
	location     *LocationController
	locationV2   *LocationV2Controller
	userGuide    *UserGuideController
	faq          *FAQController
	version      *VersionController
	requestStats *RequestStatsController
	statsSvc     services.RequestStatsService
}

func NewHandler(
	locationSvc services.LocationService,
	locationV2Svc services.LocationV2Service,
	userGuideSvc services.UserGuideService,
	faqSvc services.FAQService,
	versionSvc services.VersionService,
	requestStatsSvc services.RequestStatsService,
) *Handler {
	return &Handler{
		location:     NewLocationController(locationSvc),
		locationV2:   NewLocationV2Controller(locationV2Svc),
		userGuide:    NewUserGuideController(userGuideSvc),
		faq:          NewFAQController(faqSvc),
		version:      NewVersionController(versionSvc),
		requestStats: NewRequestStatsController(requestStatsSvc),
		statsSvc:     requestStatsSvc,
	}
}

// RequestStatsCollector returns the shared in-process collector used by HTTP middleware.
func (h *Handler) RequestStatsCollector() *requeststats.Collector {
	if h == nil || h.statsSvc == nil {
		return nil
	}
	return h.statsSvc.Collector()
}

// RegisterHTTP mounts hub routes with gateway path parity.
func (h *Handler) RegisterHTTP(mux *http.ServeMux) {
	h.location.RegisterHTTP(mux)
	h.locationV2.RegisterHTTP(mux)
	h.userGuide.RegisterHTTP(mux)
	h.faq.RegisterHTTP(mux)
	h.version.RegisterHTTP(mux)
	h.requestStats.RegisterHTTP(mux)
}
