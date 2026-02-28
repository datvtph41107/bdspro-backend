package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	httpRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"method", "path"},
	)

	httpResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"method", "path"},
	)

	concurrentRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "concurrent_requests",
			Help: "Current number of concurrent requests",
		},
	)

	uploadFileSize = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "upload_file_size_bytes",
			Help: "Uploaded file size distribution",
			Buckets: []float64{
				1024,       // 1KB
				10240,      // 10KB
				102400,     // 100KB
				1048576,    // 1MB
				10485760,   // 10MB
				52428800,   // 50MB
				104857600,  // 100MB
				524288000,  // 500MB
				1073741824, // 1GB
			},
		},
	)

	uploadDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "upload_duration_seconds",
			Help:    "File upload duration",
			Buckets: []float64{0.1, 0.5, 1, 5, 10, 30, 60},
		},
	)

	cacheHitsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total cache hits",
		},
		[]string{"cache_type"},
	)

	cacheMissesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total cache misses",
		},
		[]string{"cache_type"},
	)

	storageOperationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "storage_operations_total",
			Help: "Total storage operations",
		},
		[]string{"operation", "storage_type"},
	)

	videoProcessingDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "video_processing_duration_seconds",
			Help:    "Video processing duration",
			Buckets: []float64{1, 5, 10, 30, 60, 300, 600},
		},
	)
)

func init() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		httpRequestSize,
		httpResponseSize,
		concurrentRequests,
		uploadFileSize,
		uploadDuration,
		cacheHitsTotal,
		cacheMissesTotal,
		storageOperationsTotal,
		videoProcessingDuration,
	)
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		method := c.Request.Method

		// Track concurrent requests
		concurrentRequests.Inc()
		defer concurrentRequests.Dec()

		// Track request size
		if c.Request.ContentLength > 0 {
			httpRequestSize.WithLabelValues(method, path).
				Observe(float64(c.Request.ContentLength))
		}

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get response size
		responseSize := 0
		if w, ok := c.Writer.(gin.ResponseWriter); ok {
			responseSize = w.Size()
		}

		// Record metrics
		status := strconv.Itoa(c.Writer.Status())
		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)
		httpResponseSize.WithLabelValues(method, path).Observe(float64(responseSize))

		// Special metrics for uploads
		if path == "/v1/files" && method == "POST" {
			if sizeStr := c.GetHeader("Content-Length"); sizeStr != "" {
				if size, err := strconv.ParseFloat(sizeStr, 64); err == nil {
					uploadFileSize.Observe(size)
				}
			}
			uploadDuration.Observe(duration)
		}
	}
}

// Cache metrics helpers
func RecordCacheHit(cacheType string) {
	cacheHitsTotal.WithLabelValues(cacheType).Inc()
}

func RecordCacheMiss(cacheType string) {
	cacheMissesTotal.WithLabelValues(cacheType).Inc()
}

func RecordStorageOperation(operation, storageType string) {
	storageOperationsTotal.WithLabelValues(operation, storageType).Inc()
}

func RecordVideoProcessing(duration time.Duration) {
	videoProcessingDuration.Observe(duration.Seconds())
}

// Custom metrics endpoint
func MetricsHandler(c *gin.Context) {
	metrics := prometheus.DefaultGatherer

	c.Writer.Header().Set("Content-Type", "text/plain")
	promhttp.HandlerFor(metrics, promhttp.HandlerOpts{}).ServeHTTP(c.Writer, c.Request)
}
