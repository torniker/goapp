package request

import (
	"io"
	"net/http"
)

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

// MethodCode returns method
func (h *HTTP) MethodCode() int {
	switch h.req.Method {
	case "GET":
		return GET
	case "POST":
		return POST
	case "PUT":
		return PUT
	case "DELETE":
		return DELETE
	}
	return 0
}

// Input returns request body
func (h *HTTP) Input() io.ReadCloser {
	return h.req.Body
}

// Query returns request Query params
func (h *HTTP) Query() map[string][]string {
	return h.req.URL.Query()
}

// Path returns request path
func (h *HTTP) Path() string {
	return h.req.URL.Path
}
