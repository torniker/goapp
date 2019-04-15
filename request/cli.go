package request

import (
	"encoding/json"
	"io"
	"net/url"
)

// NewCLI create an instance of CLI request
func NewCLI(a Action, url *url.URL, input io.Reader) *CLI {
	return &CLI{
		action: a,
		path:   NewPath(url),
		input:  input,
		flags:  make(map[string][]string),
	}
}

// CLI is type for CLI requests
type CLI struct {
	action Action
	path   *Path
	input  io.Reader
	flags  map[string][]string
}

// Action returns action type
func (c *CLI) Action() Action {
	return c.action
}

// Bind returns request body
func (c *CLI) Bind(v interface{}) error {
	decoder := json.NewDecoder(c.input)
	return decoder.Decode(&v)
}

// Flags returns cli flags
func (c *CLI) Flags() map[string][]string {
	return c.flags
}

// SetFlag sets flag
func (c *CLI) SetFlag(key, val string) {
	c.flags[key] = []string{val}
}

// Path returns url path
func (c *CLI) Path() *Path {
	return c.path
}
