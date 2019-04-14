package request

import (
	"io"
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
	MethodCode() int
	Input() io.ReadCloser
	Path() string
	Query() map[string][]string
}
