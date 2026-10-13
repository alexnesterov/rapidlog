package journalgrpc

import (
	"fmt"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	journalv1.JournalServiceClient
	conn *grpc.ClientConn
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to service: %w", err)
	}

	return &Client{
		JournalServiceClient: journalv1.NewJournalServiceClient(conn),
		conn: conn,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}