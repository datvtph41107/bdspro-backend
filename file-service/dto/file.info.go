package dto

// FileInfo chứa thông tin về file đã lưu
type FileInfo struct {
	ID            uint64 `json:"id"`
	FileName      string `json:"filename"`
	RelativePath  string `json:"path"`
	AbsolutePath  string `json:"absolutePath"`
	ThumbnailPath string `json:"thumbnailPath"`
	Extension     string `json:"extension"`
	Size          int64  `json:"size"`
	Hash          string `json:"hash"`
	ContentType   string `json:"contentType"`
	Description   string `json:"description"`
	Mine          string `json:"mine"`
}
