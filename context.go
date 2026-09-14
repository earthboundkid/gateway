package gateway

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/earthboundkid/gateway/v2"
)

// RequestContext returns the APIGatewayProxyRequestContext value stored in ctx.
//
//go:fix inline
func RequestContext(ctx context.Context) (events.APIGatewayProxyRequestContext, bool) {
	return gateway.RequestContext(ctx)
}
