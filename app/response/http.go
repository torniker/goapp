package response

import (
	"encoding/json"
	"net/http"

	"github.com/torniker/goapp/app/logger"
)

// NewHTTP create an instance of HTTP responder
func NewHTTP(w http.ResponseWriter) *HTTP {
	return &HTTP{
		writer:  w,
		headers: make(map[string]string),
	}
}

// HTTP is type for HTTP responses
type HTTP struct {
	writer    http.ResponseWriter
	headers   map[string]string
	status    int
	committed bool
	output    interface{}
}

// SetStatus sets response status code
func (h *HTTP) SetStatus(status int) {
	h.status = status
	h.writer.WriteHeader(status)
}

// SetHeader sets header for responder
func (h *HTTP) SetHeader(key, val string) {
	h.headers[key] = val
	h.writer.Header().Add(key, val)
}

// Status returns response status code
func (h *HTTP) Status() int {
	return h.status
}

// Commited returns response status
func (h *HTTP) Commited() bool {
	return h.committed
}

// Write commits and writes data into the response body
func (h *HTTP) Write(body interface{}) error {
	h.committed = true
	h.output = body
	logger.Infof("<--- status: %v, body: %#v", h.status, body)
	return json.NewEncoder(h.writer).Encode(body)
}

// Output returns response output
func (h *HTTP) Output() interface{} {
	return h.output
}
