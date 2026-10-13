package identitygrpc

import (
	"fmt"

	identityv1 "github.com/alexnesterov/rapidlog-api/api/gen/identity/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	identityv1.IdentityServiceClient
	conn *grpc.ClientConn
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to service: %w", err)
	}

	return &Client{
		IdentityServiceClient: identityv1.NewIdentityServiceClient(conn),
		conn: conn,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
