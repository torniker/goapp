package app

import (
	"io"

	"github.com/torniker/goapp/app/request"
	"github.com/torniker/goapp/app/response"
)

// HandlerFunc defines handler function
type HandlerFunc func(*Ctx, string) error

// Ctx is struct where information for each request is stored
type Ctx struct {
	App      *App
	request  request.Request
	response response.Response
	path     *path
}

// Method returns request method
func (ctx *Ctx) Method() string {
	return ctx.request.Method()
}

// RequestBody returns request body
func (ctx *Ctx) RequestBody() io.ReadCloser {
	return ctx.request.Body()
}

// Segment returns URL segment by index if exists or empty string
func (ctx *Ctx) Segment(index int) string {
	if len(ctx.path.segments) < index+1 {
		return ""
	}
	return ctx.path.segments[index]
}

// SetSegmentIndex sets current segment index
func (ctx *Ctx) SetSegmentIndex(index int) {
	ctx.path.index = index
}

// Next calls pathed function with next url segment and increases segment index by 1
func (ctx *Ctx) Next(f HandlerFunc) {
	if ctx.response.Commited() {
		return
	}
	ctx.path.index++
	err := f(ctx, ctx.path.Next())
	if err != nil {
		ctx.Error(err)
	}
}

// JSON responses with json body
func (ctx *Ctx) JSON(body interface{}) error {
	ctx.response.SetHeader("Content-Type", "application/json")
	return ctx.response.Write(body)
}
