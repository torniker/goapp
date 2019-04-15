package app

import (
	"github.com/torniker/goapp/request"
	"github.com/torniker/goapp/response"
)

// HandlerFunc defines handler function
type HandlerFunc func(*Ctx) error

// Ctx is struct where information for each request is stored
type Ctx struct {
	App      *App
	Request  request.Request
	Response response.Response
}

// Create handles checks if the request method and calls HandlerFunc
func (ctx *Ctx) Create(f HandlerFunc) {
	if ctx.Request.Action() == request.CREATE {
		ctx.call(f)
	}
}

// Read handles checks if the request method and calls HandlerFunc
func (ctx *Ctx) Read(f HandlerFunc) {
	if ctx.Request.Action() == request.READ {
		ctx.call(f)
	}
}

// Update handles checks if the request method and calls HandlerFunc
func (ctx *Ctx) Update(f HandlerFunc) {
	if ctx.Request.Action() == request.UPDATE {
		ctx.call(f)
	}
}

// Delete handles checks if the request method and calls HandlerFunc
func (ctx *Ctx) Delete(f HandlerFunc) {
	if ctx.Request.Action() == request.DELETE {
		ctx.call(f)
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
	ctx.Request.Path().Increment()
	ctx.call(f)
}

// JSON responses with json body
func (ctx *Ctx) JSON(body interface{}) error {
	ctx.Response.SetHeader("Content-Type", "application/json")
	return ctx.Response.Write(body)
}
