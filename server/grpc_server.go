package server

import (
	"context"
	"net"
	"rolando/cmd/log"
	"rolando/cmd/services"
	"rolando/config"
	"rolando/server/analytics"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	ChainsService *services.ChainsService
}

func NewGrpcServer(chainsService *services.ChainsService) *GrpcServer {
	return &GrpcServer{
		ChainsService: chainsService,
	}
}

func (s *GrpcServer) Start() {
	lis, err := net.Listen("tcp", "127.0.0.1:"+(config.GrpcPort))
	if err != nil {
		log.Log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	analytics.RegisterAnalyticsServer(grpcServer, analytics.NewAnalyticsServer(s.ChainsService))

	log.Log.Infof("gRPC server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Log.Fatalf("failed to serve: %v", err)
	}
}

// Shutdown gracefully shuts down the gRPC server.
func (s *GrpcServer) Shutdown(ctx context.Context) error {
	return nil
}
