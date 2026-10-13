package grpcclient

import (
	"context"
	"fmt"

	identityv1 "github.com/alexnesterov/rapidlog-api/api/gen/identity/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IdentityClient struct {
	conn *grpc.ClientConn
	client identityv1.IdentityServiceClient
}

func NewIdentityClient(addr string) (*IdentityClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to service: %w", err)
	}

	return &IdentityClient{
		conn: conn,
		client: identityv1.NewIdentityServiceClient(conn),
	}, nil
}

func (c *IdentityClient) Close() {
	c.conn.Close()
}

func (c *IdentityClient) ResolveSession(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	resp, err := c.client.ResolveSession(ctx, &identityv1.ResolveSessionRequest{SessionId: id.String()})
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(resp.GetUserId())
}
