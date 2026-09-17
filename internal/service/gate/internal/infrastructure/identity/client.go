package identity

import (
	"context"

	identityv1 "github.com/alexnesterov/rapidlog-api/gen/identity/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/gate/internal/domain/port"
	"github.com/google/uuid"
)

type Client struct {
	api identityv1.IdentityServiceClient
}

func NewClient(client identityv1.IdentityServiceClient) *Client {
	return &Client{api: client}
}

func (c *Client) ResolveSession(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	resp, err := c.api.ResolveSession(ctx, &identityv1.ResolveSessionRequest{SessionId: id.String()})
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(resp.GetUserId())
}

var _ port.IdentityClient = &Client{}
