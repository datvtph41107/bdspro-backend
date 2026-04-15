package metrics

// // shared/pkg/metrics/metrics.go

// import (
//     "github.com/prometheus/client_golang/prometheus"
//     "github.com/prometheus/client_golang/prometheus/promauto"
// )

// // Lý do dùng promauto thay vì prometheus.MustRegister:
// // promauto tự đăng ký vào default registry khi khai báo,
// // không cần init function, không có registration race condition

// var (
//     // Histogram thay vì Counter cho latency:
//     // Histogram cho phép tính p50, p95, p99 — Counter chỉ cho average
//     // p99 latency cao trong khi average thấp → có 1% user trải nghiệm rất tệ
//     RequestDuration = promauto.NewHistogramVec(
//         prometheus.HistogramOpts{
//             Name:    "gateway_request_duration_seconds",
//             Help:    "Request duration in seconds",
//             Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
//         },
//         []string{"method", "service", "status"},
//     )

//     RequestsTotal = promauto.NewCounterVec(
//         prometheus.CounterOpts{
//             Name: "gateway_requests_total",
//             Help: "Total number of requests",
//         },
//         []string{"method", "service", "status"},
//     )

//     // Gauge cho connections hiện tại — cần theo dõi để phát hiện connection leak
//     ActiveConnections = promauto.NewGaugeVec(
//         prometheus.GaugeOpts{
//             Name: "gateway_active_connections",
//             Help: "Number of active connections",
//         },
//         []string{"service"},
//     )

//     // Upstream health gauge
//     UpstreamUp = promauto.NewGaugeVec(
//         prometheus.GaugeOpts{
//             Name: "gateway_upstream_up",
//             Help: "Whether upstream service is up (1=up, 0=down)",
//         },
//         []string{"service", "addr"},
//     )
// )
