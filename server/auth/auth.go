package auth

import (
	"context"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type AuthServerImpl struct {
	AuthServer
}

func NewAuthServerImpl() *AuthServerImpl {
	return &AuthServerImpl{}
}

func (s *AuthServerImpl) SignToken(ctx context.Context, req *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	return nil, nil
}
