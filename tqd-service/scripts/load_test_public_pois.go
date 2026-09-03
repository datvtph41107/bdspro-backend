//go:build ignore

// Load test GET /v2/tqd/public/pois.
//
// Usage:
//
//	go run ./scripts/load_test_public_pois.go \
//	  -base http://localhost:8000 \
//	  -concurrency 20 \
//	  -requests 200 \
//	  -search "cà phê"
//
// Latency được ghi DB qua TrackSearchLatency (flush ~10s).
// Xem trên admin:
//
//	GET /v2/hub/search-latency/dashboard?path=/v2/tqd/public/pois&lastDates=1
//	GET /v2/hub/search-latency/daily?path=/v2/tqd/public/pois&lastDates=1&gap=hours
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type result struct {
	status   int
	duration time.Duration
	err      error
}

func main() {
	base := flag.String("base", "http://localhost:8000", "gateway base URL")
	path := flag.String("path", "/v2/tqd/public/pois", "API path")
	search := flag.String("search", "cà phê", "query search (cần non-empty để track keyword; latency luôn track)")
	mode := flag.String("mode", "suggest", "mode query")
	size := flag.Int("size", 20, "size query")
	concurrency := flag.Int("concurrency", 20, "số worker song song")
	requests := flag.Int("requests", 200, "tổng số request")
	timeout := flag.Duration("timeout", 30*time.Second, "HTTP timeout mỗi request")
	flag.Parse()

	if *concurrency < 1 {
		*concurrency = 1
	}
	if *requests < 1 {
		fmt.Fprintln(os.Stderr, "requests must be >= 1")
		os.Exit(1)
	}

	endpoint, err := url.Parse(*base + *path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid URL: %v\n", err)
		os.Exit(1)
	}
	q := endpoint.Query()
	q.Set("search", *search)
	q.Set("mode", *mode)
	q.Set("size", fmt.Sprintf("%d", *size))
	endpoint.RawQuery = q.Encode()

	client := &http.Client{Timeout: *timeout}
	jobs := make(chan int, *requests)
	results := make(chan result, *requests)

	var startedAt time.Time
	var wg sync.WaitGroup
	var inFlight atomic.Int64

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				inFlight.Add(1)
				start := time.Now()
				status, reqErr := doRequest(client, endpoint.String())
				results <- result{
					status:   status,
					duration: time.Since(start),
					err:      reqErr,
				}
				inFlight.Add(-1)
			}
		}()
	}

	startedAt = time.Now()
	for i := 0; i < *requests; i++ {
		jobs <- i
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	durations := make([]time.Duration, 0, *requests)
	statusCount := map[int]int{}
	var errCount int
	var okCount int

	for r := range results {
		if r.err != nil {
			errCount++
			statusCount[0]++
			continue
		}
		statusCount[r.status]++
		durations = append(durations, r.duration)
		if r.status >= 200 && r.status < 300 {
			okCount++
		}
	}
	elapsed := time.Since(startedAt)

	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })

	summary := map[string]any{
		"url":           endpoint.String(),
		"concurrency":   *concurrency,
		"requests":      *requests,
		"ok":            okCount,
		"transportErr":  errCount,
		"elapsedMs":     elapsed.Milliseconds(),
		"rps":           float64(*requests) / elapsed.Seconds(),
		"status":        statusCount,
		"latencyMs":     latencyStats(durations),
		"adminHint": map[string]string{
			"dashboard": fmt.Sprintf("%s/v2/hub/search-latency/dashboard?path=%s&lastDates=1", *base, url.QueryEscape(*path)),
			"daily":     fmt.Sprintf("%s/v2/hub/search-latency/daily?path=%s&lastDates=1&gap=hours", *base, url.QueryEscape(*path)),
			"note":      "Đợi ~10–15s sau khi chạy để hub flush latency vào DB trước khi gọi admin API.",
		},
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(summary)

	if okCount == 0 {
		os.Exit(1)
	}
}

func doRequest(client *http.Client, rawURL string) (int, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode, nil
}

func latencyStats(ds []time.Duration) map[string]any {
	if len(ds) == 0 {
		return map[string]any{"count": 0}
	}
	var sum time.Duration
	for _, d := range ds {
		sum += d
	}
	return map[string]any{
		"count": len(ds),
		"min":   ms(ds[0]),
		"avg":   ms(sum / time.Duration(len(ds))),
		"p50":   ms(percentile(ds, 50)),
		"p95":   ms(percentile(ds, 95)),
		"p99":   ms(percentile(ds, 99)),
		"max":   ms(ds[len(ds)-1]),
	}
}

func percentile(ds []time.Duration, p int) time.Duration {
	if len(ds) == 0 {
		return 0
	}
	if p <= 0 {
		return ds[0]
	}
	if p >= 100 {
		return ds[len(ds)-1]
	}
	idx := (p * len(ds) / 100)
	if idx >= len(ds) {
		idx = len(ds) - 1
	}
	return ds[idx]
}

func ms(d time.Duration) float64 {
	return float64(d.Microseconds()) / 1000.0
}
