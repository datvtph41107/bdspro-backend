package client

import (
	"context"
	clients "pb/clients"
	sharepb "pb/types/shared"
	transactionpb "pb/types/transaction"
	"tx/env"
)

type UserClient struct {
	*clients.UserGrpcClient
}

func NewUserClient() *UserClient {
	return &UserClient{
		UserGrpcClient: clients.MakeUserGrpcClient(env.ENV_RELEASE),
	}
}

func (c *UserClient) MapProfileToTransaction(ctx context.Context, transactions []*transactionpb.TransactionResponse) error {
	ids := make(map[uint64]struct{})
	for _, transaction := range transactions {
		if transaction.FromOf == transactionpb.OwnerType_OWNER_TYPE_USER {
			ids[transaction.FromId] = struct{}{}
		}
	}

	profiles, err := c.GetMapByIDs(ctx, ids)
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

func (c *UserClient) MapProfileToDealContractPb(ctx context.Context, contract *transactionpb.DealContractResponse) {
	ids := make([]uint64, 1)
	ids[0] = uint64(contract.CustomerId)
	profiles, err := c.GetProfileByIds(ctx, ids)
	if err != nil {
		return
	}
	if len(profiles.Profiles) > 0 {
		profile := profiles.Profiles[0]
		profile.Phone = ""
		contract.Customer = profile
	}
}

func (c *UserClient) MapProfileToDealContractPbs(ctx context.Context, contracts []*transactionpb.DealContractItem) {
	ids := make([]uint64, 0)
	for _, contract := range contracts {
		ids = append(ids, uint64(contract.CustomerId))
	}
	profiles, err := c.GetProfileByIds(ctx, ids)
	if err != nil {
		return
	}

	profilesMap := make(map[uint64]*sharepb.ProfileItem)
	for _, profile := range profiles.Profiles {
		profilesMap[profile.Id] = profile
	}

	for _, contract := range contracts {
		profile, ok := profilesMap[contract.CustomerId]
		if ok {
			contract.Customer = profile
			contract.CustomerName = profile.FullName
		}
	}
}
