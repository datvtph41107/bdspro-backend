package filemedia

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"file/dto"
)

type videoStorage interface {
	BuildFilePath(
		fileName string,
		accessID *uint64,
		fileID uint64,
		withExt bool,
		relativePaths ...string,
	) (string, string, error)

	StoredPath(path string) (string, error)
	ResolveStoredPath(path string) (string, error)
	RemovePhysicalPath(path string) error
	HashFile(path string) (string, error)
}

type VideoInput struct {
	SourcePath  string
	FileName    string
	ContentType string
	Size        int64
	AccessID    *uint64
	FileID      uint64
}

// Processor owns video-specific filesystem/ffmpeg effects. SourcePath is
// borrowed: the caller owns creation and cleanup of the temporary input.
type Processor struct {
	storage videoStorage
}

func NewProcessor(storage videoStorage) *Processor {
	return &Processor{storage: storage}
}

func (p *Processor) Process(input VideoInput) (_ *dto.FileInfo, err error) {
	if p == nil || p.storage == nil {
		return nil, errors.New("video storage is unavailable")
	}
	if strings.TrimSpace(input.SourcePath) == "" {
		return nil, errors.New("video source path is empty")
	}

	videoPath, ext, err := p.storage.BuildFilePath(
		input.FileName,
		input.AccessID,
		input.FileID,
		false,
		"video",
	)
	if err != nil {
		return nil, err
	}

	var thumbnailPhysicalPath string
	committed := false
	defer func() {
		if committed {
			return
		}
		_ = p.storage.RemovePhysicalPath(videoPath)
		if thumbnailPhysicalPath != "" {
			_ = p.storage.RemovePhysicalPath(thumbnailPhysicalPath)
		}
	}()

	if err := os.MkdirAll(videoPath, 0o755); err != nil {
		return nil, fmt.Errorf("create video output directory: %w", err)
	}

	if _, err := os.Stat(input.SourcePath); err != nil {
		return nil, fmt.Errorf("stat video source: %w", err)
	}

	probe := exec.Command("ffmpeg", "-i", input.SourcePath)
	var probeStderr bytes.Buffer
	probe.Stderr = &probeStderr
	// ffmpeg -i commonly exits non-zero after printing stream metadata. The
	// metadata itself is the compatibility theorem used here.
	_ = probe.Run()

	videoOK, audioOK := supportedCodecs(probeStderr.String())
	if !videoOK || !audioOK {
		return nil, errors.New("unsupported video codec: require h264 video and aac audio")
	}

	playlistPath := filepath.Join(videoPath, "playlist.m3u8")
	segmentPattern := filepath.Join(videoPath, "%08d.ts")
	transcode := exec.Command(
		"ffmpeg",
		"-i", input.SourcePath,
		"-c", "copy",
		"-f", "hls",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-hls_segment_filename", segmentPattern,
		playlistPath,
	)
	var transcodeStderr bytes.Buffer
	transcode.Stderr = &transcodeStderr
	if err := transcode.Run(); err != nil {
		return nil, fmt.Errorf("create HLS media: %w: %s", err, boundedOutput(transcodeStderr.String()))
	}

	hash, err := p.storage.HashFile(playlistPath)
	if err != nil {
		return nil, fmt.Errorf("hash HLS playlist: %w", err)
	}

	thumbnailBase, _, err := p.storage.BuildFilePath(
		input.FileName,
		input.AccessID,
		input.FileID,
		false,
		"document",
	)
	if err != nil {
		return nil, fmt.Errorf("build thumbnail path: %w", err)
	}
	thumbnailPhysicalPath = thumbnailBase + ".jpg"

	if err := os.MkdirAll(filepath.Dir(thumbnailPhysicalPath), 0o755); err != nil {
		return nil, fmt.Errorf("create thumbnail directory: %w", err)
	}

	thumbnail := exec.Command(
		"ffmpeg",
		"-i", input.SourcePath,
		"-ss", "00:00:01",
		"-vframes", "1",
		"-vf", "scale=320:-1",
		"-q:v", "2",
		"-f", "image2",
		"-y",
		thumbnailPhysicalPath,
	)
	var thumbnailStderr bytes.Buffer
	thumbnail.Stderr = &thumbnailStderr
	if err := thumbnail.Run(); err != nil {
		return nil, fmt.Errorf("generate thumbnail: %w: %s", err, boundedOutput(thumbnailStderr.String()))
	}

	storedVideoPath, err := p.storage.StoredPath(videoPath)
	if err != nil {
		return nil, fmt.Errorf("project video storage path: %w", err)
	}
	storedThumbnailPath, err := p.storage.StoredPath(thumbnailPhysicalPath)
	if err != nil {
		return nil, fmt.Errorf("project thumbnail storage path: %w", err)
	}

	committed = true
	return &dto.FileInfo{
		FileName:      input.FileName,
		RelativePath:  storedVideoPath,
		AbsolutePath:  videoPath,
		ThumbnailPath: storedThumbnailPath,
		Extension:     ext,
		Size:          input.Size,
		Hash:          hash,
		ContentType:   input.ContentType,
	}, nil
}

// Cleanup compensates a successfully processed media effect when the durable
// metadata finalization fails. Both outputs stay under the configured File
// storage authority.
func (p *Processor) Cleanup(info *dto.FileInfo) error {
	if p == nil || p.storage == nil || info == nil {
		return nil
	}

	var cleanupErr error
	if strings.TrimSpace(info.AbsolutePath) != "" {
		cleanupErr = errors.Join(cleanupErr, p.storage.RemovePhysicalPath(info.AbsolutePath))
	}
	if strings.TrimSpace(info.ThumbnailPath) != "" {
		thumbnailPath, err := p.storage.ResolveStoredPath(info.ThumbnailPath)
		if err != nil {
			cleanupErr = errors.Join(cleanupErr, fmt.Errorf("resolve thumbnail for cleanup: %w", err))
		} else {
			cleanupErr = errors.Join(cleanupErr, p.storage.RemovePhysicalPath(thumbnailPath))
		}
	}
	return cleanupErr
}

func supportedCodecs(output string) (videoOK bool, audioOK bool) {
	output = strings.ToLower(output)
	videoOK = strings.Contains(output, "video") &&
		(strings.Contains(output, "h264") || strings.Contains(output, "avc"))
	audioOK = strings.Contains(output, "audio") &&
		(strings.Contains(output, "aac") || strings.Contains(output, "mp4a"))
	return videoOK, audioOK
}

func boundedOutput(output string) string {
	output = strings.TrimSpace(output)
	const max = 1000
	if len(output) <= max {
		return output
	}
	return output[:max] + "..."
}
