package gateway

import (
	gw "github.com/earthboundkid/gateway/v2"
)

// ResponseWriter implements the http.ResponseWriter interface
// in order to support the API Gateway Lambda HTTP "protocol".
//
//go:fix inline
type ResponseWriter = gw.ResponseWriter

// NewResponse returns a new response writer to capture http output.
//
//go:fix inline
func NewResponse() *ResponseWriter {
	return gw.NewResponse()
}
