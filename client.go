package main

import (
	"context"
	"io"
	"log"
	"rolando/server/bot"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	// Establish a connection to the server
	conn, err := grpc.Dial("localhost:5051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()

	client := bot.NewBotClient(conn)

	stream, err := client.GetBotGuilds(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("could not get bot guilds: %v", err)
	}
	guilds, err := FullStreamRecv(stream)
	if err != nil {
		log.Fatalf("could not receive bot guilds: %v", err)
	}
	g := guilds[0]
	log.Printf("%v", g)
	b := &bot.BroadcastMessageRequest{
		Content: "Test",
		Guilds: []*bot.BroadcastGuildRequest{
			{
				Id:       g.Id,
				Optional: &bot.BroadcastGuildRequest_ChannelId{ChannelId: "910625971865001994"},
			},
		},
	}
	client.BroadcastMessage(context.Background(), b)
}

func FullStreamRecv[T any](stream grpc.ServerStreamingClient[T]) ([]*T, error) {
	var results []*T
	for {
		var msg T
		err := stream.RecvMsg(&msg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		results = append(results, &msg)
	}
	return results, nil
}
