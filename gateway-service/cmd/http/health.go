package cmd

import (
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type readinessState struct {
	ready atomic.Bool
}

func newReadinessState() *readinessState {
	state := &readinessState{}
	state.ready.Store(true)
	return state
}

func (s *readinessState) beginDrain() {
	if s != nil {
		s.ready.Store(false)
	}
}

func (s *readinessState) isReady() bool {
	return s != nil && s.ready.Load()
}

// registerOperationalHealth keeps liveness process-local and makes readiness
// reflect whether the Gateway is still accepting work. It intentionally does
// not couple readiness to every downstream service.
func registerOperationalHealth(router *gin.Engine, readiness *readinessState) {
	router.GET("/livez", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})
	router.GET("/readyz", func(c *gin.Context) {
		if !readiness.isReady() {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})
}
