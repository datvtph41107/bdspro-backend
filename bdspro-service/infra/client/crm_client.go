package client

import (
	"bdspro/env"
	"pb/clients"
)

type CrmClient struct {
	*clients.CrmGrpcClient
}

func NewCrmClient() *CrmClient {
	return &CrmClient{
		clients.MakeCrmGrpcClient(env.ENV_RUNTIME),
	}
}

// GetContactById retrieves contact information by ID using internal service
// TODO: Uncomment sau khi generate proto code với method GetContactById
