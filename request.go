package gateway

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	gw "github.com/earthboundkid/gateway/v2"
)

// NewRequest returns a new http.Request from the given Lambda event.
//
//go:fix inline
func NewRequest(ctx context.Context, e events.APIGatewayProxyRequest) (*http.Request, error) {
	return gw.NewRequest(ctx, e)
}
