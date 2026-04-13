package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/protomaps/go-pmtiles/pmtiles"
)

func RunTileServer(port int) error {
	// ── CLI flags ──────────────────────────────────────────────
	tilesDir := flag.String("dir", "./files/tiles", "Thư mục chứa các file .pmtiles")
	publicURL := flag.String("public-url", "", "Public URL (vd: http://localhost:8032). Bắt buộc để sinh TileJSON")
	cacheSize := flag.Int("cache", 64, "Kích thước cache (MB)")
	flag.Parse()

	if *publicURL == "" {
		*publicURL = fmt.Sprintf("http://localhost:%d", port)
	}

	logger := log.New(os.Stdout, "[pmtiles] ", log.LstdFlags)

	// ── Khởi tạo PMTiles server ────────────────────────────────
	// NewServer nhận bucket URL dạng file:///path/to/dir
	bucketURL := fmt.Sprintf("file://%s", *tilesDir)

	server, err := pmtiles.NewServer(bucketURL, "", logger, *cacheSize, *publicURL)
	if err != nil {
		logger.Fatalf("Không thể khởi tạo pmtiles server: %v", err)
	}
	server.Start()
	logger.Printf("PMTiles server đã khởi động")
	logger.Printf("Thư mục tiles : %s", *tilesDir)
	logger.Printf("Public URL    : %s", *publicURL)

	// ── HTTP handler ───────────────────────────────────────────
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	// Trang chủ: liệt kê các tileset có sẵn
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			// Start timer
			start := time.Now()

			// Gọi pmtiles server
			statusCode, headers, body := server.Get(r.Context(), r.URL.Path)
			for k, v := range headers {
				w.Header().Set(k, v)
			}
			// CORS
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Range")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// Gửi response
			w.WriteHeader(statusCode)
			_, _ = w.Write(body)

			// 🔹 Log 1 dòng giống Gin
			duration := time.Since(start)
			logger.Printf("%s | %v | %s %s | %d | %d bytes",
				time.Now().Format(time.RFC3339),
				duration,
				r.Method,
				r.URL.Path,
				statusCode,
				len(body),
			)
			return
		}

		// Trang index
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, indexHTML, *publicURL)
	})

	// ── Graceful shutdown ──────────────────────────────────────
	httpServer := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Printf("HTTP server lắng nghe tại http://0.0.0.0:%d", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("ListenAndServe: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Đang tắt server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Printf("Lỗi khi tắt server: %v", err)
	}
	logger.Println("Server đã dừng.")
	return nil
}

// ── Template trang chủ ─────────────────────────────────────────
const indexHTML = `<!DOCTYPE html>
<html lang="vi">
<head>
  <meta charset="UTF-8"/>
  <title>PMTiles Server</title>
  <style>
    body { font-family: monospace; max-width: 800px; margin: 40px auto; padding: 0 20px; }
    h1   { color: #333; }
    code { background: #f0f0f0; padding: 2px 6px; border-radius: 4px; }
    pre  { background: #f0f0f0; padding: 16px; border-radius: 6px; overflow-x: auto; }
    a    { color: #0070f3; }
  </style>
</head>
<body>
  <h1>🗺️ PMTiles Server</h1>
  <p>Server đang chạy tại <code>%s</code></p>
  <h2>API Endpoints</h2>
  <pre>
GET /{tileset}.json          → TileJSON metadata
GET /{tileset}/{z}/{x}/{y}.mvt  → Vector tile
GET /{tileset}/{z}/{x}/{y}.png  → Raster tile
GET /health                  → Health check
  </pre>
  <h2>Ví dụ</h2>
  <pre>
curl %s/output1.json
curl %s/output2.json
curl %s/output1/0/0/0.mvt
  </pre>
</body>
</html>
`
