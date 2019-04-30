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

// ResJSON is a struct for wrapping a JSON data response
type ResJSON struct {
	Data interface{} `json:"data"`
}

// JSON responses with json body
func (ctx *Ctx) JSON(body interface{}) error {
	ctx.Response.SetHeader("Content-Type", "application/json")
	return ctx.Response.Write(ResJSON{Data: body})
}

// NoContent responses with status 204 No Content
func (ctx *Ctx) NoContent() {
	ctx.Response.SuccessWithNoContent()
}

// IsGET handles checks if the request method is GET
func (ctx *Ctx) IsGET() bool {
	if ctx.Request.Action() == request.GET {
		return true
	}
	return false
}

// IsPOST handles checks if the request method is POST
func (ctx *Ctx) IsPOST() bool {
	if ctx.Request.Action() == request.POST {
		return true
	}
	return false
}

// IsPUT handles checks if the request method is PUT
func (ctx *Ctx) IsPUT() bool {
	if ctx.Request.Action() == request.PUT {
		return true
	}
	return false
}

// IsDELETE handles checks if the request method is DELETE
func (ctx *Ctx) IsDELETE() bool {
	if ctx.Request.Action() == request.DELETE {
		return true
	}
	return false
}

// IsOPTIONS handles checks if the request method is OPTIONS
func (ctx *Ctx) IsOPTIONS() bool {
	if ctx.Request.Action() == request.OPTIONS {
		return true
	}
	return false
}
