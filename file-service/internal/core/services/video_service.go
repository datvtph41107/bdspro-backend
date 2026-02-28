package core_services

import (
	"bytes"
	"context"
	"file/internal/config"
	"file/internal/ports"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type videoService struct {
	storage      ports.StorageService
	cacheRepo    ports.CacheRepository
	fileRepo     ports.FileRepository
	config       *config.Config
	videoWorkers int
	jobQueue     chan videoJob
}

type videoJob struct {
	FileID   string
	FilePath string
	Options  ports.VideoProcessOptions
	Result   chan error
}

func NewVideoService(
	storage ports.StorageService,
	cacheRepo ports.CacheRepository,
	fileRepo ports.FileRepository,
	config *config.Config,
) ports.VideoService {
	workers := 4 // Default number of workers
	if config.Queue.WorkerCount > 0 {
		workers = config.Queue.WorkerCount / 2 // Use half of workers for video processing
	}

	svc := &videoService{
		storage:      storage,
		cacheRepo:    cacheRepo,
		fileRepo:     fileRepo,
		config:       config,
		videoWorkers: workers,
		jobQueue:     make(chan videoJob, 100),
	}

	// Start video processing workers
	for i := 0; i < workers; i++ {
		go svc.videoWorker(i)
	}

	return svc
}

func (s *videoService) ProcessVideo(ctx context.Context, fileID string, options ports.VideoProcessOptions) error {
	// Get file from repository
	file, err := s.fileRepo.FindByID(ctx, parseFileID(fileID))
	if err != nil {
		return fmt.Errorf("failed to find file: %w", err)
	}

	// Create job
	job := videoJob{
		FileID:   fileID,
		FilePath: file.StoragePath,
		Options:  options,
		Result:   make(chan error, 1),
	}

	// Send to queue
	select {
	case s.jobQueue <- job:
		// Job queued successfully
		return <-job.Result
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("video processing queue is full")
	}
}

func (s *videoService) videoWorker(workerID int) {
	log.Printf("Video worker %d started", workerID)

	for job := range s.jobQueue {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)

		err := s.processVideoJob(ctx, job)
		if err != nil {
			log.Printf("Video worker %d failed to process %s: %v", workerID, job.FileID, err)
		}

		job.Result <- err
		cancel()
	}
}

func (s *videoService) processVideoJob(ctx context.Context, job videoJob) error {
	startTime := time.Now()
	log.Printf("Processing video %s", job.FileID)

	// 1. Download video to temp file
	tempInput, err := s.downloadToTempFile(ctx, job.FilePath)
	if err != nil {
		return fmt.Errorf("failed to download video: %w", err)
	}
	defer os.Remove(tempInput)

	// 2. Generate thumbnail if requested
	if job.Options.GenerateThumbnail {
		thumbnailPath, err := s.generateThumbnail(ctx, tempInput, job.FileID, ports.ThumbnailOptions{
			TimeOffset: "00:00:01",
			Width:      320,
			Height:     240,
			Quality:    80,
		})
		if err != nil {
			log.Printf("Failed to generate thumbnail: %v", err)
		} else {
			// Update file record with thumbnail path
			fileID := parseFileID(job.FileID)
			if err := s.fileRepo.UpdateThumbnail(ctx, fileID, thumbnailPath); err != nil {
				log.Printf("Failed to update thumbnail path: %v", err)
			}
		}
	}

	// 3. Create HLS streams if requested
	if job.Options.CreateHLS {
		hlsPath, err := s.createHLSStream(ctx, tempInput, job.FileID, job.Options.Qualities)
		if err != nil {
			return fmt.Errorf("failed to create HLS stream: %w", err)
		}
		log.Printf("HLS stream created at: %s", hlsPath)
	}

	// 4. Calculate video duration
	duration, err := s.getVideoDuration(ctx, tempInput)
	if err == nil {
		// Update file record with duration
		fileID := parseFileID(job.FileID)
		file, err := s.fileRepo.FindByID(ctx, fileID)
		if err == nil {
			file.Duration = duration
			s.fileRepo.Update(ctx, file)
		}
	}

	log.Printf("Video %s processed in %v", job.FileID, time.Since(startTime))
	return nil
}

func (s *videoService) GenerateThumbnail(ctx context.Context, fileID string, options ports.ThumbnailOptions) (string, error) {
	// Get file from repository
	file, err := s.fileRepo.FindByID(ctx, parseFileID(fileID))
	if err != nil {
		return "", fmt.Errorf("failed to find file: %w", err)
	}

	// Download to temp file
	tempInput, err := s.downloadToTempFile(ctx, file.StoragePath)
	if err != nil {
		return "", fmt.Errorf("failed to download video: %w", err)
	}
	defer os.Remove(tempInput)

	// Generate thumbnail
	return s.generateThumbnail(ctx, tempInput, fileID, options)
}

func (s *videoService) GetVideoStream(ctx context.Context, fileID string, quality string) (io.Reader, *ports.VideoStreamInfo, error) {
	// Get file from repository
	file, err := s.fileRepo.FindByID(ctx, parseFileID(fileID))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find file: %w", err)
	}

	// For HLS videos, return the playlist
	if strings.Contains(file.StoragePath, ".m3u8") {
		return s.getHLSStream(ctx, file.StoragePath, quality)
	}

	// For regular videos, return the file
	reader, size, err := s.storage.Retrieve(ctx, file.StoragePath)
	if err != nil {
		return nil, nil, err
	}

	return reader, &ports.VideoStreamInfo{
		ContentType:   file.MimeType,
		ContentLength: size,
		IsStream:      false,
	}, nil
}

func (s *videoService) TranscodeVideo(ctx context.Context, fileID string, profiles []ports.VideoProfile) error {
	// Implement video transcoding to multiple profiles
	// This would create different quality versions of the video
	return nil
}

// Helper methods
func (s *videoService) downloadToTempFile(ctx context.Context, storagePath string) (string, error) {
	// Create temp file
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("video_%s.mp4", uuid.New().String()))

	// Download from storage
	reader, _, err := s.storage.Retrieve(ctx, storagePath)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	// Create temp file
	file, err := os.Create(tempFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Copy data
	if _, err := io.Copy(file, reader); err != nil {
		os.Remove(tempFile)
		return "", err
	}

	return tempFile, nil
}

func (s *videoService) generateThumbnail(ctx context.Context, inputPath, fileID string, options ports.ThumbnailOptions) (string, error) {
	// Generate unique thumbnail path
	thumbnailKey := fmt.Sprintf("thumbnails/%s/%s.jpg",
		time.Now().Format("2006/01/02"),
		strings.ReplaceAll(fileID, "/", "_"))

	// Generate thumbnail with ffmpeg
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", inputPath,
		"-ss", options.TimeOffset,
		"-vframes", "1",
		"-vf", fmt.Sprintf("scale=%d:%d", options.Width, options.Height),
		"-q:v", strconv.Itoa(options.Quality),
		"-f", "image2pipe",
		"-",
	)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ffmpeg failed: %w", err)
	}

	// Store thumbnail
	metadata := map[string]string{
		"Content-Type": "image/jpeg",
		"File-ID":      fileID,
		"Generated-At": time.Now().Format(time.RFC3339),
	}

	if err := s.storage.Store(ctx, thumbnailKey, bytes.NewReader(output),
		int64(len(output)), metadata); err != nil {
		return "", fmt.Errorf("failed to store thumbnail: %w", err)
	}

	return thumbnailKey, nil
}

func (s *videoService) createHLSStream(ctx context.Context, inputPath, fileID string, qualities []string) (string, error) {
	// Create HLS output directory
	outputDir := fmt.Sprintf("hls/%s/%s",
		time.Now().Format("2006/01/02"),
		strings.ReplaceAll(fileID, "/", "_"))

	// Create master playlist
	masterPlaylist := fmt.Sprintf("%s/master.m3u8", outputDir)

	// Process each quality
	for _, quality := range qualities {
		var width, height int
		var bitrate string

		switch quality {
		case "360p":
			width, height, bitrate = 640, 360, "800k"
		case "720p":
			width, height, bitrate = 1280, 720, "2500k"
		case "1080p":
			width, height, bitrate = 1920, 1080, "5000k"
		default:
			continue
		}

		// Create quality-specific stream
		segmentPattern := fmt.Sprintf("%s/%s/segment_%%03d.ts", outputDir, quality)
		playlistPath := fmt.Sprintf("%s/%s/playlist.m3u8", outputDir, quality)

		cmd := exec.CommandContext(ctx, "ffmpeg",
			"-i", inputPath,
			"-vf", fmt.Sprintf("scale=%d:%d", width, height),
			"-c:v", "libx264",
			"-b:v", bitrate,
			"-c:a", "aac",
			"-b:a", "128k",
			"-f", "hls",
			"-hls_time", "10",
			"-hls_list_size", "0",
			"-hls_segment_filename", segmentPattern,
			playlistPath,
		)

		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to create %s stream: %w", quality, err)
		}
	}

	// Generate master playlist
	masterContent := "#EXTM3U\n"
	for _, quality := range qualities {
		masterContent += fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=800000,RESOLUTION=640x360\n%s/%s/playlist.m3u8\n", outputDir, quality)
	}

	// Store master playlist
	metadata := map[string]string{
		"Content-Type": "application/vnd.apple.mpegurl",
		"File-ID":      fileID,
	}

	if err := s.storage.Store(ctx, masterPlaylist, strings.NewReader(masterContent),
		int64(len(masterContent)), metadata); err != nil {
		return "", fmt.Errorf("failed to store master playlist: %w", err)
	}

	return masterPlaylist, nil
}

func (s *videoService) getHLSStream(ctx context.Context, storagePath, quality string) (io.Reader, *ports.VideoStreamInfo, error) {
	// For HLS, we need to return the appropriate playlist
	var streamPath string
	if quality == "" {
		// Return master playlist
		streamPath = storagePath
	} else {
		// Return quality-specific playlist
		streamPath = strings.Replace(storagePath, "master.m3u8",
			fmt.Sprintf("%s/playlist.m3u8", quality), 1)
	}

	reader, size, err := s.storage.Retrieve(ctx, streamPath)
	if err != nil {
		return nil, nil, err
	}

	return reader, &ports.VideoStreamInfo{
		ContentType:   "application/vnd.apple.mpegurl",
		ContentLength: size,
		IsStream:      true,
	}, nil
}

func (s *videoService) getVideoDuration(ctx context.Context, inputPath string) (float64, error) {
	cmd := exec.CommandContext(ctx, "ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputPath,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var duration float64
	if _, err := fmt.Sscanf(string(output), "%f", &duration); err != nil {
		return 0, err
	}

	return duration, nil
}

func parseFileID(fileID string) uint64 {
	var id uint64
	fmt.Sscanf(fileID, "%d", &id)
	return id
}
