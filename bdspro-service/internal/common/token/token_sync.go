package token

import (
	"encoding/base64"
	"encoding/json"
)

type ZSetPageToken struct {
	LastScore int64  `json:"s"`
	LastID    uint64 `json:"i"`
}

func EncodeZSetPageToken(t ZSetPageToken) (string, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func DecodeZSetPageToken(tokenStr string) (*ZSetPageToken, error) {
	if tokenStr == "" {
		return &ZSetPageToken{}, nil
	}

	b, err := base64.StdEncoding.DecodeString(tokenStr)
	if err != nil {
		return nil, err
	}

	var t ZSetPageToken
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, err
	}

	return &t, nil
}
