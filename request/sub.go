package request

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
)

// NewSub create an instance of sub request
func NewSub(method string, url *url.URL, input string) *Sub {
	return &Sub{
		method: method,
		url:    url,
		input:  input,
	}
}

// Sub is type for sub-requests
type Sub struct {
	method string
	url    *url.URL
	input  string
}

// Type returns type of request
func (s *Sub) Type() int {
	return TypeCLI
}

// MethodCode returns method
func (s *Sub) MethodCode() int {
	switch strings.ToUpper(s.method) {
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
func (s *Sub) Input() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader([]byte(s.input)))
}

// Query returns cli flags
func (s *Sub) Query() map[string][]string {
	return s.url.Query()
}

// Path returns url path
func (s *Sub) Path() string {
	return s.url.Path
}
