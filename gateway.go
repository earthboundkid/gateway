// Package gateway provides a drop-in replacement for net/http.ListenAndServe for use in AWS Lambda & API Gateway.
package gateway

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// ListenAndServe is a drop-in replacement for
// http.ListenAndServe for use within AWS Lambda.
// Because the standard addr string is not used,
// it is replaced with host, which API Gateway
// does not always send with events.
//
// ListenAndServe never returns.
func ListenAndServe(host string, h http.Handler) error {
	if h == nil {
		h = http.DefaultServeMux
	}

	lambda.Start(func(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		ctx = setHost(ctx, host)
		r, err := NewRequest(ctx, e)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		w := NewResponse()
		h.ServeHTTP(w, r)
		return w.End(), nil
	})

	panic("unreachable")
}
