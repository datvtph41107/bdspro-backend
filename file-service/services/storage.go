package services

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"file/dto"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// StorageService xử lý lưu file vào thư mục hệ thống
type StorageService struct {
	RootLocation   string
	RootPublicPath string
	RootSecurePath string
}

// NewFileSystemStorageService khởi tạo service với các đường dẫn cụ thể
func NewStorageService() *StorageService {
	rootPath := "./files"
	publicPath := "./files/public"
	securePath := "./files/secure"

	createDirectories(rootPath)
	createDirectories(publicPath)
	createDirectories(securePath)

	return &StorageService{
		RootLocation:   rootPath,
		RootPublicPath: publicPath,
		RootSecurePath: securePath,
	}
}

// GetLocation trả về thư mục lưu file (public hoặc secure)
func (s *StorageService) GetLocation(secure bool) string {
	if secure {
		return s.RootSecurePath
	}
	return s.RootPublicPath
}

// Store lưu file vào thư mục phù hợp
func (s *StorageService) Store(file *multipart.FileHeader, accessID *uint64, fileID uint64, relativePaths ...string) (*dto.FileInfo, error) {
	if file == nil {
		return nil, errors.New("không tìm thấy file")
	}

	fileName := file.Filename
	filePath, ext, err := s.BuildFilePath(fileName, accessID, fileID, true, relativePaths...)
	if err != nil {
		return nil, err
	}

	// Tạo thư mục nếu chưa tồn tại
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return nil, fmt.Errorf("không thể tạo thư mục: %v", err)
	}

	// Lưu file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	// Tạo hash SHA-256
	hash, err := computeFileHash(filePath)
	if err != nil {
		return nil, err
	}

	return &dto.FileInfo{
		FileName:     fileName,
		RelativePath: strings.TrimPrefix(filePath, s.RootLocation+"/"),
		AbsolutePath: filePath,
		Extension:    ext,
		Size:         file.Size,
		Hash:         hash,
		ContentType:  file.Header.Get("Content-Type"),
	}, nil
}

// resolve xác định đường dẫn dựa vào quyền truy cập
func (s *StorageService) resolve(secure bool, relativePaths ...string) string {
	location := s.GetLocation(secure)
	for _, path := range relativePaths {
		if path != "" {
			location = filepath.Join(location, path)
		}
	}
	return location
}

// validatePath kiểm tra xem đường dẫn có hợp lệ không
func validatePath(path, parent string) {
	rel, err := filepath.Rel(parent, path)
	if err != nil || strings.Contains(rel, "..") {
		panic("path không hợp lệ")
	}
}

// computeFileHash tạo hash SHA-256 của file
func computeFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// createDirectories tạo thư mục nếu chưa tồn tại
func createDirectories(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
}

// DeleteAll xóa tất cả file trong thư mục public
func (s *StorageService) DeleteAll() error {
	return os.RemoveAll(s.RootPublicPath)
}

func (s *StorageService) BuildFilePath(fileName string, accessID *uint64, fileID uint64, withExt bool, relativePaths ...string) (string, string, error) {
	ext := filepath.Ext(fileName)
	if fileName == "" {
		return "", ext, errors.New("không có tên file")
	}

	// Làm sạch đường dẫn
	fileName = filepath.Base(fileName)

	// Xác định thư mục lưu trữ
	secure := accessID != nil
	location := s.resolve(secure, relativePaths...)
	today := time.Now().Format("2006-01-02")         // YYYY-MM-DD
	uuid := fmt.Sprintf("%d", time.Now().UnixNano()) // UUID tạm thời

	// Đường dẫn đầy đủ
	filePath := filepath.Join(location, today, uuid, fmt.Sprintf("%d", fileID))
	if withExt {
		filePath = fmt.Sprintf("%s%s", filePath, ext)
	}
	validatePath(filePath, s.RootLocation)

	return filePath, ext, nil
}

func (s *StorageService) StoreVideo(c *gin.Context, file *multipart.FileHeader, fileName string, accessID *uint64, fileID uint64) (*dto.FileInfo, error) {
	path, ext, err := s.BuildFilePath(fileName, accessID, fileID, false, "video")
	if err != nil {
		return nil, err
	}

	// outputDir := filepath.Join("public", "videos", strconv.FormatUint(fileID, 10))
	m3u8IndexPath := filepath.Join(path, "playlist.m3u8")
	segmentPattern := filepath.Join(path, "%08d.ts")

	// Tạo thư mục
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create output directory"})
		return nil, errors.New("Failed to create output directory")
	}

	// Lưu file tạm vào /tmp
	tempInput := filepath.Join(os.TempDir(), fmt.Sprintf("temp_%s.mp4", fileID))
	if err := c.SaveUploadedFile(file, tempInput); err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save temporary input"})
		return nil, errors.New("Failed to save temporary input")
	}

	// // FFmpeg command để chunk video thành HLS .mp4
	// cmd := exec.Command("ffmpeg",
	// 	"-i", tempInput,
	// 	"-c:v", "libx264",
	// 	"-c:a", "aac",
	// 	"-preset", "veryfast",
	// 	"-movflags", "frag_keyframe+empty_moov",
	// 	"-f", "hls",
	// 	"-hls_time", "10",
	// 	"-hls_list_size", "0",
	// 	"-hls_segment_filename", segmentPattern,
	// 	m3u8IndexPath,
	// )

	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	checkCmd := exec.Command("ffmpeg", "-i", tempInput)

	var stderr bytes.Buffer
	checkCmd.Stderr = &stderr
	_ = checkCmd.Run() // bỏ qua lỗi vì ffmpeg -i không cần thành công

	output := strings.ToLower(stderr.String())

	videoOK := strings.Contains(output, "video: h264")
	audioOK := strings.Contains(output, "audio: aac")

	fmt.Println("📋 Kiểm tra codec:")
	fmt.Println(" - Video là h264:", videoOK)
	fmt.Println(" - Audio là aac:", audioOK)

	// Bước 2: Dựng command chuyển sang HLS
	var cmd *exec.Cmd
	// numCPU := runtime.NumCPU()
	// threads := strconv.Itoa(numCPU)

	if videoOK && audioOK {
		// ✅ Nếu đúng codec, không cần encode lại
		cmd = exec.Command("ffmpeg",
			"-i", tempInput,
			"-c", "copy",
			"-f", "hls",
			"-hls_time", "10",
			"-hls_list_size", "0",
			"-hls_segment_filename", segmentPattern, // ví dụ: "segment_%03d.ts"
			m3u8IndexPath,
		)
		fmt.Println("🚀 Dùng -c copy (không encode lại).")
	} else {
		// giảm tải cho server không encode lại video. phần này cho fe làm
		// ❌ Nếu codec sai, encode lại
		// cmd = exec.Command("ffmpeg",
		// 	"-i", tempInput,
		// 	"-threads", threads, // Thêm dòng này
		// 	"-c:v", "libx264",
		// 	"-c:a", "aac",
		// 	"-preset", "veryfast",
		// 	"-f", "hls",
		// 	"-hls_time", "10",
		// 	"-hls_list_size", "0",
		// 	"-hls_segment_filename", segmentPattern, // ví dụ: "segment_%03d.ts"
		// 	m3u8IndexPath,
		// )
		// fmt.Println("⚠️ Encode lại do codec không đúng.")
		return nil, errors.New("codec không đúng, yêu cầu video h264 và audio aac")
	}

	if err := cmd.Run(); err != nil {
		os.Remove(tempInput)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("FFmpeg failed: %v", err)})
		return nil, errors.New("FFmpeg failed")
	}

	hash, err := computeFileHash(m3u8IndexPath)
	if err != nil {
		return nil, err
	}

	// Tạo đường dẫn thumbnail
	// thumbnailPath := filepath.Join(path, "thumbnail.jpg")
	thumbnailPath, _, err := s.BuildFilePath(fileName, accessID, fileID, false, "document")
	if err != nil {
		fmt.Printf("🔴 Lỗi khi tạo đường dẫn thumbnail: %v", err)
	}

	// Lệnh ffmpeg: lấy 1 frame tại giây thứ 1
	cmd = exec.Command("ffmpeg",
		"-i", tempInput, // input video file
		"-ss", "00:00:01", // lấy frame tại 1 giây
		"-vframes", "1", // chỉ lấy 1 frame
		"-vf", "scale=320:-1", // resize thumbnail chiều rộng 320px (chiều cao tự động theo tỉ lệ)
		"-q:v", "2", // chất lượng hình ảnh (1 tốt nhất, 31 tệ nhất)
		thumbnailPath,
	)

	// Thực thi lệnh
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to generate thumbnail: %v - %s", err, stderr.String())
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Xóa file gốc
	os.Remove(tempInput)

	return &dto.FileInfo{
		FileName:      fileName,
		RelativePath:  strings.TrimPrefix(path, s.RootLocation+"/"),
		AbsolutePath:  path,
		ThumbnailPath: thumbnailPath,
		Extension:     ext,
		Size:          file.Size,
		Hash:          hash,
		ContentType:   file.Header.Get("Content-Type"),
	}, nil
}

func ParseVideoPath(path string) (*dto.FileInfo, error) {
	// Chuẩn hóa đường dẫn
	clean := filepath.ToSlash(filepath.Clean(path)) // Convert `\` to `/` nếu trên Windows

	// Tách các phần
	parts := strings.Split(clean, "/")
	if len(parts) < 6 {
		return nil, fmt.Errorf("path không hợp lệ: %s", path)
	}

	// Lấy fileName và tách phần mở rộng
	fileName := parts[len(parts)-1]
	fileID := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	ext := strings.TrimPrefix(filepath.Ext(fileName), ".")

	info := &dto.FileInfo{
		AbsolutePath: path,
		Scope:        parts[1], // public
		Type:         parts[2], // video
		Date:         parts[3], // 2025-07-14
		UUID:         parts[4], // 1752485034177081000
		FileID:       fileID,   // 619
		Ext:          ext,      // mp4
	}
	return info, nil
}
