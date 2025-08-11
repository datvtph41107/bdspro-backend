package clients

// type OrganizationMemberGrpcClient struct {
// 	Client organizationpb.OrganizationMemberServiceClient
// }

// func NewOrganizationMemberGrpcClient(env_runtime string) *OrganizationMemberGrpcClient {
// 	organizationMemberClient := &OrganizationMemberGrpcClient{}
// 	go func() {
// 		url := "localhost"
// 		name := "organization"
// 		if env_runtime == "release" {
// 			url = name
// 		}

// 		conn, err := grpc.Dial(fmt.Sprintf("%s:8207", url),
// 			grpc.WithInsecure(),
// 			grpc.WithBackoffMaxDelay(5*time.Second),
// 			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
// 		)
// 		if err != nil {
// 			log.Fatalf("failed to dial: %v", err)
// 		}

// 		ctx, _ := context.WithCancel(context.Background())

// 		// Khởi chạy goroutine để theo dõi trạng thái kết nối
// 		go _middleware.MonitorConnection(ctx, name, conn)

// 		client := organizationpb.NewOrganizationMemberServiceClient(conn)
// 		organizationMemberClient.Client = client
// 	}()
// 	return organizationMemberClient
// }

// func (c *OrganizationMemberGrpcClient) CheckUsersInOrganization(ctx context.Context, organizationID uint32, userIds []uint32) ([]*organizationpb.CheckUserInOrganization, error) {
// 	if c.Client == nil {
// 		return nil, errors.New("organization member client not avaiable")
// 	}

// 	users, err := c.Client.CheckUsersInOrganization(ctx, &organizationpb.CheckUsersInOrganizationRequest{OrganizationId: organizationID, UserIds: userIds})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return users.Data, nil
// }
