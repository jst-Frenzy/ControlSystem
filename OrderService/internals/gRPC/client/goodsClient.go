package client

import (
	"context"
	"fmt"
	gen "github.com/jst-Frenzy/ControlSystem/protobuf/gen/goods"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type GoodsClient interface {
	GetItemQuantityAndPrice(ctx context.Context, itemID string) (*gen.ItemQuantityAndPriceResponse, error)
	Close() error
}

type goodsClient struct {
	conn   *grpc.ClientConn
	client gen.GoodsServiceClient
}

func NewGoodsClient(addr string) (GoodsClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed ro connect to goods service: %w", err)
	}

	return &goodsClient{
		conn:   conn,
		client: gen.NewGoodsServiceClient(conn),
	}, nil
}

func (c *goodsClient) GetItemQuantityAndPrice(ctx context.Context, itemID string) (*gen.ItemQuantityAndPriceResponse, error) {
	req := &gen.ItemQuantityAndPriceRequest{ItemId: itemID}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.client.GetItemQuantityAndPrice(ctx, req)
}

func (c *goodsClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}
