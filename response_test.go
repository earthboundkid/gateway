package gateway

import (
	"bytes"
	"testing"

	"github.com/earthboundkid/assert"
)

func Test_JSON_isTextMime(t *testing.T) {
	be := assert.FailsNow(t)
	be.True(isTextMime("application/json"))
	be.True(isTextMime("application/json; charset=utf-8"))
	be.True(isTextMime("Application/JSON"))
}

func Test_XML_isTextMime(t *testing.T) {
	be := assert.FailsNow(t)
	be.True(isTextMime("application/xml"))
	be.True(isTextMime("application/xml; charset=utf-8"))
	be.True(isTextMime("ApPlicaTion/xMl"))
}

func TestResponseWriter_Header(t *testing.T) {
	be := assert.FailsNow(t)
	w := NewResponse()
	w.Header().Set("Foo", "bar")
	w.Header().Set("Bar", "baz")

	var buf bytes.Buffer
	w.header.Write(&buf)

	be.Equal("Bar: baz\r\nFoo: bar\r\n", buf.String())
}

func TestResponseWriter_multiHeader(t *testing.T) {
	be := assert.FailsNow(t)
	w := NewResponse()
	w.Header().Set("Foo", "bar")
	w.Header().Set("Bar", "baz")
	w.Header().Add("X-APEX", "apex1")
	w.Header().Add("X-APEX", "apex2")

	var buf bytes.Buffer
	w.header.Write(&buf)

	be.Equal("Bar: baz\r\nFoo: bar\r\nX-Apex: apex1\r\nX-Apex: apex2\r\n", buf.String())
}

func TestResponseWriter_Write_text(t *testing.T) {
	types := []string{
		"text/x-custom",
		"text/plain",
		"text/plain; charset=utf-8",
		"application/json",
		"application/json; charset=utf-8",
		"application/xml",
		"image/svg+xml",
	}

	for _, kind := range types {
		t.Run(kind, func(t *testing.T) {
			be := assert.FailsNow(t)
			w := NewResponse()
			w.Header().Set("Content-Type", kind)
			w.Header().Set("Double-Header", "1")
			w.Header().Add("double-header", "2")
			w.Write([]byte("hello world\n"))

			e := w.End()
			be.Equal(200, e.StatusCode)
			be.Equal("hello world\n", e.Body)
			be.Equal(kind, e.Headers["Content-Type"])
			be.SlicesEqual([]string{"1", "2"}, e.MultiValueHeaders["Double-Header"])
			be.False(e.IsBase64Encoded)
			be.True(<-w.CloseNotify())
		})
	}
}

func TestResponseWriter_Write_binary(t *testing.T) {
	be := assert.FailsNow(t)
	w := NewResponse()
	w.Header().Set("Content-Type", "image/png")
	w.Write([]byte("data"))

	e := w.End()
	be.Equal(200, e.StatusCode)
	be.Equal("ZGF0YQ==", e.Body)
	be.Equal("image/png", e.Headers["Content-Type"])
	be.True(e.IsBase64Encoded)
}

func TestResponseWriter_Write_gzip(t *testing.T) {
	be := assert.FailsNow(t)
	w := NewResponse()
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Encoding", "gzip")
	w.Write([]byte("data"))

	e := w.End()
	be.Equal(200, e.StatusCode)
	be.Equal("ZGF0YQ==", e.Body)
	be.Equal("text/plain", e.Headers["Content-Type"])
	be.True(e.IsBase64Encoded)
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	be := assert.FailsNow(t)

	w := NewResponse()
	w.WriteHeader(404)
	w.Write([]byte("Not Found\n"))

	e := w.End()
	be.Equal(404, e.StatusCode)
	be.Equal("Not Found\n", e.Body)
	be.Equal("text/plain; charset=utf8", e.Headers["Content-Type"])
}
