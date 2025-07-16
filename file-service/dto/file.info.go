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

	Scope  string // "public"
	Type   string // "video"
	Date   string // "2025-07-14"
	UUID   string // "1752485034177081000"
	FileID string // "619"
	Ext    string // "mp4"
}
