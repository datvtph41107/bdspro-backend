package client

import (
	"context"
	clients "pb/clients"
	transactionpb "pb/types/transaction"
	"tx/env"
)

type UserClient struct {
	Client *clients.UserGrpcClient
}

func NewUserClient() *UserClient {
	return &UserClient{
		Client: clients.MakeUserGrpcClient(env.ENV_RELEASE),
	}
}

func (c *UserClient) MapProfileToTransaction(ctx context.Context, transactions []*transactionpb.TransactionResponse) error {
	ids := make(map[uint64]struct{})
	// for _, transaction := range transactions {
	// 	if transaction.FromO == transactionpb.OwnerType_OWNER_TYPE_USER {
	// 		ids[transaction.FromId] = struct{}{}
	// 	}
	// }

	profiles, err := c.Client.GetMapByIDs(ctx, ids)
	if err != nil {
		return err
	}

	for _, transaction := range transactions {
		profile, ok := profiles[transaction.FromId]
		if ok && transaction.FromOf == transactionpb.OwnerType_OWNER_TYPE_USER {
			transaction.Profile = profile
		}
	}

	return nil
}
