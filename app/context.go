package app

import (
	"encoding/json"
	"net/http"
)

// HandlerFunc defines handler function
type HandlerFunc func(*Ctx, string) error

// Ctx is struct where information for each request is stored
type Ctx struct {
	App          *App
	request      *http.Request
	response     http.ResponseWriter
	path         path
	hasResponded bool
}

// Method returns request method
func (ctx *Ctx) Method() string {
	return ctx.request.Method
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
	if ctx.hasResponded {
		return
	}
	ctx.path.index++
	err := f(ctx, ctx.path.Next())
	if err != nil {
		ctx.Error(err)
	}
}

// JSON responses with json body
func (ctx *Ctx) JSON(body interface{}) {
	ctx.hasResponded = true
	ctx.response.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(ctx.response).Encode(body)
	if err != nil {
		panic(err)
	}
}
