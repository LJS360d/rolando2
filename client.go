package main

import (
	"context"
	"fmt"
	"log"

	"rolando/server/analytics"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.NewClient("localhost:51902", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new Analytics client
	client := analytics.NewAnalyticsClient(conn)

	// Call GetChainAnalytics
	testChainID := &wrapperspb.StringValue{Value: "123"} // Sample Chain ID
	resp, err := client.GetChainAnalytics(context.Background(), testChainID)
	if err != nil {
		log.Fatalf("could not get chain analytics: %v", err)
	}
	fmt.Printf("Chain Analytics: %+v\n", resp)

	// Call GetAllChainsAnalytics
	respAll, err := client.GetAllChainsAnalytics(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("could not get all chains analytics: %v", err)
	}
	fmt.Printf("All Chain Analytics: %+v\n", respAll)
}
