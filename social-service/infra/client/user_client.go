package client

import (
	"context"
	"fmt"
	"log"
	pb_social "pb/types/social"
	pb_user "pb/types/user"
	"social/env"
	"time"

	"google.golang.org/grpc"
)

// @bind: social/internal/interface.UserClient
type UserClient struct {
	client pb_user.ProfileServiceClient
}

func UnaryClientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	elapsed := time.Since(start)

	log.Printf("[gRPC Client] Method: %s, Duration: %s, Error: %v", method, elapsed, err)
	return err
}

func NewUserClient() *UserClient {
	userClient := &UserClient{client: nil}
	go func() {
		url := "localhost"
		if env.ENV_RUNTIME == "release" {
			url = "user_grpc"
		}

		conn, err := grpc.Dial(fmt.Sprintf("%s:8201", url),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(UnaryClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}

		ctx, _ := context.WithCancel(context.Background())

		// Khởi chạy goroutine để theo dõi trạng thái kết nối
		go monitorConnection(ctx, conn)

		client := pb_user.NewProfileServiceClient(conn)
		userClient.client = client
	}()

	return userClient
}

func monitorConnection(ctx context.Context, conn *grpc.ClientConn) {
	for {
		state := conn.GetState()
		log.Printf("[gRPC] Connection state: %v", state)

		// Wait until the state changes before continuing the loop
		if !conn.WaitForStateChange(ctx, state) {
			log.Println("[gRPC] Context canceled or connection closed")
			return
		}
	}
}

func (c *UserClient) IsFriends(ctx context.Context, profileID uint64, friendIDs []uint64) (bool, error) {
	return true, nil
}

func (c *UserClient) PbUserToDomain(ctx context.Context, newsFeeds []*pb_social.NewsFeed) {
	if c.client == nil {
		return
	}
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

	profileIDs := make([]uint64, 0, len(profileIDSet))
	for id := range profileIDSet {
		profileIDs = append(profileIDs, id)
	}

	profiles, err := c.client.GetProfileByIds(ctx, &pb_user.GetProfileByIdsRequest{Ids: profileIDs})
	if err != nil {
		return
	}

	profileMap := make(map[uint64]*pb_user.ProfileItem, len(profiles.Profiles))
	for _, profile := range profiles.Profiles {
		profileMap[profile.Id] = profile
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
		nf.FriendTags = []*pb_user.ProfileItem{}
		for _, friendTagID := range nf.FriendTagIds {
			if profile, ok := profileMap[friendTagID]; ok {
				nf.FriendTags = append(nf.FriendTags, profile)
			}
		}
	}

}
