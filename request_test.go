package gateway

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/carlmjohnson/be"
)

func TestNewRequest_path(t *testing.T) {
	e := events.APIGatewayProxyRequest{
		Path: "/pets/luna",
	}

	r, err := NewRequest(context.Background(), e)
	be.NilErr(t, err)

	be.Equal(t, "GET", r.Method)
	be.Equal(t, `/pets/luna`, r.URL.Path)
	be.Equal(t, `/pets/luna`, r.URL.String())
	be.Equal(t, `/pets/luna`, r.RequestURI)
}

func TestNewRequest_method(t *testing.T) {
	e := events.APIGatewayProxyRequest{
		HTTPMethod: "DELETE",
		Path:       "/pets/luna",
	}

	r, err := NewRequest(context.Background(), e)
	be.NilErr(t, err)

	be.Equal(t, "DELETE", r.Method)
}

func TestNewRequest_queryString(t *testing.T) {
	e := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/pets",
		QueryStringParameters: map[string]string{
			"order":  "desc",
			"fields": "name,species",
		},
	}

	r, err := NewRequest(context.Background(), e)
	be.NilErr(t, err)

	be.Equal(t, `/pets?fields=name%2Cspecies&order=desc`, r.URL.String())
	be.Equal(t, `desc`, r.URL.Query().Get("order"))
}

func TestNewRequest_multiValueQueryString(t *testing.T) {
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

	r, err := NewRequest(context.Background(), e)
	be.NilErr(t, err)

	be.Equal(t, `/pets?fields=name%2Cspecies&multi_arr%5B%5D=arr1&multi_arr%5B%5D=arr2&multi_fields=name&multi_fields=species&order=desc`, r.URL.String())
	be.AllEqual(t, []string{"name", "species"}, r.URL.Query()["multi_fields"])
	be.AllEqual(t, []string{"arr1", "arr2"}, r.URL.Query()["multi_arr[]"])
	be.Equal(t, `/pets?fields=name%2Cspecies&multi_arr%5B%5D=arr1&multi_arr%5B%5D=arr2&multi_fields=name&multi_fields=species&order=desc`, r.RequestURI)
}

func TestNewRequest_remoteAddr(t *testing.T) {
	e := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/pets",
		RequestContext: events.APIGatewayProxyRequestContext{
			Identity: events.APIGatewayRequestIdentity{
				SourceIP: "1.2.3.4",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	be.NilErr(t, err)

	be.Equal(t, `1.2.3.4`, r.RemoteAddr)
}

func TestNewRequest_header(t *testing.T) {
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
	ctx := context.Background()
	ctx = setHost(ctx, "xxx")
	r, err := NewRequest(ctx, e)
	be.NilErr(t, err)

	be.Equal(t, `example.com`, r.Host)
	be.Equal(t, `prod`, r.Header.Get("X-Stage"))
	be.Equal(t, `1234`, r.Header.Get("X-Request-Id"))
	be.Equal(t, `18`, r.Header.Get("Content-Length"))
	be.Equal(t, `application/json`, r.Header.Get("Content-Type"))
	be.Equal(t, `bar`, r.Header.Get("X-Foo"))
}

func TestNewRequest_host(t *testing.T) {
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
	ctx := context.Background()
	ctx = setHost(ctx, "example.com")
	r, err := NewRequest(ctx, e)
	be.NilErr(t, err)
	be.Equal(t, `example.com`, r.Host)
}

func TestNewRequest_multiHeader(t *testing.T) {
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

	r, err := NewRequest(context.Background(), e)
	be.NilErr(t, err)

	be.Equal(t, `example.com`, r.Host)
	be.Equal(t, `prod`, r.Header.Get("X-Stage"))
	be.Equal(t, `1234`, r.Header.Get("X-Request-Id"))
	be.Equal(t, `18`, r.Header.Get("Content-Length"))
	be.Equal(t, `application/json`, r.Header.Get("Content-Type"))
	be.Equal(t, `bar`, r.Header.Get("X-Foo"))
	be.AllEqual(t, []string{"apex1", "apex2"}, r.Header["X-APEX"])
	be.AllEqual(t, []string{"apex-1", "apex-2"}, r.Header["X-APEX-2"])
}

func TestNewRequest_body(t *testing.T) {
	e := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/pets",
		Body:       `{ "name": "Tobi" }`,
	}

	r, err := NewRequest(context.Background(), e)
	be.NilErr(t, err)

	b, err := ioutil.ReadAll(r.Body)
	be.NilErr(t, err)

	be.Equal(t, `{ "name": "Tobi" }`, string(b))
}

func TestNewRequest_bodyBinary(t *testing.T) {
	e := events.APIGatewayProxyRequest{
		HTTPMethod:      "POST",
		Path:            "/pets",
		Body:            `aGVsbG8gd29ybGQK`,
		IsBase64Encoded: true,
	}

	r, err := NewRequest(context.Background(), e)
	be.NilErr(t, err)

	b, err := ioutil.ReadAll(r.Body)
	be.NilErr(t, err)

	be.Equal(t, "hello world\n", string(b))
}

func TestNewRequest_context(t *testing.T) {
	e := events.APIGatewayProxyRequest{}
	ctx := context.WithValue(context.Background(), "key", "value")
	r, err := NewRequest(ctx, e)
	be.NilErr(t, err)
	v := r.Context().Value("key").(string)
	be.Equal(t, "value", v)
}
