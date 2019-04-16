package request

import (
	"encoding/json"
	"io"
	"net/url"

	"github.com/torniker/goapp/logger"
)

// NewRequest create an instance of sub request
func NewRequest(a Action, url *url.URL, input io.Reader) *Req {
	return &Req{
		action: a,
		path:   NewPath(url),
		input:  input,
		flags:  make(map[string][]string),
	}
}

// Req is type for sub-requests
type Req struct {
	action Action
	path   *Path
	input  io.Reader
	data   interface{}
	flags  map[string][]string
}

// Action returns action
func (r Req) Action() Action {
	return r.action
}

// Bind binds request input to pathed object
func (r Req) Bind(v interface{}) error {
	logger.Infof("input: %v", r.data)
	if r.data != nil {
		v = r.data
		return nil
	}
	decoder := json.NewDecoder(r.input)
	return decoder.Decode(&v)
}

// Flags returns optional params
func (r Req) Flags() map[string][]string {
	return r.flags
}

// Path returns url path
func (r Req) Path() *Path {
	return r.path
}

// SetAction sets action
func (r *Req) SetAction(a Action) *Req {
	r.action = a
	return r
}

// SetPath sets path
func (r *Req) SetPath(v *url.URL) *Req {
	r.path = NewPath(v)
	return r
}

// SetInput sets input
func (r *Req) SetInput(v io.Reader) *Req {
	r.input = v
	return r
}

// SetData sets data
func (r *Req) SetData(v interface{}) *Req {
	r.data = v
	return r
}

// SetFlags sets flags
func (r *Req) SetFlags(flags map[string][]string) *Req {
	r.flags = flags
	return r
}
