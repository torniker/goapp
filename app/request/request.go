package request

import (
	"io"
	"net/http"
	"net/url"
)

// list of request types
const (
	TypeHTTP int = 1 + iota
	TypeCLI
)

// list of request methods
const (
	GET int = 1 << iota
	POST
	PUT
	DELETE
)

// Request defines methods for request object
type Request interface {
	Type() int
	Method() string
	Body() io.ReadCloser
	URL() *url.URL
}

// NewHTTP create an instance of HTTP request
func NewHTTP(r *http.Request) *HTTP {
	return &HTTP{
		req: r,
	}
}

// HTTP is type for HTTP requests
type HTTP struct {
	req *http.Request
}

// Type returns type of request
func (h *HTTP) Type() int {
	return TypeHTTP
}

// Method returns method
func (h *HTTP) Method() string {
	return h.req.Method
}

// Body returns request body
func (h *HTTP) Body() io.ReadCloser {
	return h.req.Body
}

// URL returns request url
func (h *HTTP) URL() *url.URL {
	return h.req.URL
}
