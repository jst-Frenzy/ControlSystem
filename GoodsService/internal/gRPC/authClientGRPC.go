package gRPC

import (
	"context"
	"errors"
	"fmt"
	gen "github.com/jst-Frenzy/ControlSystem/protobuf/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type AuthClient struct {
	conn   *grpc.ClientConn
	client gen.AuthServiceClient
}

func NewAuthClient(addr string) (*AuthClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to connect to auth service: %v", err))
	}

	return &AuthClient{
		conn:   conn,
		client: gen.NewAuthServiceClient(conn),
	}, nil
}

func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*gen.ValidateTokenResponse, error) {
	req := &gen.ValidateTokenRequest{AccessToken: token}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.client.ValidateToken(ctx, req)
}

func (c *AuthClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
