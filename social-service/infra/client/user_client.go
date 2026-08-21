package client

import (
	"context"
	"errors"
	"pb/clients"
	sharepb "pb/types/shared"
	pb_social "pb/types/social"
	"social/internal/dto"
)

// @bind: social/internal/interface.UserClient
type UserClient struct {
	client *clients.UserGrpcClient
}

func NewUserClient(rpcClient *clients.UserGrpcClient) *UserClient {
	return &UserClient{client: rpcClient}
}

func (c *UserClient) IsFriends(ctx context.Context, profileID uint64, friendIDs []uint64) (bool, error) {
	return true, nil
}

func (c *UserClient) GetMapByIDs(ctx context.Context, profileIDSet map[uint64]struct{}) (map[uint64]*sharepb.ProfileItem, error) {
	if c.client == nil {
		return nil, errors.New("user client not avaiable")
	}

	profileIDs := make([]uint64, 0, len(profileIDSet))
	for id := range profileIDSet {
		profileIDs = append(profileIDs, id)
	}

	profileMap, err := c.client.GetMapByIDs(ctx, profileIDSet)
	if err != nil {
		return nil, err
	}

	return profileMap, nil
}

func (c *UserClient) GetUser(ctx context.Context, profileID uint64) (*dto.UserProfile, error) {
	profile, err := c.client.Client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: []uint64{profileID}})
	if err != nil {
		return nil, err
	}

	if len(profile.Profiles) == 0 {
		return nil, errors.New("profile not found")
	}

	return &dto.UserProfile{
		ID:       profile.Profiles[0].Id,
		FullName: profile.Profiles[0].FullName,
		Avatar:   profile.Profiles[0].Avatar,
	}, nil
}

func (c *UserClient) GetUsers(ctx context.Context, profileIDs []uint64) ([]*dto.UserProfile, error) {
	profile, err := c.client.Client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: profileIDs})
	if err != nil {
		return nil, err
	}

	if len(profile.Profiles) == 0 {
		return nil, errors.New("profile not found")
	}

	results := make([]*dto.UserProfile, len(profile.Profiles))
	for i, profile := range profile.Profiles {
		results[i] = &dto.UserProfile{
			ID:       profile.Id,
			FullName: profile.FullName,
			Avatar:   profile.Avatar,
		}
	}

	return results, nil
}

// -- mapping pb to domain --

func (c *UserClient) PbUserToNewsFeed(ctx context.Context, newsFeeds []*pb_social.NewsFeed) {
	profileIDSet := make(map[uint64]struct{})
	for _, nf := range newsFeeds {
		if nf.CreatedBy != nil {
			profileIDSet[*nf.CreatedBy] = struct{}{}
		}
		if nf.FriendTagIds != nil {
			for _, friendTagID := range nf.FriendTagIds {
				profileIDSet[friendTagID] = struct{}{}
			}
		}
	}

	profileMap, err := c.GetMapByIDs(ctx, profileIDSet)
	if err != nil {
		return
	}

	// Gán Author và FriendTags cho từng newsFeed
	for _, nf := range newsFeeds {
		// Gán Author
		if nf.CreatedBy != nil {
			if profile, ok := profileMap[*nf.CreatedBy]; ok {
				nf.Author = profile
			}
		}

		// Gán FriendTags
		nf.FriendTags = []*sharepb.ProfileItem{}
		for _, friendTagID := range nf.FriendTagIds {
			if profile, ok := profileMap[friendTagID]; ok {
				nf.FriendTags = append(nf.FriendTags, profile)
			}
		}
	}

}

func (c *UserClient) PbUserToComment(ctx context.Context, comments []*pb_social.CommentResponse) (map[uint64]*sharepb.ProfileItem, error) {
	profileIDSet := make(map[uint64]struct{})
	for _, comment := range comments {
		if comment.CreatedBy != nil {
			profileIDSet[*comment.CreatedBy] = struct{}{}
		}
	}

	profileMap, err := c.GetMapByIDs(ctx, profileIDSet)
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		if comment.CreatedBy != nil {
			if profile, ok := profileMap[*comment.CreatedBy]; ok {
				comment.Author = profile
			}
		}
	}

	return profileMap, nil
}

func (c *UserClient) PbUserToReels(ctx context.Context, newsFeeds []*pb_social.ReelResponse) {
	profileIDSet := make(map[uint64]struct{})
	for _, nf := range newsFeeds {
		if nf.CreatedBy != 0 {
			profileIDSet[nf.CreatedBy] = struct{}{}
		}
	}

	profileMap, err := c.GetMapByIDs(ctx, profileIDSet)
	if err != nil {
		return
	}

	// Gán Author và FriendTags cho từng newsFeed
	for _, nf := range newsFeeds {
		// Gán Author
		if nf.CreatedBy != 0 {
			if profile, ok := profileMap[nf.CreatedBy]; ok {
				nf.AuthorName = profile.FullName
				nf.AuthorAvatar = profile.Avatar
			}
		}
	}

}
