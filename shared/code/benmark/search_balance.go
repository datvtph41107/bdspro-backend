package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

var searchTexts = []string{
	"iphone", "laptop", "hoa", "cha", "docker",
	"kubernetes", "realm", "database", "microservice", "gateway",
}

type Stats struct {
	totalRequests int64
	success       int64
	failed        int64
	totalLatency  int64 // nanoseconds
	minLatency    int64
	maxLatency    int64
}

func worker(ctx context.Context, id int, m int, wg *sync.WaitGroup, client *http.Client, stats *Stats) {
	defer wg.Done()

	for i := 0; i < m; i++ {

		// Random sleep 100ms - 1s
		sleep := time.Duration(100+rand.Intn(500)) * time.Millisecond
		time.Sleep(sleep)

		text := searchTexts[rand.Intn(len(searchTexts))]
		u := "http://14.225.210.29:8000/v2/user/profile/search/public?text=" + url.QueryEscape(text)

		start := time.Now()
		resp, err := client.Get(u)
		latency := time.Since(start)

		atomic.AddInt64(&stats.totalRequests, 1)
		atomic.AddInt64(&stats.totalLatency, latency.Nanoseconds())

		if err != nil {
			atomic.AddInt64(&stats.failed, 1)
			fmt.Printf("Worker %d | ERROR | %v\n", id, err)
			continue
		}

		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			atomic.AddInt64(&stats.success, 1)
		} else {
			atomic.AddInt64(&stats.failed, 1)
		}

		// update min/max
		updateMinMax(stats, latency.Nanoseconds())

		fmt.Printf("Worker %d | q=%s | %v\n", id, text, latency)
	}
}

func updateMinMax(stats *Stats, latency int64) {
	for {
		min := atomic.LoadInt64(&stats.minLatency)
		if min == 0 || latency < min {
			if atomic.CompareAndSwapInt64(&stats.minLatency, min, latency) {
				break
			}
		} else {
			break
		}
	}

	for {
		max := atomic.LoadInt64(&stats.maxLatency)
		if latency > max {
			if atomic.CompareAndSwapInt64(&stats.maxLatency, max, latency) {
				break
			}
		} else {
			break
		}
	}
}

// func main() {
// 	rand.Seed(time.Now().UnixNano())

// 	n := 300 // số worker
// 	m := 20  // mỗi worker chạy m lần

// 	client := &http.Client{
// 		Timeout: 5 * time.Second,
// 	}

// 	stats := &Stats{}
// 	wg := sync.WaitGroup{}
// 	ctx := context.Background()

// 	startAll := time.Now()

// 	for i := 0; i < n; i++ {
// 		wg.Add(1)
// 		go worker(ctx, i, m, &wg, client, stats)
// 	}

// 	wg.Wait()

// 	totalTime := time.Since(startAll).Seconds()

// 	totalReq := atomic.LoadInt64(&stats.totalRequests)
// 	success := atomic.LoadInt64(&stats.success)
// 	failed := atomic.LoadInt64(&stats.failed)
// 	totalLatency := atomic.LoadInt64(&stats.totalLatency)

// 	avgLatency := time.Duration(totalLatency / totalReq)
// 	minLatency := time.Duration(atomic.LoadInt64(&stats.minLatency))
// 	maxLatency := time.Duration(atomic.LoadInt64(&stats.maxLatency))

// 	qps := float64(totalReq) / totalTime

// 	fmt.Println("\n===== RESULT =====")
// 	fmt.Println("Total Requests:", totalReq)
// 	fmt.Println("Success:", success)
// 	fmt.Println("Failed:", failed)
// 	fmt.Println("Total Time:", totalTime, "seconds")
// 	fmt.Println("QPS:", qps)
// 	fmt.Println("Avg Latency:", avgLatency)
// 	fmt.Println("Min Latency:", minLatency)
// 	fmt.Println("Max Latency:", maxLatency)
// }
