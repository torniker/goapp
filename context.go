package wrap

import (
	"github.com/torniker/wrap/request"
	"github.com/torniker/wrap/response"
)

// HandlerFunc defines handler function
type HandlerFunc func(*Ctx) error

// Ctx is struct where information for each request is stored
type Ctx struct {
	Prog     *Prog
	Request  request.Request
	Response response.Response
	Store    map[string]interface{}
	User     Userer
}

// Post handles checks if the request method and calls HandlerFunc
func (ctx *Ctx) Post(f HandlerFunc) {
	if ctx.Request.Action() == request.POST {
		ctx.call(f)
	}
}

// Get handles checks if the request method and calls HandlerFunc
func (ctx *Ctx) Get(f HandlerFunc) {
	if ctx.Request.Action() == request.GET {
		ctx.call(f)
	}
}

// Put handles checks if the request method and calls HandlerFunc
func (ctx *Ctx) Put(f HandlerFunc) {
	if ctx.Request.Action() == request.PUT {
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
