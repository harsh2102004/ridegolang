package billing

import (
	"context"
	"log"
	"os"

	pb "org/api-core/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Client pb.BillingServiceClient

func Connect() {
	addr := os.Getenv("BILLING_GRPC_ADDR")

	if addr == "" {
		addr = "localhost:50052"
	}

	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	Client = pb.NewBillingServiceClient(conn)

	log.Println("Connected to billing service")
}

func GetUserSubscription(
	userID string,
) (string, string, error) {

	res, err := Client.GetUserSubscription(
		context.Background(),
		&pb.GetUserSubscriptionRequest{
			UserId: userID,
		},
	)

	if err != nil {
		return "", "", err
	}

	return res.Plan, res.Status, nil
}
