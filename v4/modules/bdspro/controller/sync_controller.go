package controller

import (
	"net/http"

	"userv4/modules/bdspro/services"
)

// SyncController owns GET /v2/bdspro/ts.
type SyncController struct {
	sync services.SyncService
}

func NewSyncController(sync services.SyncService) *SyncController {
	return &SyncController{sync: sync}
}

func (c *SyncController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /v2/bdspro/ts", c.httpGetAppTimestamps)
}

func (c *SyncController) httpGetAppTimestamps(w http.ResponseWriter, r *http.Request) {
	result, err := c.sync.GetAppTimestamps(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	// data = timestamps map; no pagination/timestamp at envelope root.
	if result == nil {
		writeOK(w, map[string]int64{})
		return
	}
	writeOK(w, result.Timestamps)
}
