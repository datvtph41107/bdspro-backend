package controller

import (
	"net/http"

	"userv4/internal/requeststats"
	"userv4/modules/hub/services"
)

// Handler is the modules/hub HTTP entry for all hub capabilities.
type Handler struct {
	location          *LocationController
	locationV2        *LocationV2Controller
	userGuide         *UserGuideController
	faq               *FAQController
	version           *VersionController
	requestStats      *RequestStatsController
	errorLog          *ErrorLogController
	applink           *ApplinkController
	systemConfig      *SystemConfigController
	eventQueue        *EventQueueController
	interactiveEvent  *InteractiveEventController
	clickURL          *ClickURLController
	updateData        *UpdateDataController
	keywordSearch     *KeywordSearchController
	searchLatency     *SearchLatencyController
	apiKey            *APIKeyController
	statsSvc          services.RequestStatsService
}

func NewHandler(
	locationSvc services.LocationService,
	locationV2Svc services.LocationV2Service,
	userGuideSvc services.UserGuideService,
	faqSvc services.FAQService,
	versionSvc services.VersionService,
	requestStatsSvc services.RequestStatsService,
	errorLogSvc services.ErrorLogService,
	applinkSvc services.ApplinkService,
	systemConfigSvc services.SystemConfigService,
	eventQueueSvc services.EventQueueService,
	interactiveEventSvc services.InteractiveEventService,
	clickURLSvc services.ClickURLService,
	updateDataSvc services.UpdateDataService,
	keywordSearchSvc services.KeywordSearchService,
	searchLatencySvc services.SearchLatencyService,
	apiKeySvc services.APIKeyService,
) *Handler {
	return &Handler{
		location:         NewLocationController(locationSvc),
		locationV2:       NewLocationV2Controller(locationV2Svc),
		userGuide:        NewUserGuideController(userGuideSvc),
		faq:              NewFAQController(faqSvc),
		version:          NewVersionController(versionSvc),
		requestStats:     NewRequestStatsController(requestStatsSvc),
		errorLog:         NewErrorLogController(errorLogSvc),
		applink:          NewApplinkController(applinkSvc),
		systemConfig:     NewSystemConfigController(systemConfigSvc),
		eventQueue:       NewEventQueueController(eventQueueSvc),
		interactiveEvent: NewInteractiveEventController(interactiveEventSvc),
		clickURL:         NewClickURLController(clickURLSvc),
		updateData:       NewUpdateDataController(updateDataSvc),
		keywordSearch:    NewKeywordSearchController(keywordSearchSvc),
		searchLatency:    NewSearchLatencyController(searchLatencySvc),
		apiKey:           NewAPIKeyController(apiKeySvc),
		statsSvc:         requestStatsSvc,
	}
}

// RequestStatsCollector returns the shared in-process collector used by HTTP middleware.
func (h *Handler) RequestStatsCollector() *requeststats.Collector {
	if h == nil || h.statsSvc == nil {
		return nil
	}
	return h.statsSvc.Collector()
}

// RegisterHTTP mounts all hub routes.
func (h *Handler) RegisterHTTP(mux *http.ServeMux) {
	h.location.RegisterHTTP(mux)
	h.locationV2.RegisterHTTP(mux)
	h.userGuide.RegisterHTTP(mux)
	h.faq.RegisterHTTP(mux)
	h.version.RegisterHTTP(mux)
	h.requestStats.RegisterHTTP(mux)
	h.errorLog.RegisterHTTP(mux)
	h.applink.RegisterHTTP(mux)
	h.systemConfig.RegisterHTTP(mux)
	h.eventQueue.RegisterHTTP(mux)
	h.interactiveEvent.RegisterHTTP(mux)
	h.clickURL.RegisterHTTP(mux)
	h.updateData.RegisterHTTP(mux)
	h.keywordSearch.RegisterHTTP(mux)
	h.searchLatency.RegisterHTTP(mux)
	h.apiKey.RegisterHTTP(mux)
}
