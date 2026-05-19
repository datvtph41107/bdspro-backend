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

	_redis "common/redis"
	"tqd/config"

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

	config.LoadConfig()
	logger := log.New(os.Stdout, "[pmtiles] ", log.LstdFlags)
	redisSvc := _redis.NewRedisService()

	blockedTilesets := make(map[string]struct{})
	for _, name := range validatePMTilesDir(*tilesDir, logger) {
		blockedTilesets[name] = struct{}{}
	}

	// ── Khởi tạo PMTiles server ────────────────────────────────
	bucketURL := fmt.Sprintf("file://%s", *tilesDir)

	server, err := pmtiles.NewServer(bucketURL, "", logger, *cacheSize, *publicURL)
	if err != nil {
		logger.Fatalf("Không thể khởi tạo pmtiles server: %v", err)
	}
	server.Start()
	logger.Printf("PMTiles server đã khởi động")
	logger.Printf("Thư mục tiles : %s", *tilesDir)
	logger.Printf("Public URL    : %s", *publicURL)
	logger.Printf("Tile ci_*     : header %s / query %s → Redis ss:k:{sessionK}", tileSessionIDHeader, tileSessionIDQuery)

	// ── HTTP handler ───────────────────────────────────────────
	corsHeaders := fmt.Sprintf("Content-Type, Range, %s", tileSessionIDHeader)
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", corsHeaders)
			w.Header().Set("Access-Control-Expose-Headers", tileEncryptedHeader)
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			start := time.Now()
			tileset := tilesetNameFromPath(r.URL.Path)
			if tileset != "" {
				if _, blocked := blockedTilesets[tileset]; blocked {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusServiceUnavailable)
					_, _ = w.Write([]byte(`{"error":"tileset unavailable: archive failed validation"}`))
					logTileRequest(logger, start, r, http.StatusServiceUnavailable, 0)
					return
				}
			}

			statusCode, headers, body := server.Get(r.Context(), r.URL.Path)

			if isCIEncryptedTileset(tileset) && statusCode >= 200 && statusCode < 300 && len(body) > 0 {
				sId := sIdFromRequest(r)
				enc, err := tileEncryptorFromSessionID(r.Context(), redisSvc, sId)
				if err != nil || enc == nil {
					logger.Printf("ci_* tile: invalid sId=%q err=%v", sId, err)
					writeTileUnauthorized(w)
					logTileRequest(logger, start, r, http.StatusUnauthorized, 0)
					return
				}
				encrypted, err := encryptTileBody(r.Context(), enc, body)
				if err != nil {
					logger.Printf("AES encrypt failed for %s: %v", r.URL.Path, err)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"error":"tile encryption failed"}`))
					return
				}
				body = encrypted
				stripHeadersAfterEncrypt(headers)
				headers[tileEncryptedHeader] = tileEncryptedHeaderVal
			}

			for k, v := range headers {
				w.Header().Set(k, v)
			}
			w.WriteHeader(statusCode)
			_, _ = w.Write(body)
			logTileRequest(logger, start, r, statusCode, len(body))
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, indexHTML, *publicURL, *publicURL, *publicURL)
	})

	httpServer := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      corsMiddleware(mux),
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

func logTileRequest(logger *log.Logger, start time.Time, r *http.Request, statusCode, bodyLen int) {
	logger.Printf("%s | %v | %s %s | %d | %d bytes",
		time.Now().Format(time.RFC3339),
		time.Since(start),
		r.Method,
		r.URL.Path,
		statusCode,
		bodyLen,
	)
}

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

Tileset ci_*: header X-Session-Id hoặc query sId (= sessionK từ POST /v2/user/profile/session)
  </pre>
  <h2>Ví dụ</h2>
  <pre>
curl -H "X-Session-Id: 8472910384729103847" %s/ci_layer711.json
curl "%s/ci_layer711.json?sId=8472910384729103847"
  </pre>
</body>
</html>
`
