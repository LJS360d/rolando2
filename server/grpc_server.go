package server

import (
	"context"
	"net"
	"rolando/cmd/log"
	"rolando/cmd/services"
	"rolando/config"
	"rolando/server/analytics"
	"rolando/server/auth"
	"rolando/server/bot"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	ChainsService  *services.ChainsService
	DiscordSession *discordgo.Session
}

func NewGrpcServer(chainsService *services.ChainsService, discordSession *discordgo.Session) *GrpcServer {
	return &GrpcServer{
		ChainsService:  chainsService,
		DiscordSession: discordSession,
	}
}

func (s *GrpcServer) Start() {
	lis, err := net.Listen("tcp", "127.0.0.1:"+(config.GrpcPort))
	if err != nil {
		log.Log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	// Register services
	analytics.RegisterAnalyticsServer(grpcServer, analytics.NewAnalyticsServer(s.ChainsService))
	bot.RegisterBotServer(grpcServer, bot.NewBotServer(s.DiscordSession))
	auth.RegisterAuthServer(grpcServer, auth.NewAuthServerImpl())

	log.Log.Infof("gRPC server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Log.Fatalf("failed to serve: %v", err)
	}
}

// Shutdown gracefully shuts down the gRPC server.
func (s *GrpcServer) Shutdown(ctx context.Context) error {
	return nil
}
