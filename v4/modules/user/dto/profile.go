package dto

// ProfileInfoV3Request is the query shape for GET /v3/user/profile/{access}/{id}.
type ProfileInfoV3Request struct {
	Timestamp int64 `form:"timestamp"`
}
