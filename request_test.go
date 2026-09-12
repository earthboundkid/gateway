package gateway

import (
	"context"
	"io"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/earthboundkid/assert"
)

func TestNewRequest_path(t *testing.T) {
	be := assert.FailsNow(t)
	e := events.APIGatewayProxyRequest{
		Path: "/pets/luna",
	}

	r := be.OK(NewRequest(context.Background(), e))

	be.Equal(r.Method, "GET")
	be.Equal(r.URL.Path, `/pets/luna`)
	be.Equal(r.URL.String(), `/pets/luna`)
	be.Equal(r.RequestURI, `/pets/luna`)
}

func TestNewRequest_method(t *testing.T) {
	be := assert.FailsNow(t)
	e := events.APIGatewayProxyRequest{
		HTTPMethod: "DELETE",
		Path:       "/pets/luna",
	}

	r := be.OK(NewRequest(t.Context(), e))

	be.Equal(r.Method, "DELETE")
}

func TestNewRequest_queryString(t *testing.T) {
	be := assert.FailsNow(t)
	e := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/pets",
		QueryStringParameters: map[string]string{
			"order":  "desc",
			"fields": "name,species",
		},
	}

	r := be.OK(NewRequest(t.Context(), e))

	be.Equal(r.URL.String(), `/pets?fields=name%2Cspecies&order=desc`)
	be.Equal(r.URL.Query().Get("order"), `desc`)
}

func TestNewRequest_multiValueQueryString(t *testing.T) {
	be := assert.FailsNow(t)
	e := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/pets",
		MultiValueQueryStringParameters: map[string][]string{
			"multi_fields": {"name", "species"},
			"multi_arr[]":  {"arr1", "arr2"},
		},
		QueryStringParameters: map[string]string{
			"order":  "desc",
			"fields": "name,species",
		},
	}

	r := be.OK(NewRequest(t.Context(), e))

	be.Equal(`/pets?fields=name%2Cspecies&multi_arr%5B%5D=arr1&multi_arr%5B%5D=arr2&multi_fields=name&multi_fields=species&order=desc`, r.URL.String())
	be.SlicesEqual([]string{"name", "species"}, r.URL.Query()["multi_fields"])
	be.SlicesEqual([]string{"arr1", "arr2"}, r.URL.Query()["multi_arr[]"])
	be.Equal(`/pets?fields=name%2Cspecies&multi_arr%5B%5D=arr1&multi_arr%5B%5D=arr2&multi_fields=name&multi_fields=species&order=desc`, r.RequestURI)
}

func TestNewRequest_remoteAddr(t *testing.T) {
	be := assert.FailsNow(t)
	e := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/pets",
		RequestContext: events.APIGatewayProxyRequestContext{
			Identity: events.APIGatewayRequestIdentity{
				SourceIP: "1.2.3.4",
			},
		},
	}

	r := be.OK(NewRequest(t.Context(), e))

	be.Equal(`1.2.3.4`, r.RemoteAddr)
}

func TestNewRequest_header(t *testing.T) {
	be := assert.FailsNow(t)
	e := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/pets",
		Body:       `{ "name": "Tobi" }`,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Foo":        "bar",
			"Host":         "example.com",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			RequestID: "1234",
			Stage:     "prod",
		},
	}
	ctx := setHost(t.Context(), "xxx")
	r := be.OK(NewRequest(ctx, e))

	be.Equal(`example.com`, r.Host)
	be.Equal(`prod`, r.Header.Get("X-Stage"))
	be.Equal(`1234`, r.Header.Get("X-Request-Id"))
	be.Equal(`18`, r.Header.Get("Content-Length"))
	be.Equal(`application/json`, r.Header.Get("Content-Type"))
	be.Equal(`bar`, r.Header.Get("X-Foo"))
}

func TestNewRequest_host(t *testing.T) {
	be := assert.FailsNow(t)
	e := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/pets",
		Body:       `{ "name": "Tobi" }`,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Foo":        "bar",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			RequestID: "1234",
			Stage:     "prod",
		},
	}
	ctx := setHost(t.Context(), "example.com")
	r := be.OK(NewRequest(ctx, e))
	be.Equal(r.Host, `example.com`)
}

func TestNewRequest_multiHeader(t *testing.T) {
	be := assert.FailsNow(t)
	e := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/pets",
		Body:       `{ "name": "Tobi" }`,
		MultiValueHeaders: map[string][]string{
			"X-APEX":   {"apex1", "apex2"},
			"X-APEX-2": {"apex-1", "apex-2"},
		},
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Foo":        "bar",
			"Host":         "example.com",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			RequestID: "1234",
			Stage:     "prod",
		},
	}

	r := be.OK(NewRequest(t.Context(), e))

	be.Equal(`example.com`, r.Host)
	be.Equal(`prod`, r.Header.Get("X-Stage"))
	be.Equal(`1234`, r.Header.Get("X-Request-Id"))
	be.Equal(`18`, r.Header.Get("Content-Length"))
	be.Equal(`application/json`, r.Header.Get("Content-Type"))
	be.Equal(`bar`, r.Header.Get("X-Foo"))
	be.SlicesEqual([]string{"apex1", "apex2"}, r.Header["X-APEX"])
	be.SlicesEqual([]string{"apex-1", "apex-2"}, r.Header["X-APEX-2"])
}

func TestNewRequest_body(t *testing.T) {
	be := assert.FailsNow(t)
	e := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/pets",
		Body:       `{ "name": "Tobi" }`,
	}

	r := be.OK(NewRequest(t.Context(), e))
	b := be.OK(io.ReadAll(r.Body))
	be.Equal(string(b), `{ "name": "Tobi" }`)
}

func TestNewRequest_bodyBinary(t *testing.T) {
	be := assert.FailsNow(t)
	e := events.APIGatewayProxyRequest{
		HTTPMethod:      "POST",
		Path:            "/pets",
		Body:            `aGVsbG8gd29ybGQK`,
		IsBase64Encoded: true,
	}

	r := be.OK(NewRequest(t.Context(), e))
	b := be.OK(io.ReadAll(r.Body))
	be.Equal(string(b), "hello world\n")
}

func TestNewRequest_context(t *testing.T) {
	be := assert.FailsNow(t)
	e := events.APIGatewayProxyRequest{}
	ctx := context.WithValue(t.Context(), "key", "value")
	r := be.OK(NewRequest(ctx, e))
	v := r.Context().Value("key").(string)
	be.Equal(v, "value")
}
