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

	handler     *HandlerFunc
	elseHandler *HandlerFunc
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

// GET handles checks if the request method and calls HandlerFunc
func (ctx *Ctx) GET(f HandlerFunc) {
	if ctx.request.MethodCode() == request.GET {
		ctx.handler = &f
	}
}

// POST handles checks if the request method and calls HandlerFunc
func (ctx *Ctx) POST(f HandlerFunc) {
	if ctx.request.MethodCode() == request.POST {
		ctx.handler = &f
	}
}

// PUT handles checks if the request method and calls HandlerFunc
func (ctx *Ctx) PUT(f HandlerFunc) {
	if ctx.request.MethodCode() == request.PUT {
		ctx.handler = &f
	}
}

// DELETE handles checks if the request method and calls HandlerFunc
func (ctx *Ctx) DELETE(f HandlerFunc) {
	if ctx.request.MethodCode() == request.DELETE {
		ctx.handler = &f
	}
}

// ELSE handles checks if the request method and calls HandlerFunc
func (ctx *Ctx) ELSE(f HandlerFunc) {
	ctx.elseHandler = &f
}

// Do calls handler func is not nil or fallbacks to default handler
func (ctx *Ctx) Do() {
	if ctx.handler != nil {
		ctx.call(*ctx.handler)
	} else if ctx.elseHandler != nil {
		ctx.call(*ctx.elseHandler)
	} else {
		ctx.call(ctx.App.DefaultHandler)
	}
}

func (ctx *Ctx) call(f HandlerFunc) {
	if ctx.response.Commited() {
		return
	}
	err := f(ctx, ctx.path.Next())
	if err != nil {
		ctx.Error(err)
	}
}

// Next calls pathed function with next url segment and increases segment index by 1
func (ctx *Ctx) Next(f HandlerFunc) {
	ctx.path.index++
	ctx.call(f)
}

// JSON responses with json body
func (ctx *Ctx) JSON(body interface{}) error {
	ctx.response.SetHeader("Content-Type", "application/json")
	return ctx.response.Write(body)
}
