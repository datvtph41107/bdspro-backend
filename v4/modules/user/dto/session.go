package dto

// TileSessionResponse mirrors userpb.GenTileSessionTokenResponse for HTTP.
type TileSessionResponse struct {
	SessionK          string `json:"sessionK"`
	ExpiresIn         int64  `json:"expiresIn"`
	SessionEncryptKey string `json:"sessionEncryptKey"`
}
