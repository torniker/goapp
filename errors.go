package wrap

import (
	"fmt"
	"net/http"

	"github.com/torniker/wrap/logger"
)

// Error checks error type and responses accordingly
func (ctx *Ctx) Error(err error) {
	switch err.(type) {
	case ErrorBadRequest:
		ctx.Response.SetStatus(http.StatusBadRequest)
		e := err.(ErrorBadRequest)
		ctx.Response.Write(e)
	case ErrorUnauthorized:
		ctx.Response.SetStatus(http.StatusUnauthorized)
		e := err.(ErrorUnauthorized)
		ctx.Response.Write(e)
	case ErrorMethodNotAllowed:
		ctx.Response.SetStatus(http.StatusMethodNotAllowed)
		e := err.(ErrorMethodNotAllowed)
		ctx.Response.Write(e)
	case ErrorStatusNotFound:
		ctx.Response.SetStatus(http.StatusNotFound)
		e := err.(ErrorStatusNotFound)
		ctx.Response.Write(e)
	case ErrorInternalServerError:
		ctx.Response.SetStatus(http.StatusInternalServerError)
		e := err.(ErrorInternalServerError)
		ctx.Response.Write(e)
	case ErrorUnprocessableEntity:
		ctx.Response.SetStatus(http.StatusUnprocessableEntity)
		e := err.(ErrorUnprocessableEntity)
		ctx.Response.Write(e)
	default:
		ctx.Response.SetStatus(http.StatusInternalServerError)
		ctx.Response.Write(ErrorInternalServerError{Message: err.Error()})
	}
}

// NotFound response 404
func (ctx *Ctx) NotFound() error {
	e := ErrorStatusNotFound{
		Message:  "not found",
		Internal: fmt.Sprintf("url: %v not found", ctx.Request.Path().URL().Path),
	}
	logger.Error(e.Internal)
	return e
}

// Unauthorized response 401
func (ctx *Ctx) Unauthorized() error {
	e := ErrorUnauthorized{
		Message:  "unauthorized",
		Internal: fmt.Sprintf("user: %v is unauthorized to request: %v", ctx.User, ctx.Request.Path().URL().Path),
	}
	logger.Error(e.Internal)
	return e
}

// InternalServerError response 500
func (ctx *Ctx) InternalServerError(err error) error {
	e := ErrorInternalServerError{
		Message:  "internal server error",
		Internal: err.Error(),
	}
	logger.Error(e.Internal)
	return e
}

// BadRequest response 400
func (ctx *Ctx) BadRequest(message string) error {
	e := ErrorBadRequest{
		Message:  message,
		Internal: fmt.Sprintf("bad request: %#v, message: %v", ctx.Request, message),
	}
	logger.Error(e.Internal)
	return e
}

// MethodNotAllowed response 405
func (ctx *Ctx) MethodNotAllowed() error {
	e := ErrorMethodNotAllowed{
		Message:  "not allowed",
		Internal: fmt.Sprintf("user: %v is not allowed to request: %v", ctx.User, ctx.Request.Path().URL().Path),
	}
	logger.Error(e.Internal)
	return e
}

// UnprocessableEntity respnse 422
func (ctx *Ctx) UnprocessableEntity(code int, message string) error {
	e := ErrorUnprocessableEntity{
		Code:     code,
		Message:  message,
		Internal: fmt.Sprintf("UnprocessableEntity code: %v, message: %v ", code, message),
	}
	logger.Error(e.Internal)
	return e
}

// ErrorBadRequest type for bad request
type ErrorBadRequest struct {
	Message  string `json:"message"`
	Internal string `json:"-"`
}

func (e ErrorBadRequest) Error() string {
	return e.Message
}

// ErrorUnauthorized type for Unauthorized
type ErrorUnauthorized struct {
	Message  string `json:"message"`
	Internal string `json:"-"`
}

func (e ErrorUnauthorized) Error() string {
	return e.Message
}

// ErrorMethodNotAllowed type for not allowed
type ErrorMethodNotAllowed struct {
	Message  string `json:"message"`
	Internal string `json:"-"`
}

func (e ErrorMethodNotAllowed) Error() string {
	return e.Message
}

// ErrorStatusNotFound type for not found
type ErrorStatusNotFound struct {
	Message  string `json:"message"`
	Internal string `json:"-"`
}

func (e ErrorStatusNotFound) Error() string {
	return e.Message
}

// ErrorInternalServerError type for internal server error
type ErrorInternalServerError struct {
	Message  string `json:"message"`
	Internal string `json:"-"`
}

func (e ErrorInternalServerError) Error() string {
	return e.Message
}

// ErrorUnprocessableEntity type for validation errors
type ErrorUnprocessableEntity struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Internal string `json:"-"`
}

func (e ErrorUnprocessableEntity) Error() string {
	return e.Message
}
