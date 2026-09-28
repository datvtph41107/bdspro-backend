package dto

// ProfileInfoV3Request is GET /v3/user/profile/{access}/{id}.
type ProfileInfoV3Request struct {
	Access    string `path:"access"`
	ID        uint64 `path:"id"`
	Timestamp int64  `form:"timestamp"`
}
