package response

import (
	"encoding/json"
	"net/http"

	"github.com/torniker/goapp/app/logger"
)

// list of response types
const (
	TypeHTTP int = 1 + iota
	TypeCLI
)

// Response defines methods for response object
type Response interface {
	SetStatus(status int)
	Status() int
	SetHeader(key, val string)
	Commited() bool
	Write(interface{}) error
}

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
	logger.Infof("<--- status: %v, body: %#v", h.status, body)
	return json.NewEncoder(h.writer).Encode(body)
}
