package gateway

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	gw "github.com/earthboundkid/gateway/v2"
)

// RequestContext returns the APIGatewayProxyRequestContext value stored in ctx.
//
//go:fix inline
func RequestContext(ctx context.Context) (events.APIGatewayProxyRequestContext, bool) {
	return gw.RequestContext(ctx)
}
