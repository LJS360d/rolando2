package analytics

import (
	"context"
	"rolando/cmd/model"
	"rolando/cmd/services"

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
	rawAnalytics := analyzer.GetRawAnalytics()
	chainAnalytics := &ChainAnalytics{
		ComplexityScore: int32(rawAnalytics.ComplexityScore),
		Gifs:            int32(rawAnalytics.Gifs),
		Images:          int32(rawAnalytics.Images),
		Videos:          int32(rawAnalytics.Videos),
		ReplyRate:       int32(rawAnalytics.ReplyRate),
		Words:           int32(rawAnalytics.Words),
		Messages:        int32(rawAnalytics.Messages),
		Bytes:           int64(rawAnalytics.Size),
	}
	return chainAnalytics, nil
}

func (s *AnalyticsServerImpl) GetAllChainsAnalytics(ctx context.Context, req *emptypb.Empty) (*ChainAnalyticsList, error) {
	chains, err := s.chainsService.GetAllChains()
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	var chainAnalyticsList ChainAnalyticsList
	for _, chain := range chains {
		analyzer := model.NewMarkovChainAnalyzer(chain)
		rawAnalytics := analyzer.GetRawAnalytics()
		chainAnalytics := &ChainAnalytics{
			ComplexityScore: int32(rawAnalytics.ComplexityScore),
			Gifs:            int32(rawAnalytics.Gifs),
			Images:          int32(rawAnalytics.Images),
			Videos:          int32(rawAnalytics.Videos),
			ReplyRate:       int32(rawAnalytics.ReplyRate),
			Words:           int32(rawAnalytics.Words),
			Messages:        int32(rawAnalytics.Messages),
			Bytes:           int64(rawAnalytics.Size),
		}
		chainAnalyticsList.Analytics = append(chainAnalyticsList.Analytics, chainAnalytics)
	}
	return &chainAnalyticsList, nil
}
