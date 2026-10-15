package dto

import "encoding/json"

// RefreshResponse is the refresh-token payload.
type RefreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// RequestLoginResponse is the OTP / login request step response.
type RequestLoginResponse struct {
	AuthID     uint64 `json:"authId,omitempty"`
	Provider   string `json:"provider,omitempty"`
	OAuthID    string `json:"oauthId,omitempty"`
	Existed    bool   `json:"existed"`
	SmsChannel string `json:"smsChannel"`
}

// MarshalJSON hides RefreshToken when serializing JSON.
func (a AuthLoginResponse) MarshalJSON() ([]byte, error) {
	type Alias AuthLoginResponse
	return json.Marshal(&struct {
		Alias
		RefreshToken string `json:"refreshToken,omitempty"`
	}{
		Alias:        (Alias)(a),
		RefreshToken: "",
	})
}

// NewRefreshResponse creates RefreshResponse.
func NewRefreshResponse(accessToken string) RefreshResponse {
	return RefreshResponse{AccessToken: accessToken}
}
