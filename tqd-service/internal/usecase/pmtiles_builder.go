package usecase

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"tqd/internal/interface/repo"
)

type BuildProgress struct {
	Status       string
	Message      string
	CurrentZoom  int32
	TilesWritten int64
	TotalTiles   int64
	OutputPath   string
}

// BuildPMTilesFile generates PMTiles using tippecanoe with better progress handling and options.
func BuildPMTilesFile(
	ctx context.Context,
	regionRepo repo.RegionRepository,
	layerID uint64,
	minZoom, maxZoom uint32,
	outputDir string,
	progressFn func(BuildProgress),
) error {

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", outputDir, err)
	}

	outputPath := filepath.Join(outputDir, fmt.Sprintf("ci_layer%d.pmtiles", layerID))

	// Temp NDJSON file
	tmpFile, err := os.CreateTemp(outputDir, fmt.Sprintf("ci_layer%d_*.ndjson", layerID))
	if err != nil {
		return fmt.Errorf("create temp ndjson: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// ── Stream from DB to temp file ─────────────────────────────────────
	rows, err := regionRepo.StreamGeoJSON(ctx, layerID)
	if err != nil {
		return fmt.Errorf("stream geojson: %w", err)
	}
	defer rows.Close()

	writer := bufio.NewWriterSize(tmpFile, 2<<20) // 2MB buffer
	var count int64

	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return fmt.Errorf("scan row: %w", err)
		}
		if len(raw) == 0 {
			continue
		}

		if _, err := writer.Write(raw); err != nil {
			return fmt.Errorf("write feature: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("write newline: %w", err)
		}

		count++

		// Progress nhẹ nhàng hơn
		if count%5_000 == 0 || count == 1 {
			if err := writer.Flush(); err != nil {
				return fmt.Errorf("flush: %w", err)
			}
			// progressFn(BuildProgress{
			// 	Status:       "processing",
			// 	Message:      fmt.Sprintf("Streamed %d features to NDJSON", count),
			// 	TilesWritten: count,
			// })
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error after %d features: %w", count, err)
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("final flush: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("layer %d has no valid geometry", layerID)
	}

	// progressFn(BuildProgress{
	// 	Status:       "processing",
	// 	Message:      fmt.Sprintf("Starting tippecanoe with %d features...", count),
	// 	TilesWritten: count,
	// })

	// ── Tippecanoe command (cải thiện đáng kể) ─────────────────────────
	args := []string{
		"-o", outputPath,
		"--force",
		"--layer", fmt.Sprintf("layer%d", layerID),
		"--read-parallel",      // tốt khi input là file NDJSON
		"--no-feature-limit",   // quan trọng với polygon phức tạp
		"--no-tile-size-limit", // tránh drop quá mạnh ở low zoom
		"--drop-densest-as-needed",
		"--extend-zooms-if-still-dropping",
	}

	// Zoom logic
	if minZoom > 0 {
		args = append(args, "--minimum-zoom", fmt.Sprintf("%d", minZoom))
	} else {
		args = append(args, "--minimum-zoom", "6")
	}
	slog.InfoContext(ctx, fmt.Sprintf("maxZoom=%d", maxZoom))
	if maxZoom > 0 {
		args = append(args, "--maximum-zoom", fmt.Sprintf("%d", maxZoom))
	} else {
		args = append(args, "--maximum-zoom", "16") // tippecanoe tự guess maxzoom tốt hơn
	}

	args = append(args, tmpPath)

	cmd := exec.CommandContext(ctx, "tippecanoe", args...)

	// Pipe stdout & stderr
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start tippecanoe: %w", err)
	}

	// Đọc output không block cmd
	go readOutput("tippecanoe:stdout", stdoutPipe, progressFn)
	go readOutput("tippecanoe:stderr", stderrPipe, progressFn)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("tippecanoe failed (exit code %d): %w", cmd.ProcessState.ExitCode(), err)
	}

	// ── Done ───────────────────────────────────────────────────────────
	progressFn(BuildProgress{
		Status:       "done",
		Message:      fmt.Sprintf("PMTiles generated successfully from %d features", count),
		TilesWritten: count,
		OutputPath:   outputPath,
	})

	return nil
}

// BuildFamilyPMTilesFile generates a single PMTiles file for all layers in a family.
func BuildFamilyPMTilesFile(
	ctx context.Context,
	regionRepo repo.RegionRepository,
	familyID uint64,
	minZoom, maxZoom uint32,
	outputDir string,
	progressFn func(BuildProgress),
) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", outputDir, err)
	}

	outputPath := filepath.Join(outputDir, fmt.Sprintf("ci_family%d.pmtiles", familyID))
	tmpFile, err := os.CreateTemp(outputDir, fmt.Sprintf("ci_family%d_*.ndjson", familyID))
	if err != nil {
		return fmt.Errorf("create temp ndjson: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	rows, err := regionRepo.StreamGeoJSONByFamilyID(ctx, familyID)
	if err != nil {
		return fmt.Errorf("stream geojson by family: %w", err)
	}
	defer rows.Close()

	writer := bufio.NewWriterSize(tmpFile, 2<<20)
	var count int64

	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return fmt.Errorf("scan row: %w", err)
		}
		if len(raw) == 0 {
			continue
		}
		if _, err := writer.Write(raw); err != nil {
			return fmt.Errorf("write feature: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("write newline: %w", err)
		}
		count++
		if count%5_000 == 0 || count == 1 {
			if err := writer.Flush(); err != nil {
				return fmt.Errorf("flush: %w", err)
			}
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error after %d features: %w", count, err)
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("final flush: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("family %d has no valid geometry", familyID)
	}

	args := []string{
		"-o", outputPath,
		"--force",
		"--layer", fmt.Sprintf("family%d", familyID),
		"--read-parallel",
		"--no-feature-limit",
		"--no-tile-size-limit",
		"--drop-densest-as-needed",
		"--extend-zooms-if-still-dropping",
	}
	if minZoom > 0 {
		args = append(args, "--minimum-zoom", fmt.Sprintf("%d", minZoom))
	} else {
		args = append(args, "--minimum-zoom", "6")
	}
	if maxZoom > 0 {
		args = append(args, "--maximum-zoom", fmt.Sprintf("%d", maxZoom))
	} else {
		args = append(args, "--maximum-zoom", "16")
	}
	args = append(args, tmpPath)

	cmd := exec.CommandContext(ctx, "tippecanoe", args...)
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start tippecanoe: %w", err)
	}
	go readOutput("tippecanoe:stdout", stdoutPipe, progressFn)
	go readOutput("tippecanoe:stderr", stderrPipe, progressFn)
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("tippecanoe failed (exit code %d): %w", cmd.ProcessState.ExitCode(), err)
	}

	progressFn(BuildProgress{
		Status:       "done",
		Message:      fmt.Sprintf("PMTiles generated successfully from %d features", count),
		TilesWritten: count,
		OutputPath:   outputPath,
	})
	return nil
}

// Helper đọc output mà không dễ block
func readOutput(prefix string, r io.Reader, progressFn func(BuildProgress)) {
	scanner := bufio.NewScanner(r)
	// Tăng buffer nếu có dòng rất dài (geojson properties)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		slog.Info(fmt.Sprintf("[%s] %s", prefix, line))

		// Optional: parse progress nếu tippecanoe in ra % hoặc zoom
		// Ví dụ: nếu thấy "Zoom level" thì update CurrentZoom
		if progressFn != nil {
			// Bạn có thể cải tiến thêm parser ở đây
			progressFn(BuildProgress{
				Status:  "processing",
				Message: line,
			})
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		slog.Error(fmt.Sprintf("[%s] scanner error: %v", prefix, err))
	}
}
