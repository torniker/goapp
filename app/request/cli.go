package request

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
)

// NewCLI create an instance of CLI request
func NewCLI(method string, url *url.URL, input string) *CLI {
	return &CLI{
		method: method,
		url:    url,
		input:  input,
	}
}

// CLI is type for CLI requests
type CLI struct {
	method string
	url    *url.URL
	input  string
}

// Type returns type of request
func (c *CLI) Type() int {
	return TypeCLI
}

// MethodCode returns method
func (c *CLI) MethodCode() int {
	switch strings.ToUpper(c.method) {
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
func (c *CLI) Input() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader([]byte(c.input)))
}

// Query returns cli flags
func (c *CLI) Query() map[string][]string {
	return c.url.Query()
}

// Path returns url path
func (c *CLI) Path() string {
	return c.url.Path
}
