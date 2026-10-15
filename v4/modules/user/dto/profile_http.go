package dto

// ProfileInfoV3Request is GET /v3/user/profile/{access}/{id}.
type ProfileInfoV3Request struct {
	Access    string `path:"access"`
	ID        uint64 `path:"id"`
	Timestamp int64  `form:"timestamp"`
}

// ProfileInfoRequest is GET /v2/user/profile/info/{id}.
type ProfileInfoRequest struct {
	ID uint64 `path:"id"`
}

// SearchProfilePublicRequest is GET /v2/user/profile/search/public.
type SearchProfilePublicRequest struct {
	Text string `form:"text"`
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
}

// FindByPhoneRequest is GET /v2/user/profile/find-by-phone.
type FindByPhoneRequest struct {
	Phone string `form:"phone"`
}

// StatusOnlineRequest is GET /v2/user/status-online/{id}.
type StatusOnlineRequest struct {
	ID uint64 `path:"id"`
}
