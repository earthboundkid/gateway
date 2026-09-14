package gateway

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/earthboundkid/gateway/v2"
)

// NewRequest returns a new http.Request from the given Lambda event.
//
//go:fix inline
func NewRequest(ctx context.Context, e events.APIGatewayProxyRequest) (*http.Request, error) {
	return gateway.NewRequest(ctx, e)
}
