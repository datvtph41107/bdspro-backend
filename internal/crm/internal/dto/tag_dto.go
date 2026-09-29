package dto

type AdminTagRequest struct {
	ID        uint64
	Name      string
	IsActive  bool
	IsDefault bool
}

type TagListRequest struct {
	Page uint32
	Size uint32
	Text string
}

type AdminTagListRequest struct {
	Page uint32
	Size uint32
	Text string
}
