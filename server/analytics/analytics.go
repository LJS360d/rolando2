package analytics

import (
	"context"
	"rolando/cmd/model"
	"rolando/cmd/services"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type AnalyticsServerImpl struct {
	AnalyticsServer
	chainsService *services.ChainsService
}

func NewAnalyticsServer(chainsService *services.ChainsService) *AnalyticsServerImpl {
	return &AnalyticsServerImpl{
		chainsService: chainsService,
	}
}

func (s *AnalyticsServerImpl) GetChainAnalytics(ctx context.Context, req *wrapperspb.StringValue) (*ChainAnalytics, error) {
	chainId := req.GetValue()
	chain, err := s.chainsService.GetChain(chainId)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	analyzer := model.NewMarkovChainAnalyzer(chain)
	chainDoc, err := s.chainsService.GetChainDocument(chainId)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	rawAnalytics := analyzer.GetRawAnalytics()
	chainAnalytics := &ChainAnalytics{
		ComplexityScore: uint32(rawAnalytics.ComplexityScore),
		Gifs:            uint32(rawAnalytics.Gifs),
		Images:          uint32(rawAnalytics.Images),
		Videos:          uint32(rawAnalytics.Videos),
		ReplyRate:       uint32(rawAnalytics.ReplyRate),
		Words:           uint32(rawAnalytics.Words),
		Messages:        uint32(rawAnalytics.Messages),
		Bytes:           uint64(rawAnalytics.Size),
		Id:              chain.ID,
		Name:            chainDoc.Name,
	}
	return chainAnalytics, nil
}

func (s *AnalyticsServerImpl) GetAllChainsAnalytics(req *emptypb.Empty, stream grpc.ServerStreamingServer[ChainAnalytics]) error {
	// Retrieve chains
	chains, err := s.chainsService.GetAllChains()
	if err != nil {
		return status.Error(codes.FailedPrecondition, err.Error())
	}

	// Start streaming the actual data (chain analytics)
	for _, chain := range chains {
		analyzer := model.NewMarkovChainAnalyzer(chain)
		chainDoc, err := s.chainsService.GetChainDocument(chain.ID)
		if err != nil {
			return status.Error(codes.FailedPrecondition, err.Error())
		}

		// Get the analytics data
		rawAnalytics := analyzer.GetRawAnalytics()
		chainAnalytics := &ChainAnalytics{
			ComplexityScore: uint32(rawAnalytics.ComplexityScore),
			Gifs:            uint32(rawAnalytics.Gifs),
			Images:          uint32(rawAnalytics.Images),
			Videos:          uint32(rawAnalytics.Videos),
			ReplyRate:       uint32(rawAnalytics.ReplyRate),
			Words:           uint32(rawAnalytics.Words),
			Messages:        uint32(rawAnalytics.Messages),
			Bytes:           uint64(rawAnalytics.Size),
			Id:              chain.ID,
			Name:            chainDoc.Name,
		}

		// Send each chain's analytics as part of the stream
		if err := stream.Send(chainAnalytics); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
	return nil
}
