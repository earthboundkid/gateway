// Package gateway provides a drop-in replacement for net/http.ListenAndServe for use in AWS Lambda & API Gateway.
package gateway

import (
	"net/http"

	gw "github.com/earthboundkid/gateway/v2"
)

// ListenAndServe is a drop-in replacement for
// http.ListenAndServe for use within AWS Lambda.
// Because the standard addr string is not used,
// it is replaced with host, which API Gateway
// does not always send with events.
//
// ListenAndServe never returns.
//
//go:fix inline
func ListenAndServe(host string, h http.Handler) error {
	return gw.ListenAndServe(host, h)
}
