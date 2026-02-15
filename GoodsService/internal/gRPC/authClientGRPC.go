package gRPC

import (
	"context"
	"fmt"
	gen "github.com/jst-Frenzy/ControlSystem/protobuf/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

//go:generate mockgen -source=authClientGRPC.go -destination=../mocks/MockAuthClient.go -package=mock_GoodService

type AuthClient interface {
	ValidateToken(ctx context.Context, token string) (*gen.ValidateTokenResponse, error)
	Close() error
}

type authClient struct {
	conn   *grpc.ClientConn
	client gen.AuthServiceClient
}

func NewAuthClient(addr string) (AuthClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	return &authClient{
		conn:   conn,
		client: gen.NewAuthServiceClient(conn),
	}, nil
}

func (c *authClient) ValidateToken(ctx context.Context, token string) (*gen.ValidateTokenResponse, error) {
	req := &gen.ValidateTokenRequest{AccessToken: token}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.client.ValidateToken(ctx, req)
}

func (c *authClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
