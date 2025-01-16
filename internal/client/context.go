package client

import (
	"context"
	"github.com/mcoot/crosswordgame-go/internal/utils"
)

const (
	ContextKeyClient utils.ContextKey = "client"
)

func AddClientToContext(ctx context.Context, c *Client) context.Context {
	return context.WithValue(ctx, ContextKeyClient, c)
}

func GetClient(ctx context.Context) *Client {
	return ctx.Value(ContextKeyClient).(*Client)
}
