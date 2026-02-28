package core_services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"file/internal/config"
	"file/internal/core/domain"
	"file/internal/pkg/utils"
	"file/internal/ports"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type fileService struct {
	fileRepo     ports.FileRepository
	accessRepo   ports.FileAccessRepository
	sessionRepo  ports.UploadSessionRepository
	cacheRepo    ports.CacheRepository
	storage      ports.StorageService
	config       *config.Config
	videoService ports.VideoService
}

func NewFileService(
	fileRepo ports.FileRepository,
	accessRepo ports.FileAccessRepository,
	sessionRepo ports.UploadSessionRepository,
	cacheRepo ports.CacheRepository,
	storage ports.StorageService,
	config *config.Config,
	videoService ports.VideoService,
) ports.FileService {
	return &fileService{
		fileRepo:     fileRepo,
		accessRepo:   accessRepo,
		sessionRepo:  sessionRepo,
		cacheRepo:    cacheRepo,
		storage:      storage,
		config:       config,
		videoService: videoService,
	}
}

func (s *fileService) Upload(ctx context.Context, req ports.UploadRequest) (*ports.UploadResponse, error) {
	// Validate file
	if err := s.validateFile(req.FileHeader, s.config.Storage.MaxUploadSize); err != nil {
		return nil, fmt.Errorf("file validation failed: %w", err)
	}

	// Read file and calculate hash
	hash, size, err := s.calculateFileHash(req.File)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate file hash: %w", err)
	}

	// Check if file already exists
	existingFile, err := s.fileRepo.FindByHash(ctx, hash)
	if err == nil && existingFile != nil {
		// File already exists, return existing file info
		return s.createUploadResponse(existingFile), nil
	}

	// Generate file metadata
	fileName := req.OriginalName
	if fileName == "" {
		fileName = req.FileHeader.Filename
	}

	extension := strings.ToLower(filepath.Ext(fileName))
	mimeType := req.FileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = utils.GetMimeType(extension)
	}

	// Generate storage key
	storageKey := s.generateStorageKey(fileName, hash)

	// Store file
	metadata := map[string]string{
		"Content-Type":  mimeType,
		"Original-Name": fileName,
		"Service-ID":    req.ServiceID,
		"Uploaded-By":   fmt.Sprintf("%d", req.UserID),
	}

	if err := s.storage.Store(ctx, storageKey, req.File, size, metadata); err != nil {
		return nil, fmt.Errorf("failed to store file: %w", err)
	}

	// Create file record
	file := &domain.File{
		UUID:          utils.GenerateUUID(),
		OriginalName:  fileName,
		StoragePath:   storageKey,
		PublicPath:    s.generatePublicPath(storageKey, req.IsPublic),
		Size:          size,
		MimeType:      mimeType,
		Extension:     extension,
		Hash:          hash,
		Metadata:      domain.JSONMap(req.Metadata),
		Status:        domain.FileStatusActive,
		UploadedBy:    req.UserID,
		ServiceID:     req.ServiceID,
		ReferenceType: req.ReferenceType,
		ReferenceID:   req.ReferenceID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Set width/height for images
	if strings.HasPrefix(mimeType, "image/") {
		// Extract image dimensions
		if dims, err := s.getImageDimensions(req.File); err == nil {
			file.Width = dims.Width
			file.Height = dims.Height
		}
		// Reset file reader
		req.File.Seek(0, 0)
	}

	// Save to database
	if err := s.fileRepo.Create(ctx, file); err != nil {
		// Clean up storage if DB save fails
		s.storage.Delete(ctx, storageKey)
		return nil, fmt.Errorf("failed to save file record: %w", err)
	}

	// Cache file info
	s.cacheFileInfo(ctx, file)

	// Process video if needed
	if strings.HasPrefix(mimeType, "video/") && s.videoService != nil {
		go s.processVideoBackground(ctx, file.ID, file.StoragePath)
	}

	return s.createUploadResponse(file), nil
}

func (s *fileService) UploadChunk(ctx context.Context, req ports.UploadChunkRequest) (*ports.UploadChunkResponse, error) {
	// Get or create upload session
	session, err := s.sessionRepo.GetSession(ctx, req.SessionID)
	if err != nil {
		// Create new session
		session = &domain.UploadSession{
			ID:           req.SessionID,
			TotalChunks:  req.TotalChunks,
			ChunkSize:    int64(len(req.ChunkData)),
			TotalSize:    int64(len(req.ChunkData) * req.TotalChunks),
			UploadedSize: int64(len(req.ChunkData)),
			Status:       domain.FileStatusUploading,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			ExpiresAt:    time.Now().Add(24 * time.Hour),
		}

		if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
			return nil, fmt.Errorf("failed to create upload session: %w", err)
		}
	} else {
		// Update existing session
		session.UploadedSize += int64(len(req.ChunkData))
		session.UpdatedAt = time.Now()

		if err := s.sessionRepo.UpdateSession(ctx, session); err != nil {
			return nil, fmt.Errorf("failed to update upload session: %w", err)
		}
	}

	// Store chunk in temporary storage
	chunkKey := fmt.Sprintf("chunks/%s/%d", req.SessionID, req.ChunkNumber)
	metadata := map[string]string{
		"Session-ID":   req.SessionID,
		"Chunk-Number": fmt.Sprintf("%d", req.ChunkNumber),
		"Total-Chunks": fmt.Sprintf("%d", req.TotalChunks),
	}

	if err := s.storage.Store(ctx, chunkKey, bytes.NewReader(req.ChunkData),
		int64(len(req.ChunkData)), metadata); err != nil {
		return nil, fmt.Errorf("failed to store chunk: %w", err)
	}

	// Check if all chunks are uploaded
	isComplete := session.UploadedSize >= session.TotalSize

	response := &ports.UploadChunkResponse{
		SessionID:     req.SessionID,
		ChunkReceived: req.ChunkNumber,
		TotalChunks:   req.TotalChunks,
		IsComplete:    isComplete,
	}

	// If complete, assemble file
	if isComplete {
		go s.assembleFileFromChunks(ctx, req.SessionID, req.FileHash)
		response.FileID = "processing" // Temporary, will be updated
	}

	return response, nil
}

func (s *fileService) Download(ctx context.Context, req ports.DownloadRequest) (io.Reader, *ports.FileInfo, error) {
	// Get file from cache or DB
	file, err := s.getFileWithCache(ctx, req.FileID)
	if err != nil {
		return nil, nil, fmt.Errorf("file not found: %w", err)
	}

	// Check access
	if file.UploadedBy != req.UserID && !s.isPublicFile(file) {
		hasAccess, err := s.accessRepo.HasAccess(ctx, file.ID, req.UserID, domain.PermissionRead)
		if err != nil || !hasAccess {
			return nil, nil, fmt.Errorf("access denied")
		}
	}

	// Validate signature for public files
	if s.isPublicFile(file) && req.Signature != "" {
		if !s.validateSignature(ctx, req.Signature, file.ID, req.UserID) {
			return nil, nil, fmt.Errorf("invalid signature")
		}
	}

	// Get file from storage
	reader, size, err := s.storage.Retrieve(ctx, file.StoragePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve file: %w", err)
	}

	// Create file info response
	fileInfo := &ports.FileInfo{
		ID:           fmt.Sprintf("%d", file.ID),
		OriginalName: file.OriginalName,
		Size:         size,
		MimeType:     file.MimeType,
		Extension:    file.Extension,
		PublicURL:    s.getPublicURL(file),
		Width:        file.Width,
		Height:       file.Height,
		Duration:     file.Duration,
		Metadata:     map[string]interface{}(file.Metadata),
		Status:       file.Status,
		CreatedAt:    file.CreatedAt,
		UpdatedAt:    file.UpdatedAt,
	}

	return reader, fileInfo, nil
}

func (s *fileService) GetFileInfo(ctx context.Context, fileID string) (*ports.FileInfo, error) {
	file, err := s.getFileWithCache(ctx, fileID)
	if err != nil {
		return nil, err
	}

	return &ports.FileInfo{
		ID:           fmt.Sprintf("%d", file.ID),
		OriginalName: file.OriginalName,
		Size:         file.Size,
		MimeType:     file.MimeType,
		Extension:    file.Extension,
		PublicURL:    s.getPublicURL(file),
		ThumbnailURL: s.getThumbnailURL(file),
		Width:        file.Width,
		Height:       file.Height,
		Duration:     file.Duration,
		Metadata:     map[string]interface{}(file.Metadata),
		Status:       file.Status,
		CreatedAt:    file.CreatedAt,
		UpdatedAt:    file.UpdatedAt,
	}, nil
}

func (s *fileService) DeleteFile(ctx context.Context, req ports.DeleteRequest) error {
	file, err := s.getFileWithCache(ctx, req.FileID)
	if err != nil {
		return err
	}

	// Check permissions
	if file.UploadedBy != req.UserID {
		hasAccess, err := s.accessRepo.HasAccess(ctx, file.ID, req.UserID, domain.PermissionDelete)
		if err != nil || !hasAccess {
			return fmt.Errorf("permission denied")
		}
	}

	// Soft delete if not forced
	if !req.Force {
		file.Status = domain.FileStatusDeleted
		file.DeletedAt = time.Now()
		if err := s.fileRepo.Update(ctx, file); err != nil {
			return fmt.Errorf("failed to soft delete file: %w", err)
		}

		// Update cache
		s.cacheFileInfo(ctx, file)
		return nil
	}

	// Hard delete
	if err := s.storage.Delete(ctx, file.StoragePath); err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Delete thumbnail if exists
	if file.ThumbnailPath != "" {
		s.storage.Delete(ctx, file.ThumbnailPath)
	}

	// Delete from database
	if err := s.fileRepo.Delete(ctx, file.ID); err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	// Delete from cache
	s.cacheRepo.Delete(ctx, s.cacheRepo.FileKey(req.FileID))

	return nil
}

func (s *fileService) GenerateAccessURL(ctx context.Context, req ports.AccessURLRequest) (string, error) {
	file, err := s.getFileWithCache(ctx, req.FileID)
	if err != nil {
		return "", err
	}

	// Check permissions
	hasAccess, err := s.accessRepo.HasAccess(ctx, file.ID, req.UserID, req.Permission)
	if err != nil || !hasAccess {
		return "", fmt.Errorf("permission denied")
	}

	// Generate signature
	signature := s.generateSignature(file.ID, req.UserID, req.Expiry)

	// Cache signature
	s.cacheRepo.Set(ctx, s.cacheRepo.SignatureKey(signature),
		[]byte(fmt.Sprintf("%d:%d", file.ID, req.UserID)), req.Expiry)

	// Generate access URL
	baseURL := s.config.Storage.CDNBaseURL
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://localhost:%s", s.config.Server.Port)
	}

	accessURL := fmt.Sprintf("%s/v1/file/%s?signature=%s", baseURL, file.PublicPath, signature)

	return accessURL, nil
}

func (s *fileService) ListFiles(ctx context.Context, req ports.ListRequest) (*ports.ListResponse, error) {
	filter := ports.FileFilter{
		ServiceID:     req.ServiceID,
		ReferenceType: req.ReferenceType,
		ReferenceID:   req.ReferenceID,
		MimeType:      req.MimeType,
	}

	pagination := ports.Pagination{
		Page:     req.Page,
		PageSize: req.PageSize,
		SortBy:   req.SortBy,
		SortDir:  req.SortDir,
	}

	files, total, err := s.fileRepo.List(ctx, filter, pagination)
	if err != nil {
		return nil, err
	}

	// Convert to FileInfo
	fileInfos := make([]*ports.FileInfo, len(files))
	for i, file := range files {
		fileInfos[i] = &ports.FileInfo{
			ID:           fmt.Sprintf("%d", file.ID),
			OriginalName: file.OriginalName,
			Size:         file.Size,
			MimeType:     file.MimeType,
			Extension:    file.Extension,
			PublicURL:    s.getPublicURL(file),
			ThumbnailURL: s.getThumbnailURL(file),
			Width:        file.Width,
			Height:       file.Height,
			Duration:     file.Duration,
			Metadata:     map[string]interface{}(file.Metadata),
			Status:       file.Status,
			CreatedAt:    file.CreatedAt,
			UpdatedAt:    file.UpdatedAt,
		}
	}

	totalPages := (total + int64(req.PageSize) - 1) / int64(req.PageSize)

	return &ports.ListResponse{
		Files:      fileInfos,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: int(totalPages),
	}, nil
}

func (s *fileService) BatchUpload(ctx context.Context, files []ports.UploadRequest) ([]*ports.UploadResponse, error) {
	// Use worker pool for concurrent uploads
	workerCount := s.config.Queue.WorkerCount
	if workerCount > len(files) {
		workerCount = len(files)
	}

	results := make([]*ports.UploadResponse, len(files))
	errors := make([]error, len(files))

	var wg sync.WaitGroup
	jobs := make(chan int, len(files))

	// Start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for idx := range jobs {
				results[idx], errors[idx] = s.Upload(ctx, files[idx])
			}
		}(i)
	}

	// Send jobs
	for i := range files {
		jobs <- i
	}
	close(jobs)

	wg.Wait()

	// Check for errors
	for _, err := range errors {
		if err != nil {
			return nil, fmt.Errorf("batch upload failed: %w", err)
		}
	}

	return results, nil
}

// Helper methods
func (s *fileService) calculateFileHash(file multipart.File) (string, int64, error) {
	hash := sha256.New()
	size, err := io.Copy(hash, file)
	if err != nil {
		return "", 0, err
	}

	// Reset file pointer
	file.Seek(0, 0)

	return hex.EncodeToString(hash.Sum(nil)), size, nil
}

func (s *fileService) validateFile(fileHeader *multipart.FileHeader, maxSize int64) error {
	// Check file size
	if fileHeader.Size > maxSize {
		return fmt.Errorf("file size exceeds limit of %d bytes", maxSize)
	}

	// Check file type
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".pdf":  true,
		".mp4":  true,
		".mov":  true,
		".avi":  true,
		".zip":  true,
	}

	if !allowedTypes[ext] {
		return fmt.Errorf("file type %s not allowed", ext)
	}

	return nil
}

func (s *fileService) generateStorageKey(fileName, hash string) string {
	timestamp := time.Now().Format("2006/01/02")
	// Remove special characters from filename
	cleanName := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '.' || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, fileName)
	return fmt.Sprintf("%s/%s/%s", timestamp, hash[:8], cleanName)
}

func (s *fileService) generatePublicPath(storageKey string, isPublic bool) string {
	if !isPublic {
		return "private/" + storageKey
	}
	return "public/" + storageKey
}

func (s *fileService) createUploadResponse(file *domain.File) *ports.UploadResponse {
	return &ports.UploadResponse{
		FileID:       fmt.Sprintf("%d", file.ID),
		OriginalName: file.OriginalName,
		Size:         file.Size,
		MimeType:     file.MimeType,
		PublicURL:    s.getPublicURL(file),
		StoragePath:  file.StoragePath,
		Hash:         file.Hash,
		CreatedAt:    file.CreatedAt,
	}
}

func (s *fileService) getPublicURL(file *domain.File) string {
	if s.config.Storage.CDNBaseURL != "" {
		return fmt.Sprintf("%s/%s", s.config.Storage.CDNBaseURL, file.PublicPath)
	}
	return file.PublicPath
}

func (s *fileService) getThumbnailURL(file *domain.File) string {
	if file.ThumbnailPath != "" && s.config.Storage.CDNBaseURL != "" {
		return fmt.Sprintf("%s/%s", s.config.Storage.CDNBaseURL, file.ThumbnailPath)
	}
	return file.ThumbnailPath
}

func (s *fileService) cacheFileInfo(ctx context.Context, file *domain.File) {
	key := s.cacheRepo.FileKey(fmt.Sprintf("%d", file.ID))
	data, _ := json.Marshal(file)
	s.cacheRepo.Set(ctx, key, data, 5*time.Minute)
}

func (s *fileService) getFileWithCache(ctx context.Context, fileID string) (*domain.File, error) {
	// Try cache first
	key := s.cacheRepo.FileKey(fileID)
	var file domain.File
	if err := s.cacheRepo.GetJSON(ctx, key, &file); err == nil && file.ID > 0 {
		return &file, nil
	}

	// Try to parse as numeric ID
	var id uint64
	if _, err := fmt.Sscanf(fileID, "%d", &id); err == nil {
		// Get by numeric ID
		f, err := s.fileRepo.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}
		s.cacheFileInfo(ctx, f)
		return f, nil
	}

	// Try UUID
	f, err := s.fileRepo.FindByUUID(ctx, fileID)
	if err != nil {
		return nil, err
	}
	s.cacheFileInfo(ctx, f)
	return f, nil
}

func (s *fileService) isPublicFile(file *domain.File) bool {
	return strings.HasPrefix(file.PublicPath, "public/")
}

func (s *fileService) validateSignature(ctx context.Context, signature string, fileID uint64, userID uint64) bool {
	// Check cache first
	cacheKey := s.cacheRepo.SignatureKey(signature)
	data, err := s.cacheRepo.Get(ctx, cacheKey)
	if err != nil || data == nil {
		return false
	}

	expected := fmt.Sprintf("%d:%d", fileID, userID)
	return string(data) == expected
}

func (s *fileService) generateSignature(fileID uint64, userID uint64, expiry time.Duration) string {
	timestamp := time.Now().Add(expiry).Unix()
	data := fmt.Sprintf("%d:%d:%d", fileID, userID, timestamp)

	// Use HMAC-SHA256
	h := hmac.New(sha256.New, []byte(s.config.Security.SignatureKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *fileService) assembleFileFromChunks(ctx context.Context, sessionID, fileHash string) {
	// Implement chunk assembly logic
	// 1. Get session
	session, err := s.sessionRepo.GetSession(ctx, sessionID)
	if err != nil {
		log.Printf("Failed to get session %s: %v", sessionID, err)
		return
	}

	// 2. Assemble chunks
	var assembledData []byte
	for i := 1; i <= session.TotalChunks; i++ {
		chunkKey := fmt.Sprintf("chunks/%s/%d", sessionID, i)
		reader, _, err := s.storage.Retrieve(ctx, chunkKey)
		if err != nil {
			log.Printf("Failed to retrieve chunk %d: %v", i, err)
			continue
		}

		chunkData, _ := io.ReadAll(reader)
		assembledData = append(assembledData, chunkData...)
		reader.Close()

		// Clean up chunk
		s.storage.Delete(ctx, chunkKey)
	}

	// 3. Create file from assembled data
	// (Implementation depends on your requirements)
	log.Printf("Assembled file from %d chunks, total size: %d bytes",
		session.TotalChunks, len(assembledData))
}

func (s *fileService) processVideoBackground(ctx context.Context, fileID uint64, storagePath string) {
	// Create new context with timeout
	processCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Process video in background
	if s.videoService != nil {
		if err := s.videoService.ProcessVideo(processCtx, fmt.Sprintf("%d", fileID), ports.VideoProcessOptions{
			GenerateThumbnail: true,
			CreateHLS:         true,
			Qualities:         []string{"360p", "720p", "1080p"},
		}); err != nil {
			log.Printf("Failed to process video %d: %v", fileID, err)
		}
	}
}

// Helper function to get image dimensions
type ImageDimensions struct {
	Width  int
	Height int
}

func (s *fileService) getImageDimensions(file multipart.File) (*ImageDimensions, error) {
	// Read first few bytes to detect format
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}

	// Reset file pointer
	file.Seek(0, 0)

	// Try to decode as JPEG/PNG/GIF
	// Note: You'll need to import appropriate image packages
	// For now, return dummy dimensions
	return &ImageDimensions{Width: 800, Height: 600}, nil
}
