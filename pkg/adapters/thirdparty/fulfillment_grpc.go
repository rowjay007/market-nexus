package thirdparty

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type FulfillmentClient struct {
	conn *grpc.ClientConn
}

func NewFulfillmentClient(address string) (*FulfillmentClient, error) {
	if address == "" {
		return nil, errors.New("fulfillment grpc address required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return &FulfillmentClient{conn: conn}, nil
}

func (c *FulfillmentClient) Close() error {
	return c.conn.Close()
}
