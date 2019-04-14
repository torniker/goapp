package app

import (
	"github.com/torniker/goapp/request"
	"github.com/torniker/goapp/response"
)

// HandlerFunc defines handler function
type HandlerFunc func(*Ctx) error

// Ctx is struct where information for each request is stored
type Ctx struct {
	App         *App
	Request     request.Request
	Response    response.Response
	CurrentPath *path

	handler     *HandlerFunc
	elseHandler *HandlerFunc
}

// Method returns request method
func (ctx *Ctx) Method() string {
	switch ctx.Request.MethodCode() {
	case request.GET:
		return "GET"
	case request.POST:
		return "POST"
	case request.PUT:
		return "PUT"
	case request.DELETE:
		return "DELETE"
	}
	return ""
}

// Segment returns URL segment by index if exists or empty string
func (ctx *Ctx) Segment(index int) string {
	if len(ctx.CurrentPath.segments) < index+1 {
		return ""
	}
	return ctx.CurrentPath.segments[index]
}

// SetSegmentIndex sets current segment index
func (ctx *Ctx) SetSegmentIndex(index int) {
	ctx.CurrentPath.index = index
}

// GET handles checks if the request method and calls HandlerFunc
func (ctx *Ctx) GET(f HandlerFunc) {
	if ctx.Request.MethodCode()&request.GET != 0 {
		ctx.handler = &f
	}
}

// POST handles checks if the request method and calls HandlerFunc
func (ctx *Ctx) POST(f HandlerFunc) {
	if ctx.Request.MethodCode()&request.POST != 0 {
		ctx.handler = &f
	}
}

// PUT handles checks if the request method and calls HandlerFunc
func (ctx *Ctx) PUT(f HandlerFunc) {
	if ctx.Request.MethodCode()&request.PUT != 0 {
		ctx.handler = &f
	}
}

// DELETE handles checks if the request method and calls HandlerFunc
func (ctx *Ctx) DELETE(f HandlerFunc) {
	if ctx.Request.MethodCode()&request.DELETE != 0 {
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
	if ctx.Response.Commited() {
		return
	}
	err := f(ctx)
	if err != nil {
		ctx.Error(err)
	}
}

// Next calls pathed function with next url segment and increases segment index by 1
func (ctx *Ctx) Next(f HandlerFunc) {
	ctx.CurrentPath.index++
	ctx.call(f)
}

// JSON responses with json body
func (ctx *Ctx) JSON(body interface{}) error {
	ctx.Response.SetHeader("Content-Type", "application/json")
	return ctx.Response.Write(body)
}
