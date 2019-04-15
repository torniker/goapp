package request

import (
	"fmt"
	"net/url"
	"strings"
)

// list of request actions
const (
	CREATE Action = 1 + iota
	READ
	UPDATE
	DELETE
)

// Action type for describing request action
type Action int

// Request defines methods for request object
type Request interface {
	Action() Action
	Path() *Path
	Bind(v interface{}) error
	SetFlag(key, value string)
	Flags() map[string][]string
}

// Path defines structure for command path
type Path struct {
	url      *url.URL
	segments []string
	index    int
}

// NewAction returns validated action
func NewAction(val int) (*Action, error) {
	a := Action(val)
	if !a.IsValid() {
		return nil, fmt.Errorf("invalid value of action")
	}
	return &a, nil
}

// IsValid checks if action has a valid value
func (a *Action) IsValid() bool {
	if a == nil {
		return false
	}
	if *a > 0 && *a < 5 {
		return true
	}
	return false
}

// NewActionFromString returns action
func NewActionFromString(a string) Action {
	switch a {
	case "CREATE":
		return CREATE
	case "READ":
		return READ
	case "UPDATE":
		return UPDATE
	case "DELETE":
		return DELETE
	}
	return Action(0)
}

// string returns string name of action
func (a Action) String() string {
	switch a {
	case CREATE:
		return "CREATE"
	case READ:
		return "READ"
	case UPDATE:
		return "UPDATE"
	case DELETE:
		return "DELETE"
	}
	return ""
}

// NewPath creates path object from url
func NewPath(url *url.URL) *Path {
	return &Path{
		url:      url,
		segments: strings.Split(url.Path, "/"),
		index:    0,
	}
}

// Next returns next segment of the path
func (p *Path) Next() string {
	if len(p.segments) < p.index+2 {
		return ""
	}
	return p.segments[p.index+1]
}

// Current returns current segment of the path
func (p *Path) Current() string {
	if len(p.segments) < p.index+1 {
		return ""
	}
	return p.segments[p.index]
}

// Increment increases path index by 1
func (p *Path) Increment() {
	p.index++
}
