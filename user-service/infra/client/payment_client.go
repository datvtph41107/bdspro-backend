package client

import (
	"fmt"
	"log"
	"user/env"

	paymentpb "pb/types/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
func NewPaymentClient() paymentpb.PaymentServiceClient {
	url := "localhost"
	if env.ENV_RUNTIME == "release" {
		url = "payment"
	}
	conn, err := grpc.NewClient(fmt.Sprintf("%s:50052", url), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to payment service: %v", err)
	}
	paymentClient := paymentpb.NewPaymentServiceClient(conn)
	return paymentClient
}