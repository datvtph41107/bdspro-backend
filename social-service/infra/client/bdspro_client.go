package client

import "context"

// @bind: social/internal/interface.BdsproClient
type BdsproClient struct {
}

func NewBdsproClient() *BdsproClient {
	return &BdsproClient{}
}

func (c *BdsproClient) OwnerPost(ctx context.Context, postID uint64, profileID uint64) error {
	return nil
}
