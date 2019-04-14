package response

import (
	"encoding/json"
	"os"

	"github.com/torniker/goapp/logger"
)

// NewSub create an instance of sub responder
func NewSub() *Sub {
	return &Sub{
		headers: make(map[string]string),
	}
}

// Sub is type for sub responses
type Sub struct {
	headers   map[string]string
	status    int
	committed bool
	output    interface{}
}

// SetStatus sets response status code
func (s *Sub) SetStatus(status int) {
	s.status = status
}

// SetHeader sets header for responder
func (s *Sub) SetHeader(key, val string) {
	s.headers[key] = val
}

// Status returns response status code
func (s *Sub) Status() int {
	return s.status
}

// Commited returns response status
func (s *Sub) Commited() bool {
	return s.committed
}

// Write commits and writes data into the response body
func (s *Sub) Write(body interface{}) error {
	if s.status == 0 {
		s.SetStatus(200)
	}
	s.committed = true
	s.output = body
	logger.Infof("<--- status: %v, body: %#v", s.status, body)
	return json.NewEncoder(os.Stdout).Encode(body)
}

// Output returns response output
func (s *Sub) Output() interface{} {
	return s.output
}
