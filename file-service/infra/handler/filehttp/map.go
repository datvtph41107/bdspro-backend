package filehttp

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var layerAllowedRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

var tileContentTypeByExt = map[string]string{
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".webp": "image/webp",
	".pbf":  "application/x-protobuf",
}

var tileCandidateExts = []string{".png", ".jpg", ".jpeg", ".webp", ".pbf"}

// registerMapRoutes restores the public map tile/config/font surface that lived
// in file-service/main.go before the migrate refactor dropped it.
func registerMapRoutes(r *gin.Engine) {
	// Tile map: /v1/map/layer/:layer/:z/:x/:y  (XYZ: z / x / y)
	// Reads files from: ./files/public/map/<layer>/<z>/<x>/<y>.(png|jpg|webp|pbf)
	r.GET("/v1/map/layer/:layer/:z/:x/:y", serveMapLayerTile)

	// Config files: /v1/map/config/<name>.json
	// Reads json config from: ./files/public/config/<name>.json
	r.GET("/v1/map/config/:name", serveMapConfig)

	r.GET("/v1/map/fonts/:fontstack/:range", serveMapFont)

	// Mid tile map: /v1/map/mid/:layer/:z/:x/:y
	// Cache files at: ./files/public/map/<layer>/<z>/<x>/<y>.(png|jpg|webp|pbf)
	// If missing, download from meeymap then cache.
	r.GET("/v1/map/mid/:layer/:z/:x/:y", serveMapMidTile)
}

func serveMapLayerTile(c *gin.Context) {
	layer := c.Param("layer")
	zStr := c.Param("z")
	xStr := c.Param("x")
	yParam := c.Param("y")

	if layer == "" || zStr == "" || xStr == "" || yParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "layer, z, y, x are required",
		})
		return
	}
	if !layerAllowedRe.MatchString(layer) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid layer",
		})
		return
	}

	z, err := strconv.ParseInt(zStr, 10, 64)
	if err != nil || z < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid z",
		})
		return
	}
	x, err := strconv.ParseInt(xStr, 10, 64)
	if err != nil || x < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid x",
		})
		return
	}

	if yParam == "." || yParam == ".." {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid y",
		})
		return
	}
	if ext := filepath.Ext(yParam); ext == "" {
		if yNum, err := strconv.ParseInt(yParam, 10, 64); err != nil || yNum < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid y",
			})
			return
		}
	}

	baseDir := filepath.Join(".", "files", "public", "map", layer, strconv.FormatInt(z, 10), strconv.FormatInt(x, 10))
	log.Printf("baseDir: %s", baseDir)

	if ext := filepath.Ext(yParam); ext != "" {
		base := strings.TrimSuffix(yParam, ext)
		y, err := strconv.ParseInt(base, 10, 64)
		if err != nil || y < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid y",
			})
			return
		}

		candidate := filepath.Join(baseDir, yParam)
		log.Printf("candidate: %s", candidate)
		// Chỉ serve file thường: nếu path là thư mục, http.ServeFile (c.File) trả 301 thêm "/" — client tile sẽ lỗi.
		if fi, statErr := os.Stat(candidate); statErr == nil && fi.Mode().IsRegular() {
			if ct, ok := tileContentTypeByExt[ext]; ok {
				c.Header("Content-Type", ct)
			}
			c.Header("Cache-Control", "public, max-age=86400")
			c.File(candidate)
			return
		}
	}

	for _, ext := range tileCandidateExts {
		candidate := filepath.Join(baseDir, yParam+ext)
		if fi, statErr := os.Stat(candidate); statErr == nil && fi.Mode().IsRegular() {
			if ct, ok := tileContentTypeByExt[ext]; ok {
				c.Header("Content-Type", ct)
			}
			c.Header("Cache-Control", "public, max-age=86400")
			c.File(candidate)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "tile not found",
	})
}

func serveMapConfig(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	if strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.Contains(name, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
		return
	}
	if filepath.Ext(name) != ".json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only .json is allowed"})
		return
	}

	fullPath := filepath.Join(".", "files", "public", "config", name)
	if _, err := os.Stat(fullPath); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Cache-Control", "public, max-age=300")
	c.File(fullPath)
}

func serveMapFont(c *gin.Context) {
	fontstack := c.Param("fontstack")
	rangeFile := c.Param("range")

	if fontstack == "" || rangeFile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
		return
	}

	if strings.Contains(fontstack, "..") || strings.Contains(fontstack, "\\") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fontstack"})
		return
	}

	fullPath := filepath.Join(".", "files", "public", "fonts", fontstack, rangeFile)

	if _, err := os.Stat(fullPath); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "font not found"})
		return
	}

	c.Header("Content-Type", "application/x-protobuf")
	c.Header("Cache-Control", "public, max-age=86400")

	c.File(fullPath)
}

func serveMapMidTile(c *gin.Context) {
	layer := c.Param("layer")
	zStr := c.Param("z")
	xStr := c.Param("x")
	yParam := c.Param("y")

	if layer == "" || zStr == "" || xStr == "" || yParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "layer, z, y, x are required"})
		return
	}
	if !layerAllowedRe.MatchString(layer) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid layer"})
		return
	}

	z, err := strconv.ParseInt(zStr, 10, 64)
	if err != nil || z < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid z"})
		return
	}
	x, err := strconv.ParseInt(xStr, 10, 64)
	if err != nil || x < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid x"})
		return
	}

	if yParam == "." || yParam == ".." {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid y"})
		return
	}
	if ext := filepath.Ext(yParam); ext == "" {
		if yNum, err := strconv.ParseInt(yParam, 10, 64); err != nil || yNum < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid y"})
			return
		}
	}

	baseDir := filepath.Join(".", "files", "public", "map", layer, zStr, xStr)

	var candidateNames []string
	if ext := filepath.Ext(yParam); ext != "" {
		base := strings.TrimSuffix(yParam, ext)
		y, err := strconv.ParseInt(base, 10, 64)
		if err != nil || y < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid y"})
			return
		}
		_ = y
		candidateNames = []string{yParam}
	} else {
		for _, ext := range tileCandidateExts {
			candidateNames = append(candidateNames, yParam+ext)
		}
	}

	for _, name := range candidateNames {
		fullPath := filepath.Join(baseDir, name)
		if fi, statErr := os.Stat(fullPath); statErr == nil && fi.Mode().IsRegular() {
			if ct, ok := tileContentTypeByExt[filepath.Ext(name)]; ok {
				c.Header("Content-Type", ct)
			}
			c.Header("Cache-Control", "public, max-age=86400")
			c.File(fullPath)
			return
		}
	}

	layerName := "HT_TOANQUOC"
	ext := "pbf"
	switch layer {
	case "vn_v1":
		layerName = "vietnam-latest-v28"
	case "KHSDD_HANOI":
		layerName = "KHSDD_HANOI"
	case "QHSDD_HANOI":
		layerName = "QHSDD_HANOI"
	case "qh_sdd":
		layerName = "qh_sdd"
	}

	download := "https://mapv3.meeymap.com/data/" + layerName + "/{z}/{x}/{y}." + ext
	log.Printf("download: %s, z=%d, x=%d, y=%s", layerName, z, x, yParam)

	client := &http.Client{Timeout: 30 * time.Second}
	buildDownloadURL := func(tileName string) (string, error) {
		if strings.Contains(download, "{z}") || strings.Contains(download, "{x}") || strings.Contains(download, "{y}") || strings.Contains(download, "{layer}") {
			u := download
			u = strings.ReplaceAll(u, "{layer}", layer)
			u = strings.ReplaceAll(u, "{z}", zStr)
			u = strings.ReplaceAll(u, "{x}", xStr)
			u = strings.ReplaceAll(u, "{y}", strings.TrimSuffix(tileName, filepath.Ext(tileName)))
			return u, nil
		}
		base, err := url.Parse(download)
		if err != nil {
			return "", err
		}
		rel := &url.URL{Path: pathJoinURL(base.Path, layer, zStr, xStr, tileName)}
		return base.ResolveReference(rel).String(), nil
	}

	for _, name := range candidateNames {
		downloadURL, err := buildDownloadURL(name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid download"})
			return
		}

		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, downloadURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		if resp.Body == nil {
			resp.Body.Close()
			continue
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			resp.Body.Close()
			continue
		}

		if err := os.MkdirAll(baseDir, 0o755); err != nil {
			resp.Body.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create cache dir"})
			return
		}

		finalPath := filepath.Join(baseDir, name)
		tmpPath := finalPath + ".tmp"
		f, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			resp.Body.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create cache file"})
			return
		}

		if resp.StatusCode == http.StatusNoContent {
			nameExt := strings.ToLower(filepath.Ext(name))
			if nameExt == ".png" {
				transparentPNG := []byte{
					0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
					0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
					0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
					0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
					0x89, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44,
					0x41, 0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00,
					0x00, 0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4,
					0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44,
					0xae, 0x42, 0x60, 0x82,
				}
				if _, err := f.Write(transparentPNG); err != nil {
					_ = f.Close()
					_ = resp.Body.Close()
					_ = os.Remove(tmpPath)
					return
				}
				_ = f.Close()
				_ = resp.Body.Close()

				if err := os.Rename(tmpPath, finalPath); err != nil {
					_ = os.Remove(tmpPath)
				}
				c.Header("Content-Type", "image/png")
				c.Header("Cache-Control", "public, max-age=86400")
				c.Data(http.StatusOK, "image/png", transparentPNG)
				return
			}

			_ = f.Close()
			_ = resp.Body.Close()

			if err := os.Rename(tmpPath, finalPath); err != nil {
				_ = os.Remove(tmpPath)
			}
			c.Status(http.StatusNoContent)
			return
		}

		ct := resp.Header.Get("Content-Type")
		if ct == "" {
			if byExt, ok := tileContentTypeByExt[filepath.Ext(name)]; ok {
				ct = byExt
			} else {
				ct = "application/octet-stream"
			}
		}
		c.Header("Content-Type", ct)
		c.Header("Cache-Control", "public, max-age=86400")

		tee := io.TeeReader(resp.Body, f)
		_, copyErr := io.Copy(c.Writer, tee)

		_ = f.Close()
		_ = resp.Body.Close()

		if copyErr != nil {
			_ = os.Remove(tmpPath)
			return
		}

		if err := os.Rename(tmpPath, finalPath); err != nil {
			_ = os.Remove(tmpPath)
		}
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "tile not found"})
}

func pathJoinURL(parts ...string) string {
	var out []string
	for _, p := range parts {
		p = strings.Trim(p, "/")
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	return "/" + strings.Join(out, "/")
}
