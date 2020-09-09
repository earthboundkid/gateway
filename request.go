package gateway

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func wraperr(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

func mergeValues(single map[string]string, multi map[string][]string) map[string][]string {
	m := make(map[string][]string, len(single)+len(multi))
	for k, v := range single {
		m[k] = []string{v}
	}
	// Let multi trump single if both are set
	for k, v := range multi {
		m[k] = v
	}
	return m
}

// NewRequest returns a new http.Request from the given Lambda event.
func NewRequest(ctx context.Context, e events.APIGatewayProxyRequest) (*http.Request, error) {
	var header http.Header = mergeValues(e.Headers, e.MultiValueHeaders)
	host := header.Get("Host")
	if host == "" {
		host = getHost(ctx)
	}
	hostURL := &url.URL{
		Scheme: header.Get("X-Forwarded-Proto"),
		Host:   host,
		Path:   "/",
	}

	// path
	u, err := hostURL.Parse(e.Path)
	if err != nil {
		return nil, wraperr(err, "parsing path")
	}

	{
		var query url.Values = mergeValues(
			e.QueryStringParameters, e.MultiValueQueryStringParameters)
		u.RawQuery = query.Encode()
	}

	// base64 encoded body
	body := e.Body
	if e.IsBase64Encoded {
		b, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return nil, wraperr(err, "decoding base64 body")
		}
		body = string(b)
	}

	// new request
	req, err := http.NewRequest(e.HTTPMethod, u.String(), strings.NewReader(body))
	if err != nil {
		return nil, wraperr(err, "creating request")
	}

	// manually set RequestURI because NewRequest is for clients and req.RequestURI is for servers
	req.RequestURI = u.RequestURI()

	// remote addr
	req.RemoteAddr = e.RequestContext.Identity.SourceIP

	// header fields
	req.Header = header

	// content-length
	if req.Header.Get("Content-Length") == "" && body != "" {
		req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	}

	// custom fields
	req.Header.Set("X-Request-Id", e.RequestContext.RequestID)
	req.Header.Set("X-Stage", e.RequestContext.Stage)

	// custom context values
	req = req.WithContext(newContext(ctx, e))

	// xray support
	if traceID := ctx.Value("x-amzn-trace-id"); traceID != nil {
		req.Header.Set("X-Amzn-Trace-Id", fmt.Sprintf("%v", traceID))
	}

	return req, nil
}
