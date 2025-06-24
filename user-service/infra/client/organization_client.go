package client

import (
	"fmt"
	"log"
	"user/env"

	organizationpb "pb/types/organization"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewOrganizationClient() organizationpb.OrganizationServiceClient {
	url := "localhost"
	if env.ENV_RUNTIME == "release" {
		url = "organization"
	}
	conn, err := grpc.NewClient(fmt.Sprintf("%s%s", url, viper.GetString("grpc.organization.host")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to organization service: %v", err)
	}
	orgClient := organizationpb.NewOrganizationServiceClient(conn)
	return orgClient
}

func NewOrganizationMemberClient() organizationpb.OrganizationMemberServiceClient {
	// conn, err := grpc.NewClient(viper.GetString("grpc.organization.host"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	url := "localhost"
	if env.ENV_RUNTIME == "release" {
		url = "organization"
	}
	conn, err := grpc.NewClient(fmt.Sprintf("%s%s", url, viper.GetString("grpc.group.host")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to organization service: %v", err)
	}
	orgClient := organizationpb.NewOrganizationMemberServiceClient(conn)
	return orgClient
}

func NewGroupClient() organizationpb.GroupServiceClient {
	url := "localhost"
	if env.ENV_RUNTIME == "release" {
		url = "organization"
	}
	conn, err := grpc.NewClient(fmt.Sprintf("%s%s", url, viper.GetString("grpc.group.host")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to group service: %v", err)
	}
	groupClient := organizationpb.NewGroupServiceClient(conn)
	return groupClient
}
