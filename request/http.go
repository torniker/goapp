package request

import (
	"encoding/json"
	"net/http"
)

// NewHTTP create an instance of HTTP request
func NewHTTP(r *http.Request) *HTTP {
	return &HTTP{
		path: NewPath(r.URL),
		req:  r,
	}
}

// HTTP is type for HTTP requests
type HTTP struct {
	path *Path
	req  *http.Request
}

// Action returns method
func (h HTTP) Action() Action {
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

// Bind build object from input
func (h HTTP) Bind(v interface{}) error {
	decoder := json.NewDecoder(h.req.Body)
	return decoder.Decode(&v)
}

// Flags returns request Query params
func (h HTTP) Flags() map[string][]string {
	return h.req.URL.Query()
}

// Path returns request path
func (h HTTP) Path() *Path {
	return h.path
}
