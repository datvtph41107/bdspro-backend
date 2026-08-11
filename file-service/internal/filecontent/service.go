package filecontent

import (
	"bytes"
	_utils "common/utils"
	"context"
	"errors"
	"file/models"
	"file/repositories"
	"file/services"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	ErrInvalidPath  = errors.New("file path không hợp lệ")
	ErrNoActor      = errors.New("secure file cần profile hiện tại")
	ErrAccessDenied = errors.New("không có quyền đọc secure file")
)

const MaxContentBytes = 64 << 20

/**
 * UploadInput là dữ liệu tối thiểu cần để lưu một file từ service nội bộ.
 */
type UploadInput struct {
	Content     []byte
	Filename    string
	ContentType string
	Description string
	Secure      bool
	ProfileID   uint64
}

type SavedFile struct {
	ID          uint64
	Path        string
	Filename    string
	ContentType string
	Size        int64
}

type ReadFile struct {
	Content     []byte
	ContentType string
}

/**
 * Service xử lý lưu và đọc file cho gRPC nội bộ.
 *
 * Không nhận gin.Context và không đọc HTTP header.
 * repositories/StorageService hiện là bridge sang code file-service cũ;
 * khi migrate xong có thể thay bridge này mà không đổi gRPC contract.
 */
type Service struct {
	files    *repositories.FileRepository
	access   *repositories.AccessRepository
	storage  *services.StorageService
	xorKey   string
	document string
}

func NewService(
	files *repositories.FileRepository,
	access *repositories.AccessRepository,
	storage *services.StorageService,
	xorKey string,
) *Service {
	return &Service{
		files:    files,
		access:   access,
		storage:  storage,
		xorKey:   xorKey,
		document: "document",
	}
}

/**
 * UploadFile lưu file và trả path cùng format mà source hiện tại đang sử dụng.
 *
 * Nếu secure=true, profile hiện tại phải tồn tại để tạo quyền đọc file.
 */
func (s *Service) UploadFile(ctx context.Context, input UploadInput) (SavedFile, error) {
	if ctx == nil {
		return SavedFile{}, errors.New("context is nil")
	}
	if len(input.Content) == 0 {
		return SavedFile{}, errors.New("file content is empty")
	}
	if len(input.Content) > MaxContentBytes {
		return SavedFile{}, errors.New("file content exceeds 64 MiB")
	}
	input.Filename = strings.TrimSpace(input.Filename)
	if input.Filename == "" {
		return SavedFile{}, errors.New("filename is empty")
	}
	if input.Secure && input.ProfileID == 0 {
		return SavedFile{}, ErrNoActor
	}

	fileEntity := &models.FileEntity{Description: input.Description}
	if err := s.files.Save(fileEntity); err != nil {
		return SavedFile{}, fmt.Errorf("create file record: %w", err)
	}

	var accessID *uint64
	if input.Secure {
		profileID := input.ProfileID
		accessID = &profileID
	}

	fileInfo, err := s.storage.StoreContent(
		bytes.NewReader(input.Content),
		input.Filename,
		input.ContentType,
		int64(len(input.Content)),
		accessID,
		fileEntity.ID,
		s.document,
	)
	if err != nil {
		_ = s.files.DeleteByID(fileEntity.ID)
		return SavedFile{}, fmt.Errorf("store file content: %w", err)
	}

	encryptedPath := _utils.XorEncode(fileInfo.RelativePath, s.xorKey)
	fileEntity.SetFields(fileInfo)
	if input.Secure {
		fileEntity.Path = "s" + encryptedPath
	} else {
		fileEntity.Path = "p" + encryptedPath
	}

	if err := s.files.Save(fileEntity); err != nil {
		_ = os.Remove(fileInfo.AbsolutePath)
		_ = s.files.DeleteByID(fileEntity.ID)
		return SavedFile{}, fmt.Errorf("save file info: %w", err)
	}

	if input.Secure {
		if err := s.access.Create(&models.AccessEntity{
			FileID:   fileEntity.ID,
			AccessID: input.ProfileID,
		}); err != nil {
			_ = os.Remove(fileInfo.AbsolutePath)
			_ = s.files.DeleteByID(fileEntity.ID)
			return SavedFile{}, fmt.Errorf("save file access: %w", err)
		}
	}

	return SavedFile{
		ID:          fileEntity.ID,
		Path:        fileEntity.Path,
		Filename:    fileInfo.FileName,
		ContentType: fileInfo.ContentType,
		Size:        fileInfo.Size,
	}, nil
}

/**
 * ReadFile đọc nội dung bằng path đã lưu trong database.
 *
 * Đây là service-to-service call đã qua trusted gRPC metadata nên không dùng
 * HTTP signature. Secure file vẫn bắt buộc Actor và file_access tương ứng.
 */
func (s *Service) ReadFile(ctx context.Context, encodedPath string, profileID uint64) (ReadFile, error) {
	if ctx == nil {
		return ReadFile{}, errors.New("context is nil")
	}
	if len(encodedPath) < 2 {
		return ReadFile{}, ErrInvalidPath
	}
	prefix := encodedPath[0]
	if prefix != 'p' && prefix != 's' {
		return ReadFile{}, ErrInvalidPath
	}

	decodedPath, err := _utils.XorDecode(encodedPath[1:], s.xorKey)
	if err != nil {
		return ReadFile{}, ErrInvalidPath
	}
	cleanPath := filepath.Clean(decodedPath)
	if cleanPath == "." || strings.Contains(cleanPath, "..") {
		return ReadFile{}, ErrInvalidPath
	}
	root := filepath.Clean(s.storage.RootLocation)
	absolutePath := cleanPath
	if !filepath.IsAbs(absolutePath) && !strings.HasPrefix(cleanPath, root+string(os.PathSeparator)) {
		absolutePath = filepath.Join(root, cleanPath)
	}

	rel, err := filepath.Rel(root, absolutePath)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return ReadFile{}, ErrInvalidPath
	}
	expectedRoot := filepath.Clean(s.storage.RootPublicPath)
	if prefix == 's' {
		expectedRoot = filepath.Clean(s.storage.RootSecurePath)
	}
	expectedRel, err := filepath.Rel(expectedRoot, absolutePath)
	if err != nil || expectedRel == ".." || strings.HasPrefix(expectedRel, ".."+string(os.PathSeparator)) {
		return ReadFile{}, ErrInvalidPath
	}
	if prefix == 's' {
		if profileID == 0 {
			return ReadFile{}, ErrNoActor
		}
		fileID, err := fileIDFromPath(decodedPath)
		if err != nil {
			return ReadFile{}, ErrInvalidPath
		}
		allowed, err := s.access.ExistsAccess(int64(fileID), int64(profileID))
		if err != nil {
			return ReadFile{}, fmt.Errorf("check file access: %w", err)
		}
		if !allowed {
			return ReadFile{}, ErrAccessDenied
		}
	}

	content, err := os.ReadFile(absolutePath)
	if err != nil {
		return ReadFile{}, err
	}
	contentType := mime.TypeByExtension(strings.ToLower(filepath.Ext(absolutePath)))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return ReadFile{Content: content, ContentType: contentType}, nil
}

func fileIDFromPath(path string) (uint64, error) {
	base := filepath.Base(filepath.Clean(path))
	name := strings.TrimSuffix(base, filepath.Ext(base))
	fileID, err := strconv.ParseUint(name, 10, 64)
	if err != nil || fileID == 0 {
		return 0, ErrInvalidPath
	}
	return fileID, nil
}
